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
