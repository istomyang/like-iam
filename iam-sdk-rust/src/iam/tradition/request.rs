use std::collections::HashMap;

use crate::api::authz::AuthInfo;

pub struct AuthzRequest {
    pub auth_info: Option<AuthInfo>,
    pub url: String,
    pub method: String,
    pub headers: Option<HashMap<String, String>>,
    pub version: String,
}

impl AuthzRequest {
    
}
