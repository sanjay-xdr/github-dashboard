package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/handlers"
)

func Routes() http.Handler {

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	v1.Get("/pullrequestdata", handlers.GetPullRequestData)
	v1.Get("/repodata", handlers.GetRepoData)
	v1.Post("/merged-prs", handlers.GetMergedPRByDate)
	v1.Get("/fetch-prs", handlers.FetchPRData)
	v1.Get("/fetch-workflow", handlers.FetchWorkflowData)

	r.Mount("/api/v1", v1)

	/* example of creating v2 of this
	v2 := chi.NewRouter()
	r.Mount("/api/v2", v2)
	*/

	return r
}
