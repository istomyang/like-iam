pub mod authz;
mod format;
pub mod login;
pub mod meta;
pub mod policy;
pub mod secret;
pub mod user;

use crate::error::Error;

/// Db abstracts the hook for db operations.
pub(self) trait Db {
    fn _before_create(&mut self) -> Result<(), Error>;
    fn _before_update(&mut self) -> Result<(), Error>;
    fn _after_find(&mut self) -> Result<(), Error>;
}
