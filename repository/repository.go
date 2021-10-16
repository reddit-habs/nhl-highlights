package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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
`

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite3", "file:"+path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) scanGames(rows *sql.Rows) ([]*Game, error) {
	defer rows.Close()
	games := make([]*Game, 0)

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(
			&game.GameID,
			&game.Date,
			&game.Type,
			&game.Away,
			&game.Home,
			&game.Season,
			&game.Recap,
			&game.Extended,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func (r *Repository) GetGame(gameID int64) (*Game, error) {
	rows, err := r.db.Query("SELECT game_id, date, type, away, home, season, recap, extended FROM games WHERE game_id = ?", gameID)
	if err != nil {
		return nil, err
	}

	games, err := r.scanGames(rows)
	if err != nil {
		return nil, err
	}

	var game *Game
	if len(games) > 0 {
		game = games[0]
	}

	return game, nil
}

func (r *Repository) GetGamesMissingContent() ([]*Game, error) {
	rows, err := r.db.Query("SELECT game_id, date, type, away, home, season, recap, extended FROM games WHERE recap IS NULL or extended IS NULL")
	if err != nil {
		return nil, err
	}

	return r.scanGames(rows)
}

func (r *Repository) GetGames() ([]*Game, error) {
	rows, err := r.db.Query("SELECT game_id, date, type, away, home, season, recap, extended FROM games ORDER BY date DESC")
	if err != nil {
		return nil, err
	}

	return r.scanGames(rows)
}

func (r *Repository) UpsertGame(game *Game) error {
	result, err := r.db.Exec(
		`INSERT INTO games (game_id, date, type, away, home, season, recap, extended)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(game_id) DO UPDATE SET
			date=excluded.date,
			type=excluded.type,
			away=excluded.away,
			home=excluded.home,
			season=excluded.season,
			recap=excluded.recap,
			extended=excluded.extended`,
		game.GameID, game.Date, game.Type, game.Away, game.Home,
		game.Season, game.Recap, game.Extended,
	)
	if err != nil {
		return err
	}
	aff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if aff != 1 {
		return fmt.Errorf("no rows modified")
	}
	return nil
}
