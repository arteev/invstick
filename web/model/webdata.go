package model

//A Webdata data for rendering
type Webdata struct {
	Locale    string
	Locales   []string
	Templates []string
}

//Data for rendering
var Data Webdata
