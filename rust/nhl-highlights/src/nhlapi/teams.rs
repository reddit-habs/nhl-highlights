use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct TeamsResponse {
    pub teams: Vec<Team>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Team {
    pub id: i32,
    pub name: String,
    #[serde(rename = "abbreviation")]
    pub abbrev: String,
    #[serde(rename = "teamName")]
    pub team_name: String,
    #[serde(rename = "locationName")]
    pub location: String,
}
