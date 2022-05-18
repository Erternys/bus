use std::collections::BTreeMap;

use super::errors::{ParserError, ParserErrorKind};

#[derive(Debug)]
pub struct Response {
  pub status: u16,
  pub version: String,
  pub message: String,

  pub headers: super::Headers,
  pub body: Vec<u8>
}

impl Response {
  pub fn parse(res_str: &str) -> Result<Self, ParserError> {
    let mut res = Self::default();

    let mut lines = res_str.split_inclusive("\r\n");
    match lines.next() {
      Some(line) => {
        let line = line.trim();
        if let [version, status, message] = Vec::from_iter(line.splitn(3, " ")).as_slice() {
          res.status = match status.parse() {
            Ok(v) => v,
            Err(_) => return Err(ParserError::new(ParserErrorKind::Parsing, "the data cannot be parsed"))
          };
          res.message = message.to_string();
          res.version = version.to_string();
        } else {
          return Err(ParserError::new(ParserErrorKind::Parsing, "the data cannot be parsed"))
        }
      },
      None => return Err(ParserError::new(ParserErrorKind::Parsing, "the data cannot be parsed"))
    }
    while let Some(line) = lines.next() {
      if let [header, value] = Vec::from_iter(line.splitn(2, ":")).as_slice() {
        res
          .headers
          .insert(header.trim().to_string(), value.trim().to_string());
      } else if line.len() == 2 {
        break
      } else {
        return Err(ParserError::new(ParserErrorKind::Parsing, "the data cannot be parsed"))
      }
    }

    while let Some(line) = lines.next() {
      res.body.append(&mut Vec::from(line))
    }

    Ok(res)
  }
}

impl Default for Response {
  fn default() -> Self {
    Self {
      status: 100,
      version: String::from("HTTP/1.1"),
      message: String::new(),
      headers: BTreeMap::new(),
      body: Vec::new()
    }
  }
}
