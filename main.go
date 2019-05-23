package main

import (
	"net/http"

	"github.com/thejayhaykid/home-library/models"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v1"

	"encoding/json"
	"strconv"

	"github.com/gobuffalo/envy"
	sessions "github.com/goincremental/negroni-sessions"
	gmux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/yosssi/ace"
)

// Global Variables
var db *sql.DB
var dbmap *gorp.DbMap

// Note: Don't store your key in your source code. Pass it via an
// environmental variable, or flag (or both), and don't accidentally commit it
// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
var store = sessions.NewCookieStore(envy.Get("SESSION_KEY"))

// LoginPage determines if login page contains an error
type LoginPage struct {
	Error string
}

func initDb() {
	db, _ = sql.Open("sqlite3", "dev.db")

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(models.Book{}, "books").SetKeys(true, "pk")
	dbmap.AddTableWithName(models.User{}, "users").SetKeys(false, "username")
	dbmap.CreateTablesIfNotExists()
}

func verifyDatabase(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

func getBookCollection(books *[]models.Book, sortCol string, filterByClass string, username string, w http.ResponseWriter) bool {
	if sortCol == "" {
		sortCol = "pk"
	}
	where := " where user=?"
	if filterByClass == "fiction" {
		where += " and classification between '800' and '900'"
	} else if filterByClass == "nonfiction" {
		where += " and classification not between '800' and '900'"
	}
	if _, err := dbmap.Select(books, "select * from books"+where+" order by "+sortCol, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

func getStringFromSession(r *http.Request, key string) string {
	var strVal string
	session, _ := store.Get(r, "session-name")

	if val := session.Values[key]; val != nil {
		strVal = val.(string)
	}
	return strVal
}

func verifyUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/login" {
		next(w, r)
		return
	}

	if username := getStringFromSession(r, "User"); username != "" {
		if user, _ := dbmap.Get(models.User{}, username); user != nil {
			next(w, r)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// LogoutHandler will log a user out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session.Values["User"] = nil
	session.Values["Filter"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	initDb()

	mux := gmux.NewRouter()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := models.Page{Books: []models.Book{}, Filter: getStringFromSession(r, "Filter"), models.User: getStringFromSession(r, "User")}
		if !getBookCollection(&p.Books, getStringFromSession(r, "SortBy"), getStringFromSession(r, "Filter"), p.User, w) {
			return
		}

		if err = template.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var results []models.SearchResult
		var err error

		if results, err = search(r.FormValue("search")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		var book models.ClassifyBookResponse
		var err error

		if book, err = find(r.FormValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		b := models.Book{
			PK:             -1,
			Title:          book.BookData.Title,
			Author:         book.BookData.Author,
			Classification: book.Classification.MostPopular,
			ID:             r.FormValue("id"),
			User:           getStringFromSession(r, "User"),
		}

		if err = dbmap.Insert(&b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("PUT")

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		var b []models.Book
		if !getBookCollection(&b, r.FormValue("sortBy"), getStringFromSession(r, "Filter"), getStringFromSession(r, "User"), w) {
			return
		}

		sessions.GetSession(r).Set("SortBy", r.FormValue("sortBy"))

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET").Queries("sortBy", "{sortBy:title|author|classification}")

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		var b []models.Book
		if !getBookCollection(&b, getStringFromSession(r, "SortBy"), r.FormValue("filter"), getStringFromSession(r, "User"), w) {
			return
		}

		sessions.GetSession(r).Set("Filter", r.FormValue("filter"))

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET").Queries("filter", "{filter:all|fiction|nonfiction}")

	mux.HandleFunc("/books/{pk}", func(w http.ResponseWriter, r *http.Request) {
		pk, _ := strconv.ParseInt(gmux.Vars(r)["pk"], 10, 64)
		var b models.Book
		if err := dbmap.SelectOne(&b, "select * from books where pk=? and user=?", pk, getStringFromSession(r, "User")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if _, err := dbmap.Delete(&b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var p LoginPage
		if r.FormValue("register") != "" {
			secret, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
			user := models.User{r.FormValue("username"), secret}
			if err := dbmap.Insert(&user); err != nil {
				p.Error = err.Error()
			} else {
				session.GetSession(r).Set("User", user.Username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		} else if r.FormValue("login") != "" {
			user, err := dbmap.Get(models.User{}, r.FormValue("username"))
			if err != nil {
				p.Error = err.Error()
			} else if user == nil {
				p.Error = "No such user found with Username: " + r.FormValue("username")
			} else {
				u := user.(*models.User)
				if err = bcrypt.CompareHashAndPassword(u.Secret, []byte(r.FormValue("password"))); err != nil {
					p.Error = err.Error()
				} else {
					sessions.GetSession(r).Set("User", u.Username)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
		}

		template, err := ace.Load("templates/login", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = template.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/logout", LougoutHandler)

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(envy.Get("SESSION_KEY"))))
	e.Use(verifyDatabase)
	e.Use(verifyUser)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.UseHandler(mux)
	e.Static("/", "/assets")
	e.Logger.Fatal(e.Start(":8080"))
}
