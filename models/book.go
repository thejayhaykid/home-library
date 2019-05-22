package models

import (
	"encoding/xml"
	"home-library/models"
	"net/url"
)

// Book is the default struct of a book
type Book struct {
	PK             int64  `db:"pk"`
	Title          string `db:"title"`
	Author         string `db:"author"`
	Classification string `db:"classification"`
	ID             string `db:"id"`
	User           string `db:"user"`
}

// ClassifyBookResponse will Get pertinent information from XML
type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

// Find a book from a string
func Find(id string) (models.ClassifyBookResponse, error) {
	var c models.ClassifyBookResponse
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?summary=true&owi=" + url.QueryEscape(id))

	if err != nil {
		return models.ClassifyBookResponse{}, err
	}

	err = xml.Unmarshal(body, &c)
	return c, err
}
