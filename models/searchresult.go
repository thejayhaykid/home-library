package models

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SearchResult gets values that are displayed on the page from XML received
type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

// ClassifySearchResponse Puts received XML into struct
type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

// Search for a specific book
func Search(query string) ([]SearchResult, error) {
	var c ClassifySearchResponse
	body, err := ClassifyAPI("http://classify.oclc.org/classify2/Classify?summary=true&title=" + url.QueryEscape(query))

	if err != nil {
		return []SearchResult{}, err
	}

	err = xml.Unmarshal(body, &c)
	return c.Results, err
}

// ClassifyAPI will classify the API
func ClassifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(url); err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
