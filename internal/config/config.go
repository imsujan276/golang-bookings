package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/imsujan276/golang-bookings/internal/models"
)

// AppConfig holds the application config i.e. global
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChannel   chan models.MailData
}
