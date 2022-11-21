use std::collections::{BTreeMap, HashMap};

use attohttpc::Session;

use super::*;

pub struct Client {
    sess: Session,
}

impl Client {
    pub fn new() -> Self {
        Client { sess: Session::new() }
    }

    pub fn schedule(&self, start_date: &str, end_date: &str) -> anyhow::Result<ScheduleResponse> {
        let resp = self
            .sess
            .get(format!(
                "https://statsapi.web.nhl.com/api/v1/schedule?startDate={}&endDate={}",
                start_date, end_date,
            ))
            .send()?;
        assert!(resp.is_success());
        Ok(resp.json_utf8()?)
    }

    pub fn content(&self, game_id: i64) -> anyhow::Result<ContentResponse> {
        let resp = self
            .sess
            .get(format!("https://statsapi.web.nhl.com/api/v1/game/{}/content", game_id,))
            .send()?;
        assert!(resp.is_success());
        Ok(resp.json_utf8()?)
    }

    pub fn teams(&self) -> anyhow::Result<TeamsResponse> {
        let resp = self.sess.get("https://statsapi.web.nhl.com/api/v1/teams").send()?;
        assert!(resp.is_success());
        Ok(resp.json_utf8()?)
    }
}

pub struct TeamsCache {
    teams: Vec<Team>,
    by_id: BTreeMap<i32, usize>,
    by_abbrev: HashMap<String, usize>,
}

impl TeamsCache {
    pub fn create(client: &Client) -> anyhow::Result<Self> {
        let teams = client.teams()?;
        let mut by_id = BTreeMap::new();
        let mut by_abbrev = HashMap::new();

        for (idx, team) in teams.teams.iter().enumerate() {
            by_id.insert(team.id, idx);
            by_abbrev.insert(team.abbrev.clone(), idx);
        }

        Ok(TeamsCache {
            teams: teams.teams,
            by_id,
            by_abbrev,
        })
    }

    pub fn teams(&self) -> &[Team] {
        &self.teams
    }

    pub fn get_by_id(&self, id: i32) -> Option<&Team> {
        self.by_id.get(&id).map(|&idx| &self.teams[idx])
    }

    pub fn get_by_abbrev(&self, abbrev: &str) -> Option<&Team> {
        self.by_abbrev.get(abbrev).map(|&idx| &self.teams[idx])
    }
}
