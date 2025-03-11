package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sanjay-xdr/github-dashboard/backend/internals/handlers"
)

func Routes() http.Handler {

	r := chi.NewRouter()
	// CORS setup
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins is a list of origins a cross-domain request can be executed from.
		// If the special "*" value is present in the list, all origins will be allowed.
		// An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com).
		AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},

		// AllowOriginFunc is a custom function to validate the origin. It takes the origin as argument
		// and returns true if allowed or false otherwise. If this option is set, AllowedOrigins is ignored.
		// AllowOriginFunc: func(origin string) bool { return true },

		// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// AllowedHeaders is a list of non-simple headers the client is allowed to use with cross-domain requests.
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

		// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification.
		ExposedHeaders: []string{"Link"},

		// AllowCredentials indicates whether the request can include user credentials like cookies, HTTP authentication or client-side SSL certificates.
		AllowCredentials: true,

		// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
		MaxAge: 300, // 5 minutes
	}))
	
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
