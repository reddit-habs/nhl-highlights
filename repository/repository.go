package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sbstp/nhl-highlights/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const pragmas = `
PRAGMA foreign_keys = ON;
`

const schema = `
CREATE TABLE IF NOT EXISTS games (
	game_id INTEGER PRIMARY KEY NOT NULL,
	date TEXT NOT NULL,
	type TEXT NOT NULL,
	away TEXT NOT NULL,
	home TEXT NOT NULL,
	season TEXT NOT NULL,
	recap TEXT,
	extended TEXT
);

CREATE TABLE IF NOT EXISTS highlights (
	id INTEGER NOT NULL,
	game_id INTEGER NOT NULL,
	media_url TEXT NOT NULL,
	event_id INT,
	title TEXT,
	blurb TEXT,
	description TEXT,
	PRIMARY KEY (id, game_id),
	FOREIGN KEY(game_id) REFERENCES games(game_id)
);
`

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite3", "file:"+path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(pragmas); err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetGame(gameID int64, date string) (*models.Game, error) {
	game, err := models.Games(models.GameWhere.GameID.EQ(gameID), models.GameWhere.Date.EQ(date)).One(context.TODO(), r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return game, err
}

func (r *Repository) GetGamesMissingContent(incremental bool) ([]*models.Game, error) {
	mods := []qm.QueryMod{
		qm.Expr(
			models.GameWhere.Recap.IsNull(),
			qm.Or2(models.GameWhere.Extended.IsNull()),
		),
	}
	if incremental {
		cutoff := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
		mods = append(mods, models.GameWhere.Date.GTE(cutoff))
	}
	return models.Games(mods...).All(context.TODO(), r.db)
}

func (r *Repository) GetGames() ([]*models.Game, error) {
	return models.Games(qm.Load(models.GameRels.Highlights)).All(context.TODO(), r.db)
}

func (r *Repository) UpsertGame(game *models.Game) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := game.Upsert(context.TODO(), tx, true, []string{models.GameColumns.GameID}, boil.Infer(), boil.Infer()); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpsertHighlights(highlights []*models.Highlight) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, h := range highlights {
		if err := h.Upsert(context.TODO(), tx, true, []string{models.HighlightColumns.ID, models.HighlightColumns.GameID}, boil.Infer(), boil.Infer()); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
