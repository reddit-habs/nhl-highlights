package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sbstp/nhl-highlights/repository"
)

func serve(bindAddress string) error {
	repo, err := repository.New("games.db")
	if err != nil {
		return err
	}
	defer repo.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/nhl", func(r chi.Router) {
		r.Get("/{season}/", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			renderCachedPage(w, repo, season, nil)

		})

		r.Get("/{season}/index.html", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			renderCachedPage(w, repo, season, nil)

		})

		r.Get("/{season}/{team}.html", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			team := chi.URLParam(r, "team")
			renderCachedPage(w, repo, season, &team)
		})
	})

	return http.ListenAndServe(bindAddress, r)
}

func renderCachedPage(w http.ResponseWriter, repo *repository.Repository, season string, team *string) {
	w.Header().Set("Content-Type", "text/html")
	cp, err := repo.GetCachedPage(season, team)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(cp.Content)
}
