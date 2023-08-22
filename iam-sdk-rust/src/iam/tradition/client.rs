use crate::{
    api::authz::{RequestInfo, ResponseInfo},
    error::Error,
};

#[derive(Debug)]
pub struct IAMClient<'a> {
    id: &'a str,
    key: &'a str,
}

impl<'a> IAMClient<'a> {
    pub fn from_secret(id: &'a str, key: &'a str) -> Self {
        Self { id, key }
    }

    pub fn authz(req: RequestInfo) -> Result<ResponseInfo, Error> {
        todo!()
        // Err(Default::default())
    }
}

#[test]
fn test_tradition_iam_client() {
    let c = IAMClient::from_secret("id", "key");
    println!("{:?}", c);
}
