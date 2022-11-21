mod nhlapi;

fn main() -> anyhow::Result<()> {
    let client = nhlapi::Client::new();
    let teams_cache = nhlapi::TeamsCache::create(&client)?;

    println!("{:#?}", teams_cache.teams());

    Ok(())
}
