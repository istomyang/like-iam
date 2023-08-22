use std::time::Duration;

use async_trait::async_trait;

use crate::{
    api::{
        meta::{
            CeateOperationMeta, DeleteOperationMeta, GetOperationMeta, ListOperationMeta,
            UpdateOperationMeta,
        },
        policy::{Policy, PolicyList},
        secret::{Secret, SecretList},
        user::{User, UserList},
    },
    error::{Error, ErrorKind},
    iam::config::Config,
};

use super::interface::{IPolicyClient, ISecretClient, IUserClient};

// pub struct ApiV1Client<A, B, C>
// where
//     A: IUserClient,
//     B: ISecretClient,
//     C: IPolicyClient,
// {
//     config: Config,
//     __marker: std::marker::PhantomData<(A, B, C)>,
// }
//
// impl<A, B, C> ApiV1Client<A, B, C>
// where
//     A: IUserClient,
//     B: ISecretClient,
//     C: IPolicyClient,
// {
//     pub fn new(config: Config) -> Self {
//         Self {
//             config,
//             __marker: std::marker::PhantomData,
//         }
//     }

//     fn create_client(&self) -> reqwest::Client {
//         let mut c = reqwest::ClientBuilder::new();
//         c = c.timeout(Duration::from_secs(
//             self.config.server_info.timeout_ms.into(),
//         ));
//         c = c.connect_timeout(Duration::from_secs(
//             self.config.server_info.timeout_ms.into(),
//         ));
//         c.build().expect("client-builder should has no problem.")
//     }

//     fn get_url(&self) -> String {
//         self.config.server_info.address.clone()
//     }
// }

// impl<A, B, C> IApiClient<A, B, C> for ApiV1Client<A, B, C>
// where
//     A: IUserClient,
//     B: ISecretClient,
//     C: IPolicyClient,
// {
//     fn user(&self) -> A {
//         let url = self.get_url();
//         let client = self.create_client();
//         UserClient::new(url, client)
//     }

//     fn secret(&self) -> B {
//         todo!()
//     }

//     fn policy(&self) -> C {
//         todo!()
//     }
// }

pub struct ApiClient {
    url: String,
    client: reqwest::Client,
    token: String,
}

impl ApiClient {
    pub fn new(config: &Config) -> Self {
        let mut builder = reqwest::ClientBuilder::new();
        builder = builder.timeout(Duration::from_secs(config.server_info.timeout_ms.into()));
        builder =
            builder.connect_timeout(Duration::from_secs(config.server_info.timeout_ms.into()));
        let client = builder
            .build()
            .expect("client-builder should has no problem.");
        let url = config.server_info.address.clone();
        Self {
            client,
            url,
            token: config.auth_info.token.clone(),
        }
    }

    pub fn user(&self) -> impl IUserClient {
        UserClient::new(self.url.clone(), self.client.clone(), self.token.clone())
    }

    pub fn secret(&self) -> impl ISecretClient {
        SecretClient::new(self.url.clone(), self.client.clone(), self.token.clone())
    }

    pub fn policy(&self) -> impl IPolicyClient {
        PolicyClient::new(self.url.clone(), self.client.clone(), self.token.clone())
    }
}

pub fn request_map_err(e: reqwest::Error) -> Error {
    Error::new_with_str(ErrorKind::Unknown, e.to_string())
}

pub fn jwt_header(token: String) -> (String, String) {
    ("Authorization".to_string(), format!("Bearer {}", token))
}

struct UserClient {
    client: reqwest::Client,
    url: String,
    token: String,
}

impl UserClient {
    pub fn new(root_url: String, client: reqwest::Client, token: String) -> Self {
        Self {
            client,
            url: format!("{}/users", root_url),
            token,
        }
    }
}

