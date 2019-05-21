package models

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
