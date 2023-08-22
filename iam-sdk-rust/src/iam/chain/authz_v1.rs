use std::time::Duration;

use crate::{
    api::authz::{RequestInfo, ResponseInfo},
    error::Error,
    iam::config::Config,
};

use super::api_v1::request_map_err;

pub struct AuthzClient {
    client: reqwest::Client,
    url: String,
}

impl AuthzClient {
    pub fn new(config: &Config) -> Self {
        let mut builder = reqwest::ClientBuilder::new();
        builder = builder.timeout(Duration::from_secs(config.server_info.timeout_ms.into()));
        builder =
            builder.connect_timeout(Duration::from_secs(config.server_info.timeout_ms.into()));
        let client = builder
            .build()
            .expect("client-builder should has no problem.");
        let url = format!("{}/authz", config.server_info.address.clone());
        Self { client, url }
    }

    async fn auth(&self, req: RequestInfo) -> Result<ResponseInfo, Error> {
        let res = self
            .client
            .post(&self.url)
            .json(&req)
            .send()
            .await
            .map_err(request_map_err)?;
        let r = res.json::<ResponseInfo>().await.map_err(request_map_err)?;
        Ok(r)
    }
}
