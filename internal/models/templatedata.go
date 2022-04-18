package models

import "github.com/imsujan276/golang-bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{} // any type of data
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
