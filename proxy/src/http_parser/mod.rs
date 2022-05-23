mod errors;
mod request;
mod response;

use std::collections::BTreeMap;

pub use errors::HttpErrorKind;
pub use request::Request;
pub use response::Response;

pub type Headers = BTreeMap<String, String>;

#[macro_export]
macro_rules! str_vec {
    ($($arg:tt)*) => {{
        let res = Vec::from(format!($($arg)*));
        res
    }}
}
