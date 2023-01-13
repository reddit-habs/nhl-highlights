package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/repository"
)

func serve(bindAddress string) error {
	repo, err := repository.New("games.db")
	if err != nil {
		return err
	}
	defer repo.Close()

	api := nhlapi.NewClient()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/nhl", func(r chi.Router) {
		r.Get("/clips/{gamePk:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			gamePk := chi.URLParam(r, "gamePk")
			gameID, _ := strconv.ParseInt(gamePk, 10, 64)
			renderClips(w, api, gameID)
		})

		r.Get("/current", func(w http.ResponseWriter, r *http.Request) {
			renderCachedPage(w, repo, nil, nil)
		})

		r.Get("/{season}/", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			renderCachedPage(w, repo, &season, nil)
		})

		r.Get("/{season}/index.html", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			renderCachedPage(w, repo, &season, nil)
		})

		r.Get("/{season}/{team}.html", func(w http.ResponseWriter, r *http.Request) {
			season := chi.URLParam(r, "season")
			team := chi.URLParam(r, "team")
			renderCachedPage(w, repo, &season, &team)
		})
	})

	return http.ListenAndServe(bindAddress, r)
}

func renderClips(w http.ResponseWriter, api nhlapi.Client, gameID int64) {
	clips, err := scanClips(api, gameID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(clips)
}

func renderCachedPage(w http.ResponseWriter, repo *repository.Repository, season *string, team *string) {
	if season == nil {
		x, _ := repo.GetCurrentSeason()
		log.Printf(x)
		season = &x
	}
	w.Header().Set("Content-Type", "text/html")
	cp, err := repo.GetCachedPage(*season, team)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(cp.Content)
}
