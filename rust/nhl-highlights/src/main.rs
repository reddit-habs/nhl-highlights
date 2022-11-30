use nhlapi::{ScheduleDate, ScheduleGame, TeamsCache};

use crate::nhlapi::teams;

mod nhlapi;
mod repo;

fn main() -> anyhow::Result<()> {
    simple_logger::SimpleLogger::new()
        .with_level(log::LevelFilter::Info)
        .env()
        .init()
        .unwrap();
    realMain("output", true, "", "")
}

fn realMain(output_dir: &str, incremental: bool, start_date: &str, end_date: &str) -> anyhow::Result<()> {
    let repo = repo::DB::new("games.db")?;
    let client = nhlapi::Client::new();
    let teams_cache = nhlapi::TeamsCache::create(&client)?;

    let schedule = client.schedule(start_date, end_date)?;

    for date in schedule.dates.iter() {
        log::info!("Date: {}", &date.date);
        for game in &date.games {
            log::info!(
                "Game: {} at {}",
                game.teams.away.team_id.name,
                game.teams.home.team_id.name
            );
            if repo.get_game(game.game_id, &date.date)?.is_none() {
                log::info!("Adding game {} on date {}", game.game_id, &date.date);
                repo.upsert_game(&game_from_schedule(&teams_cache, &date, &game))?;
            }
        }
    }
    Ok(())
}

fn game_from_schedule(teams_cache: &TeamsCache, date: &ScheduleDate, game: &ScheduleGame) -> repo::Game {
    repo::Game {
        game_id: game.game_id,
        date: date.date.clone(),
        type_: game.type_.clone(),
        away: teams_cache
            .get_by_id(game.teams.away.team_id.id)
            .expect("bad away team")
            .abbrev
            .clone(),
        home: teams_cache
            .get_by_id(game.teams.home.team_id.id)
            .expect("bad home team")
            .abbrev
            .clone(),
        season: game.season.clone(),
        recap: None,
        extended: None,
    }
}
