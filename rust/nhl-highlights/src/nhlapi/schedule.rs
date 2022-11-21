use serde::{Deserialize, Serialize};

pub const GAME_TYPE_PRESEASON: &str = "PR";
pub const GAME_TYPE_REGULAR: &str = "R";
pub const GAME_TYPE_PLAYOFFS: &str = "P";

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleResponse {
    pub dates: Vec<ScheduleDate>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleDate {
    pub date: String,
    pub games: Vec<ScheduleGame>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleGame {
    #[serde(rename = "gamePk")]
    pub game_id: i64,
    #[serde(rename = "gameType")]
    pub type_: String,
    pub season: String,
    pub teams: ScheduleTeams,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleTeams {
    pub away: ScheduleTeam,
    pub home: ScheduleTeam,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleTeam {
    #[serde(rename = "team")]
    pub team_id: ScheduleTeamID,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ScheduleTeamID {
    pub id: i32,
    pub name: String,
}
