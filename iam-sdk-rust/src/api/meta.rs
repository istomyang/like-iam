#![allow(dead_code)]

use crate::api::format::rfc3339_utc;
use crate::error::{Error, ErrorKind};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::{Map, Value};

use super::Db;

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct ObjectMeta {
    id: u64,
    instance_id: String,
    name: String,
    #[serde(skip)]
    pub extend: Option<Map<String, Value>>,
    extend_shadow: String,

    #[serde(with = "rfc3339_utc")]
    create_at: DateTime<Utc>,
    #[serde(with = "rfc3339_utc")]
    update_at: DateTime<Utc>,
    #[serde(with = "rfc3339_utc")]
    delete_at: DateTime<Utc>,
}

impl Default for ObjectMeta {
    fn default() -> Self {
        ObjectMeta {
            id: 123,
            instance_id: "instance_id".to_string(),
            name: "name-123fsz".to_string(),
            extend: None,
            extend_shadow: "".to_string(),
            create_at: Utc::now(),
            update_at: Utc::now(),
            delete_at: Utc::now(),
        }
    }
}

impl Db for ObjectMeta {
    fn _before_create(&mut self) -> Result<(), Error> {
        let v = Value::Object(self.extend.clone().unwrap());
        self.extend_shadow = serde_json::to_string(&v)
            .map_err(|e| Error::new_with_str(ErrorKind::JsonEncode, e.to_string()))?;
        Ok(())
    }

    fn _before_update(&mut self) -> Result<(), Error> {
        let v = Value::Object(self.extend.clone().unwrap());
        self.extend_shadow = serde_json::to_string(&v)
            .map_err(|e| Error::new_with_str(ErrorKind::JsonEncode, e.to_string()))?;
        Ok(())
    }

    fn _after_find(&mut self) -> Result<(), Error> {
        let v: Value = serde_json::from_str(&self.extend_shadow)
            .map_err(|e| Error::new_with_str(ErrorKind::JsonDecode, e.to_string()))?;
        let ex = v.as_object().unwrap();
        self.extend = Some(ex.clone());
        Ok(())
    }
}

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct TypeMeta {
    // Kind is a string value representing the REST resource this object represents.
    // Servers may infer this from the endpoint the client submits requests to.
    // Cannot be updated.
    // In CamelCase.
    // required: false
    pub kind: String,

    // APIVersion defines the versioned schema of this representation of an object.
    // Servers should convert recognized schemas to the latest internal value, and
    // may reject unrecognized values.
    pub api_version: String,
}

impl Default for TypeMeta {
    fn default() -> Self {
        Self {
            kind: "test".to_string(),
            api_version: "v1+test".to_string(),
        }
    }
}

#[test]
fn test_object_meta_ser() {
    let json_str = r#"
      {
        "id": 12345,
        "instanceId": "iam-123jbz",
        "name": "name",
        "extendShadow": "{ \"a\": 1, \"b\": \"hello, world!\"}",
        "createAt": "1996-12-19T16:39:57-08:00",
        "updateAt": "1996-12-19T16:39:57-08:00",
        "deleteAt": "1996-12-19T16:39:57-08:00"
      }
    "#;

    let mut data: ObjectMeta = serde_json::from_str(json_str).unwrap();
    println!("{:#?}", data);

    let serialized = serde_json::to_string_pretty(&data).unwrap();
    println!("{}", serialized);

    let es = &data.extend_shadow;
    let s: Value = serde_json::from_str(es).unwrap();
    let o = s.as_object().unwrap();
    data.extend = Some(o.clone());

    print!("{:?}", o);
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ListMeta {
    pub total_count: u64,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct GetOperationMeta {
    meta: TypeMeta,
}

impl Default for GetOperationMeta {
    fn default() -> Self {
        Self {
            meta: Default::default(),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CeateOperationMeta {
    #[serde(flatten)]
    meta: TypeMeta,

    // When present, indicates that modifications should not be
    // persisted. An invalid or unrecognized dryRun directive will
    // result in an error response and no further processing of the
    // request. Valid values are:
    // - All: all dry run stages will be processed
    // +optional
    dray_run: Vec<String>,
}

impl Default for CeateOperationMeta {
    fn default() -> Self {
        Self {
            meta: Default::default(),
            dray_run: vec![],
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct UpdateOperationMeta {
    meta: TypeMeta,

    // When present, indicates that modifications should not be
    // persisted. An invalid or unrecognized dryRun directive will
    // result in an error response and no further processing of the
    // request. Valid values are:
    // - All: all dry run stages will be processed
    // +optional
    dray_run: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DeleteOperationMeta {
    meta: TypeMeta,

    // +optional
    unscoped: bool,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PatchOperationMeta {
    meta: TypeMeta,

    // When present, indicates that modifications should not be
    // persisted. An invalid or unrecognized dryRun directive will
    // result in an error response and no further processing of the
    // request. Valid values are:
    // - All: all dry run stages will be processed
    // +optional
    dray_run: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ListOperationMeta {
    meta: TypeMeta,

    // LabelSelector is used to find matching REST resources.
    label_selector: Option<String>,

    // FieldSelector restricts the list of returned objects by their fields. Defaults to everything.
    field_selector: Option<String>,

    // TimeoutSeconds specifies the seconds of ClientIP type session sticky time.
    timeout_seconds: Option<i64>,

    // Offset specify the number of records to skip before starting to return the records.
    offset: Option<i64>,

    // Limit specify the number of records to be retrieved.
    limit: Option<i64>,
}
