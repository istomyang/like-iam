// use std::time::Duration;

// use crate::error::{Error, ErrorKind, Result};
// use reqwest::header::HeaderMap;

// #[derive(Debug, Clone, Copy)]
// #[non_exhaustive]
// pub(crate) enum HttpMethod {
//     Get,
//     Put,
//     Post,
//     Delete,
// }

// #[derive(Debug, Clone, Copy)]
// #[non_exhaustive]
// pub(crate) enum Resource {
//     User,
//     Secret,
//     Policy,
// }

// impl std::fmt::Display for Resource {
//     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
//         match self {
//             User => write!(f, "users"),
//             Secret => write!(f, "secrets"),
//             Policy => write!(f, "policies"),
//         }
//     }
// }

// pub(self) fn map_err(e: reqwest::Error) -> Error {
//     Error::new_with_str(ErrorKind::Unknown, e.to_string())
// }

// /// Client is used to set config for http client, which can build many requests.
// pub(crate) struct Client {
//     inner: reqwest::ClientBuilder,
// }

// impl Client {
//     pub fn new() -> Self {
//         Self {
//             inner: reqwest::ClientBuilder::new(),
//         }
//     }

//     // pub fn get(&self, url: &str) -> Result<Request> {
//     //     let c = self.inner.build().map_err(map_err)?;
//     // }

//     pub fn user_agent(mut self, ua: &str) -> Self {
//         self.inner.user_agent(ua);
//         self
//     }

//     pub fn default_headers(mut self, headers: HeaderMap) -> Self {
//         self.inner.default_headers(headers);
//         self
//     }

//     pub fn timeout(mut self, timeout_ms: u64) -> Self {
//         self.inner.timeout(Duration::from_millis(timeout_ms));
//         self.inner
//             .connect_timeout(Duration::from_millis(timeout_ms));
//         self
//     }
// }

// #[test]
// fn test_rest() {
//     let header = HeaderMap::new();
//     let client = Client::new()
//         .user_agent("User-Agent")
//         .default_headers(header);
// }

// const DEFAULT_TIMEOUT_MS: u32 = 1_000;

// pub(crate) struct Request {
//     client: Client,
// }

// impl Request {}

// pub(crate) struct Response {}

// // pub(crate) struct Request {
// //     inner: Option<reqwest::RequestBuilder>,
// //     method: HttpMethod,
// //     resource: Option<Resource>,
// //     /// sub_resource is to set manually when like /login or /logout.
// //     sub_resource: Option<String>,
// //     timeout_ms: u32,
// // }

// // impl Request {
// //     fn send(self) -> crate::error::Result<Response> {
// //         Ok(Response {})
// //     }

// //     fn new(method: HttpMethod, url: &str) -> Self {
// //         Self {
// //             inner: None,
// //             method: method,
// //             resource: None,
// //             sub_resource: None,
// //             timeout_ms: DEFAULT_TIMEOUT_MS,
// //         }
// //     }

// //     fn sub_resource(mut self, resource: &str) -> Self {
// //         self.sub_resource = Some(resource.to_string());
// //         self
// //     }

// //     fn resource(mut self, resource: Resource) -> Self {
// //         self.resource = Some(resource);
// //         self
// //     }

// //     fn timeout_ms(mut self, ms: u32) -> Self {
// //         self.timeout_ms = ms;
// //         self
// //     }
// // }

// // pub(crate) struct Request<T>
// // where
// //     T: Serializer + ?Sized,
// // {
// //     /// url should without any query parameters and sub-paths.
// //     url: String,
// //     method: HttpMethod,
// //     resource: Option<Resource>,
// //     params: Option<HashMap<String, String>>,
// //     headers: Option<HashMap<String, String>>,
// //     body: Option<T>,
// //     timeout_ms: u32,

// //     inner: reqwest::RequestBuilder,
// // }

// // impl<T> Request<T>
// // where
// //     T: Serializer + ?Sized,
// // {
// //     fn new(url: String, method: HttpMethod) -> Self {
// //         Self {
// //             url,
// //             method,
// //             resource: None,
// //             params: None,
// //             headers: None,
// //             body: None,
// //             timeout_ms: DEFAULT_TIMEOUT_MS,
// //         }
// //     }

// //     fn header<'a>(mut self, key: &'a str, value: &'a str) -> Self {}

// //     fn timeout_ms(mut self, ms: u32) -> Self {
// //         self.timeout_ms = ms;
// //         self
// //     }

// //     fn resource(mut self, resource: Resource) -> Self {
// //         self.resource = Some(resource);
// //         self
// //     }

// //     fn body(mut self, data: T) -> Self {
// //         self.body = Some(data);
// //         self
// //     }
// // }
