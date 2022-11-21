use fallible_iterator::FallibleIterator;

pub struct Game {
    game_id: i64,
    date: String,
    type_: String,
    away: String,
    home: String,
    season: String,
    recap: Option<String>,
    extended: Option<String>,
}

impl<'a> TryFrom<&'a rusqlite::Row<'a>> for Game {
    type Error = rusqlite::Error;

    fn try_from(r: &rusqlite::Row) -> Result<Self, Self::Error> {
        Ok(Game {
            game_id: r.get(0)?,
            date: r.get(1)?,
            type_: r.get(2)?,
            away: r.get(3)?,
            home: r.get(4)?,
            season: r.get(5)?,
            recap: r.get(6)?,
            extended: r.get(7)?,
        })
    }
}

pub struct DB {
    handle: rusqlite::Connection,
}

impl DB {
    pub fn new() -> anyhow::Result<DB> {
        let handle = rusqlite::Connection::open("games.db")?;
        handle.execute_batch(include_str!("schema.sql"))?;
        Ok(DB { handle })
    }

    pub fn get_game(&self, game_id: i64, date: &str) -> anyhow::Result<Game> {
        let mut stmt = self.handle.prepare(
            "SELECT game_id, date, type, away, home, season, recap, extended FROM games WHERE game_id = ? AND date = ?",
        )?;
        Ok(stmt.query_row((game_id, date), |r| r.try_into())?)
    }

    pub fn get_games(&self) -> anyhow::Result<Vec<Game>> {
        let mut stmt = self
            .handle
            .prepare("SELECT game_id, date, type, away, home, season, recap, extended FROM games")?;
        let x = stmt.query(())?.map(|r| r.try_into()).collect()?;
        Ok(x)
    }

    pub fn get_games_missing_content(&self) -> anyhow::Result<Vec<Game>> {
        let mut stmt = self
            .handle
            .prepare("SELECT game_id, date, type, away, home, season, recap, extended FROM games WHERE recap IS NULL or extended IS NULL")?;
        let x = stmt.query(())?.map(|r| r.try_into()).collect()?;
        Ok(x)
    }

    pub fn upsert_game(&self, game: &Game) -> anyhow::Result<()> {
        self.handle.execute(
            r#"
            INSERT INTO games (
                game_id, date, type, away, home, season, recap, extended
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (game_id) DO UPDATE SET
                date = excluded.date,
                type = excluded.type,
                home = excluded.home,
                away = excluded.away,
                season = excluded.season,
                recap = excluded.recap,
                extended = excluded.extended;
            "#,
            (
                game.game_id,
                &game.date,
                &game.type_,
                &game.away,
                &game.home,
                &game.season,
                &game.recap,
                &game.extended,
            ),
        )?;
        Ok(())
    }
}

/*
// Game is an object representing the database table.
type Game struct {
    GameID   int64       `boil:"game_id" json:"game_id" toml:"game_id" yaml:"game_id"`
    Date     string      `boil:"date" json:"date" toml:"date" yaml:"date"`
    Type     string      `boil:"type" json:"type" toml:"type" yaml:"type"`
    Away     string      `boil:"away" json:"away" toml:"away" yaml:"away"`
    Home     string      `boil:"home" json:"home" toml:"home" yaml:"home"`
    Season   string      `boil:"season" json:"season" toml:"season" yaml:"season"`
    Recap    null.String `boil:"recap" json:"recap,omitempty" toml:"recap" yaml:"recap,omitempty"`
    Extended null.String `boil:"extended" json:"extended,omitempty" toml:"extended" yaml:"extended,omitempty"`

    R *gameR `boil:"-" json:"-" toml:"-" yaml:"-"`
    L gameL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}
 */
