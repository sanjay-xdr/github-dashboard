package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/handlers"
)

func Routes() http.Handler {

	r := chi.NewRouter()

	v1 := chi.NewRouter()
	v1.Get("/pullrequestdata", handlers.GetPullRequestData)
	v1.Get("/repodata", handlers.GetRepoData)
	v1.Post("/merged-prs", handlers.GetMergedPRByDate)

	r.Mount("/api/v1", v1)

	/* example of creating v2 of this
	v2 := chi.NewRouter()
	r.Mount("/api/v2", v2)
	*/

	return r
}
