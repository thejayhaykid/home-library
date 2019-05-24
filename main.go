package main

import (
	"github.com/thejayhaykid/home-library/actions"
)

func main() {
	e := actions.Start()
	e.Logger.Fatal(e.Start(":1212"))
}
