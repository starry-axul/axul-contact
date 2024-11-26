package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/starry-axul/axul-contact/internal/contact"
	"github.com/starry-axul/axul-contact/pkg/bootstrap"
	"github.com/starry-axul/axul-contact/pkg/handler"
	"github.com/go-kit/kit/transport/awslambda"
	"gorm.io/gorm"
)


var db *gorm.DB
var h *awslambda.Handler

func init() {
	logger := bootstrap.SetupLogger()
	db = bootstrap.DBConnection()

	repo := contact.NewRepo(db, logger)
	service := contact.NewService(repo, nil, nil, logger)
	h = handler.NewGetHandler(contact.MakeEndpoints(service))
}

func main() {
	lambda.StartHandler(h)
}