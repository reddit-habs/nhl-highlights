package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sbstp/nhl-highlights/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const pragmas = `
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
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

CREATE TABLE IF NOT EXISTS cached_pages (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	season TEXT NOT NULL,
	team TEXT,
	content BLOB NOT NULL
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

func (r *Repository) Close() error {
	return r.db.Close()
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
	return models.Games().All(context.TODO(), r.db)
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

func (r *Repository) UpdateCachedPagesIteratively(ctx context.Context, stream chan *models.CachedPage) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := models.CachedPages().DeleteAll(context.TODO(), tx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			// error occured in producer, rollback
			return context.Cause(ctx)
		case c, ok := <-stream:
			if !ok {
				// iterator closed without error, commit
				return tx.Commit()
			}
			if err := c.Insert(context.TODO(), tx, boil.Infer()); err != nil {
				return err
			}
		}
	}
}

func (r *Repository) UpdateCachedPagges(cachedPages []*models.CachedPage) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := models.CachedPages().DeleteAll(context.TODO(), tx); err != nil {
		return err
	}

	for _, c := range cachedPages {
		if err := c.Insert(context.TODO(), tx, boil.Infer()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetCachedPage(season string, team *string) (*models.CachedPage, error) {
	return models.CachedPages(models.CachedPageWhere.Season.EQ(season), models.CachedPageWhere.Team.EQ(null.StringFromPtr(team))).One(context.TODO(), r.db)
}

func (r *Repository) GetCurrentSeason() (string, error) {
	var season string
	row := r.db.QueryRow("SELECT MAX(season) FROM games")
	if err := row.Scan(&season); err != nil {
		return "", err
	}
	return cleanupSeason(season), nil
}

func cleanupSeason(s string) string {
	if len(s) == 8 {
		return s[:4] + "-" + s[4:]
	}
	return s
}
