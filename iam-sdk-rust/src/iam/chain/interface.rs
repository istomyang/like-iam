use crate::{
    api::{
        authz::{RequestInfo, ResponseInfo},
        meta::{
            CeateOperationMeta, DeleteOperationMeta, GetOperationMeta, ListOperationMeta,
            UpdateOperationMeta,
        },
        policy::{Policy, PolicyList},
        secret::{Secret, SecretList},
        user::{User, UserList},
    },
    error::Error,
};

use async_trait::async_trait;

pub trait IApiClient<A, B, C>
where
    A: IUserClient,
    B: ISecretClient,
    C: IPolicyClient,
{
    fn user(&self) -> A;
    fn secret(&self) -> B;
    fn policy(&self) -> C;
}

#[async_trait]
pub trait IUserClient {
    async fn create(&self, user: User, opts: CeateOperationMeta) -> Result<User, Error>;
    async fn update(&self, user: User, opts: UpdateOperationMeta) -> Result<User, Error>;
    async fn delete(&self, username: String, opts: DeleteOperationMeta) -> Result<(), Error>;
    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), Error>;
    async fn get(&self, username: String, opts: GetOperationMeta) -> Result<User, Error>;
    async fn list(&self, opts: ListOperationMeta) -> Result<UserList, Error>;
}

#[async_trait]
pub trait ISecretClient {
    async fn create(&self, secret: Secret, opts: CeateOperationMeta) -> Result<Secret, Error>;
    async fn update(&self, secret: Secret, opts: UpdateOperationMeta) -> Result<Secret, Error>;
    async fn delete(&self, username: String, opts: DeleteOperationMeta) -> Result<(), Error>;
    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), Error>;
    async fn get(&self, username: String, opts: GetOperationMeta) -> Result<Secret, Error>;
    async fn list(&self, opts: ListOperationMeta) -> Result<SecretList, Error>;
}

#[async_trait]
pub trait IPolicyClient {
    async fn create(&self, policy: Policy, opts: CeateOperationMeta) -> Result<Policy, Error>;
    async fn update(&self, policy: Policy, opts: UpdateOperationMeta) -> Result<Policy, Error>;
    async fn delete(&self, username: String, opts: DeleteOperationMeta) -> Result<(), Error>;
    async fn delete_collection(
        &self,
        delete_opts: DeleteOperationMeta,
        list_opts: ListOperationMeta,
    ) -> Result<(), Error>;
    async fn get(&self, username: String, opts: GetOperationMeta) -> Result<Policy, Error>;
    async fn list(&self, opts: ListOperationMeta) -> Result<PolicyList, Error>;
}

#[async_trait]
pub trait IAuthzClient {
    async fn auth(&self, req: RequestInfo) -> Result<ResponseInfo, Error>;
}