#[async_trait]
impl IUserClient for UserClient {
    async fn create(
        &self,
        user: User,
        opts: CeateOperationMeta,
    ) -> Result<User, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .post(&self.url)
            .json(&user)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(user)
    }

    async fn update(
        &self,
        user: User,
        opts: UpdateOperationMeta,
    ) -> Result<User, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .put(&self.url)
            .json(&user)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        Ok(user)
    }

    async fn delete(
        &self,
        username: String,
        opts: DeleteOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        Ok(())
    }

    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(&self.url)
            .query(&delete_opts)
            .query(&list_opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        Ok(())
    }

    async fn get(
        &self,
        username: String,
        opts: GetOperationMeta,
    ) -> Result<User, crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let u = res.json::<User>().await.map_err(request_map_err)?;
        Ok(u)
    }

    async fn list(
        &self,
        opts: ListOperationMeta,
    ) -> Result<crate::api::user::UserList, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&self.url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let u = res.json::<UserList>().await.map_err(request_map_err)?;
        Ok(u)
    }
}

struct SecretClient {
    client: reqwest::Client,
    url: String,
    token: String,
}

impl SecretClient {
    fn new(root_url: String, client: reqwest::Client, token: String) -> Self {
        Self {
            client,
            url: format!("{}/secret", root_url),
            token,
        }
    }
}

#[async_trait]
impl ISecretClient for SecretClient {
    async fn create(
        &self,
        secret: crate::api::secret::Secret,
        opts: CeateOperationMeta,
    ) -> Result<crate::api::secret::Secret, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .post(&self.url)
            .json(&secret)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(secret)
    }

    async fn update(
        &self,
        secret: crate::api::secret::Secret,
        opts: UpdateOperationMeta,
    ) -> Result<crate::api::secret::Secret, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .put(&self.url)
            .json(&secret)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(secret)
    }

    async fn delete(
        &self,
        username: String,
        opts: DeleteOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(&url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(())
    }

    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(&self.url)
            .query(&delete_opts)
            .query(&list_opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(())
    }

    async fn get(
        &self,
        username: String,
        opts: GetOperationMeta,
    ) -> Result<crate::api::secret::Secret, crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let s = res.json::<Secret>().await.map_err(request_map_err)?;
        Ok(s)
    }

    async fn list(
        &self,
        opts: ListOperationMeta,
    ) -> Result<crate::api::secret::SecretList, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&self.url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let s = res.json::<SecretList>().await.map_err(request_map_err)?;
        Ok(s)
    }
}

struct PolicyClient {
    client: reqwest::Client,
    url: String,
    token: String,
}

impl PolicyClient {
    fn new(root_url: String, client: reqwest::Client, token: String) -> Self {
        Self {
            client,
            url: format!("{}/policy", root_url),
            token,
        }
    }
}

#[async_trait]
impl IPolicyClient for PolicyClient {
    async fn create(
        &self,
        policy: Policy,
        opts: CeateOperationMeta,
    ) -> Result<Policy, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .post(&self.url)
            .json(&policy)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(policy)
    }

    async fn update(
        &self,
        policy: Policy,
        opts: UpdateOperationMeta,
    ) -> Result<Policy, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .put(&self.url)
            .json(&policy)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(policy)
    }

    async fn delete(
        &self,
        username: String,
        opts: DeleteOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(&url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(())
    }

    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        self.client
            .delete(&self.url)
            .query(&delete_opts)
            .query(&list_opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;

        Ok(())
    }

    async fn get(
        &self,
        username: String,
        opts: GetOperationMeta,
    ) -> Result<Policy, crate::error::Error> {
        let url = format!("{}/{}", self.url, username);
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let s = res.json::<Policy>().await.map_err(request_map_err)?;
        Ok(s)
    }

    async fn list(&self, opts: ListOperationMeta) -> Result<PolicyList, crate::error::Error> {
        let (k, v) = jwt_header(self.token.clone());
        let res = self
            .client
            .get(&self.url)
            .query(&opts)
            .header(k, v)
            .send()
            .await
            .map_err(request_map_err)?;
        let s = res.json::<PolicyList>().await.map_err(request_map_err)?;
        Ok(s)
    }
}

#[test]
fn test1() {
    let s = "https://www.google.com".to_string();

    println!("{}", &s[..5]);
}
