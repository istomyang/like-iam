use super::{
    format::rfc3339_utc,
    meta::{ListMeta, ObjectMeta},
    Db,
};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct User {
    #[serde(flatten)]
    meta: ObjectMeta,

    username: String,
    password: String,
    is_admin: bool,

    #[serde(with = "rfc3339_utc")]
    login_at: DateTime<Utc>,
    total_policy: u64,
}

impl User {
    #[allow(dead_code)]
    pub fn table_name() -> &'static str {
        "user"
    }
}

impl Default for User {
    fn default() -> Self {
        User {
            meta: Default::default(),
            username: "admin".to_string(),
            password: "admin".to_string(),
            is_admin: true,
            login_at: Utc::now(),
            total_policy: 8_u64,
        }
    }
}

impl Db for User {
    fn _before_create(&mut self) -> Result<(), crate::error::Error> {
        Ok(())
    }

    fn _before_update(&mut self) -> Result<(), crate::error::Error> {
        Ok(())
    }

    fn _after_find(&mut self) -> Result<(), crate::error::Error> {
        Ok(())
    }
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct UserList {
    #[serde(flatten)]
    meta: ListMeta,

    items: Vec<User>,
}
