use super::{
    meta::{ListMeta, ObjectMeta},
    Db,
};

use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Default)]
#[serde(rename_all = "camelCase")]
pub struct PolicyInfo {
    id: String,
    description: String,
    subjects: Vec<String>,
    effect: String,
    resources: Vec<String>,
    actions: Vec<String>,
    conditions: Vec<String>,
    meta: Vec<u8>,
}

impl PolicyInfo {
    #[allow(dead_code)]
    fn to_string(&self) -> String {
        serde_json::to_string(&self).unwrap()
    }
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct Policy {
    #[serde(flatten)]
    meta: ObjectMeta,

    username: String,

    #[serde(skip)]
    pub policy: PolicyInfo,

    policy_shadow: String,
}

impl Policy {
    #[allow(dead_code)]
    pub fn table_name() -> &'static str {
        "policy"
    }
}

impl Db for Policy {
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
pub struct PolicyList {
    #[serde(flatten)]
    meta: ListMeta,

    items: Vec<Policy>,
}
