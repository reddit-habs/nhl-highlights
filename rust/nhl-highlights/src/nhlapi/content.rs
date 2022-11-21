use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ContentResponse {
    pub media: ContentMedia,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ContentMedia {
    pub epg: Vec<ContentEPG>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ContentEPG {
    title: String,
    items: Vec<ContentItem>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ContentItem {
    #[serde(rename = "type")]
    type_: String,
    playbacks: Vec<ContentPlayback>,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct ContentPlayback {
    name: String,
    url: String,
}
