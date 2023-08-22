use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

use super::format::rfc3339_utc;

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct LoginResponse {
    pub token: String,
    #[serde(with = "rfc3339_utc")]
    pub expire: DateTime<Utc>,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct UserNamePassword {
    pub username: String,
    pub password: String,
}
