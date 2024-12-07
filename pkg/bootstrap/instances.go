package bootstrap

import (
	"gorm.io/gorm"
	"github.com/ncostamagna/go-logger-hub/loghub"
	"github.com/starry-axul/axul-contact/internal/contact"
	"github.com/starry-axul/axul-contact/pkg/notify"
	"os"
)

func ContactInstance(db *gorm.DB, logger loghub.Logger) contact.Endpoints {
	n := notify.NewHttpClient(os.Getenv("NOTIFY_URL"), os.Getenv("NOTIFY_TOKEN"))
	repo := contact.NewRepo(db, logger)
	service := contact.NewService(repo, n, nil, logger)
	return contact.MakeEndpoints(service)
}
