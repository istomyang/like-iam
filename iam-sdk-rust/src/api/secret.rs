use super::{
    meta::{ListMeta, ObjectMeta},
    Db,
};

use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct Secret {
    #[serde(flatten)]
    meta: ObjectMeta,

    username: String,
    description: String,
    secret_id: String,
    secret_key: String,

    expires_at: i64,
}

impl Secret {
    #[allow(dead_code)]
    pub fn table_name() -> &'static str {
        "secret"
    }
}

impl Db for Secret {
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
pub struct SecretList {
    #[serde(flatten)]
    meta: ListMeta,

    items: Vec<Secret>,
}
