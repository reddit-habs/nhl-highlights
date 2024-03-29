package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sbstp/nhl-highlights/generate"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/render"
	"github.com/sbstp/nhl-highlights/repository"
)

func serve(ctx context.Context, bindAddress string, incremental bool) error {
	repo, err := repository.New("games.db")
	if err != nil {
		return err
	}
	defer repo.Close()

	api := nhlapi.NewClient()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/nhl", func(r chi.Router) {
		r.Get("/clips/{gamePk:[0-9]+}", htmlOrError(func(r *http.Request) (render.Render, error) {
			gamePk := chi.URLParam(r, "gamePk")
			gameID, _ := strconv.ParseInt(gamePk, 10, 64)
			return renderClips(api, gameID)
		}))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/nhl/current/", http.StatusFound)
		})

		r.Get("/current", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/nhl/current/", http.StatusFound)
		})

		r.Get("/current/", htmlOrError(func(r *http.Request) (render.Render, error) {
			return renderCachedPage(repo, nil, nil)
		}))

		r.Get("/current/index.html", htmlOrError(func(r *http.Request) (render.Render, error) {
			return renderCachedPage(repo, nil, nil)
		}))

		r.Get("/current/{team}.html", htmlOrError(func(r *http.Request) (render.Render, error) {
			team := chi.URLParam(r, "team")
			return renderCachedPage(repo, nil, &team)
		}))

		r.Get("/{season}/", htmlOrError(func(r *http.Request) (render.Render, error) {
			season := chi.URLParam(r, "season")
			return renderCachedPage(repo, &season, nil)
		}))

		r.Get("/{season}/index.html", htmlOrError(func(r *http.Request) (render.Render, error) {
			season := chi.URLParam(r, "season")
			return renderCachedPage(repo, &season, nil)
		}))

		r.Get("/{season}/{team}.html", htmlOrError(func(r *http.Request) (render.Render, error) {
			season := chi.URLParam(r, "season")
			team := chi.URLParam(r, "team")
			return renderCachedPage(repo, &season, &team)
		}))
	})

	go http.ListenAndServe(bindAddress, r)

	if incremental {
		go startIncrementalArchiveTimer(ctx)
	}

	<-ctx.Done()
	return nil
}

func htmlOrError(wrapped func(*http.Request) (render.Render, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := wrapped(r)
		if err != nil {
			log.Printf("[error] %v", err)
			http.Error(w, "Internal server error occured. The error has been logged.", http.StatusInternalServerError)
			return
		}
		content.Render(w, http.StatusOK)
	})

}

func renderClips(api nhlapi.Client, gameID int64) (render.Render, error) {
	clips, err := scanClips(api, gameID)
	if err != nil {
		return nil, err
	}
	data, err := generate.Clips(clips)
	if err != nil {
		return nil, err
	}
	return render.HTML(data), nil
}

func renderCachedPage(repo *repository.Repository, season *string, team *string) (render.Render, error) {
	if season == nil {
		x, err := repo.GetCurrentSeason()
		if err != nil {
			return nil, err
		}
		season = &x
	}
	cp, err := repo.GetCachedPage(*season, team)
	if err != nil {
		return nil, err
	}
	return render.CompressedHTML(cp.Content), nil
}

func startIncrementalArchiveTimer(ctx context.Context) {
	doArchival := func() {
		if err := archive(true, "", ""); err != nil {
			log.Printf("[error] archival timer error: %v", err)
		}
	}

	delay := time.After(time.Second * 5)
	ticker := time.NewTicker(time.Minute * 15)
	defer ticker.Stop()

	for {
		select {
		case <-delay:
			// archive 5 second after starting the time (i.e. on program start)
			doArchival()
		case <-ticker.C:
			// archive every 15 minutes
			doArchival()
		case <-ctx.Done():
			// program terminated, stop timer
			return
		}
	}
}
