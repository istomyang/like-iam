use crate::{
    api::login::{LoginResponse, UserNamePassword},
    iam::{chain::api_v1::request_map_err, config::Config},
};

use super::{api_v1::ApiClient, authz_v1::AuthzClient};

pub struct Client {
    config: Config,
}

impl Client {
    pub fn new(config: Config) -> Self {
        Self { config }
    }

    pub async fn login(
        &self,
        username: &str,
        password: &str,
    ) -> Result<LoginResponse, crate::error::Error> {
        let url = self.config.server_info.address.clone();
        let c = reqwest::Client::new();
        let res = c
            .post(url)
            .json(&UserNamePassword {
                username: username.to_string(),
                password: password.to_string(),
            })
            .send()
            .await
            .map_err(request_map_err)?;
        let l = res.json::<LoginResponse>().await.map_err(request_map_err)?;
        Ok(l)
    }

    pub async fn refresh_token(&self, token: String) -> Result<LoginResponse, crate::error::Error> {
        let url = self.config.server_info.address.clone();
        let c = reqwest::Client::new();
        let v = format!("Bearer {}", token);
        let res = c
            .post(url)
            .header("Authorization", v)
            .send()
            .await
            .map_err(request_map_err)?;
        let l = res.json::<LoginResponse>().await.map_err(request_map_err)?;
        Ok(l)
    }

    pub fn api(&self) -> ApiClient {
        #[cfg(feature = "iam_v1")]
        ApiClient::new(&self.config)
    }

    pub fn authz(&self) -> AuthzClient {
        #[cfg(feature = "iam_v1")]
        AuthzClient::new(&self.config)
    }
}

#[cfg(test)]
mod test {
    use crate::{
        api::{meta::CeateOperationMeta, user::User},
        iam::chain::interface::IUserClient,
    };

    use super::*;
    #[test]
    fn test_client() {
        let config = Config::load_from_file("path").unwrap();
        f(config);
    }

    async fn f(config: Config) {
        let user: User = Default::default();
        let opts: CeateOperationMeta = Default::default();

        let _ = Client::new(config).api().user().create(user, opts).await;
    }
}
