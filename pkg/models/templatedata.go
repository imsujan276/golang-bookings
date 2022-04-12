package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap      map[string]string
	IntMap         map[string]int
	FloatMap       map[string]float64
	Data           map[string]interface{} // any type of data
	CSRFToken      string
	FlashMessage   string
	WarningMessage string
	ErrorMessage   string
}
