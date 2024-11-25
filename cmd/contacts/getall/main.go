package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/starry-axul/axul-contact/internal/contact"
	"github.com/starry-axul/axul-contact/pkg/handler"
)

func main() {

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	repo := contact.NewRepo(nil, nil)
	service := contact.NewService(repo, nil, nil, nil)
	app := handler.NewGetAllHandler(contact.MakeEndpoints(service))
	lambda.StartHandler(app)

	log.Println(<-errs)
}
