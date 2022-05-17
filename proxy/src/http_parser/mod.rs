mod request;
mod response;

use std::collections::BTreeMap;

pub use request::Request;
pub use response::Response;

pub type Headers = BTreeMap<String, String>;
