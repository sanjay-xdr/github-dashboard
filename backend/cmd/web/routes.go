package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/handlers"
)

func Routes() http.Handler {

	r := chi.NewRouter()

	r.Get("/pullrequestdata", handlers.GetPullRequestData)
	r.Get("/repodata", handlers.GetRepoData)
	r.Get("/testresult", handlers.GetTestResult)

	return r
}
