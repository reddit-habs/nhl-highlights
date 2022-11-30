use fallible_iterator::FallibleIterator;
use rusqlite::OptionalExtension;

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Game {
    pub game_id: i64,
    pub date: String,
    pub type_: String,
    pub away: String,
    pub home: String,
    pub season: String,
    pub recap: Option<String>,
    pub extended: Option<String>,
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
    pub fn new(path: &str) -> anyhow::Result<DB> {
        let handle = rusqlite::Connection::open(path)?;
        handle.execute_batch(include_str!("schema.sql"))?;
        Ok(DB { handle })
    }

    pub fn get_game(&self, game_id: i64, date: &str) -> anyhow::Result<Option<Game>> {
        let mut stmt = self.handle.prepare(
            "SELECT game_id, date, type, away, home, season, recap, extended FROM games WHERE game_id = ? AND date = ?",
        )?;
        Ok(stmt.query_row((game_id, date), |r| r.try_into()).optional()?)
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

#[test]
fn test_game_cru() -> anyhow::Result<()> {
    let db = DB::new(":memory:")?;
    let g1 = &Game {
        game_id: 123,
        date: "2022-11-10".into(),
        type_: "PR".into(),
        away: "MTL".into(),
        home: "TOR".into(),
        season: "20212022".into(),
        recap: None,
        extended: None,
    };

    db.upsert_game(g1)?;
    let mut g2 = db.get_game(123, "2022-11-10")?.unwrap();
    assert_eq!(g1, &g2);

    g2.recap = Some("hello".into());
    db.upsert_game(&g2)?;
    let g3 = db.get_game(123, "2022-11-10")?.unwrap();
    assert_eq!(g2, g3);

    Ok(())
}
