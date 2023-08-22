#![allow(dead_code)]

use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;

use crate::error::{Error, ErrorKind};

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct RequestInfo {
    #[serde(flatten)]
    info: AuthInfo,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct AuthInfo {
    pub resource: Option<String>,
    pub action: Option<String>,
    pub subject: Option<String>,
    pub context: Option<HashMap<String, Value>>,
}

impl RequestInfo {
    pub fn new(
        resource: String,
        action: String,
        subject: String,
        context: HashMap<String, Value>,
    ) -> Self {
        Self {
            info: AuthInfo {
                resource: Some(resource),
                action: Some(action),
                subject: Some(subject),
                context: Some(context),
            },
        }
    }

    pub fn to_json(&self) -> String {
        serde_json::to_string(&self).unwrap()
    }
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct ResponseInfo {
    allowed: bool,
    denied: bool,
    reason: String,
    error: String,
}

impl ResponseInfo {
    fn from(data: String) -> Result<Self, Error> {
        let a: Self = serde_json::from_str(&data)
            .map_err(|e| Error::new_with_str(ErrorKind::JsonDecode, e.to_string()))?;
        Ok(a)
    }
}
