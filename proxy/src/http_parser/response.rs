use std::collections::BTreeMap;

use crate::{conn::Conn, str_vec};

use super::errors::{HttpError, HttpErrorKind};

#[derive(Debug)]
pub struct Response {
  pub status: u16,
  pub version: String,
  pub message: String,

  pub headers: super::Headers,
  pub body: Vec<u8>
}

impl Response {
  pub fn parse<'p>(res_str: String) -> Result<Self, HttpError<'p>> {
    let mut res = Self::default();

    let mut lines = res_str.split_inclusive("\r\n");
    match lines.next() {
      Some(line) => {
        let line = line.trim();
        if let [version, status, message] = Vec::from_iter(line.splitn(3, " ")).as_slice() {
          res.status = match status.parse() {
            Ok(v) => v,
            Err(_) => return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
          };
          res.message = message.to_string();
          res.version = version.to_string();
        } else {
          return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
        }
      },
      None => return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
    }
    while let Some(line) = lines.next() {
      if let [header, value] = Vec::from_iter(line.splitn(2, ":")).as_slice() {
        res
          .headers
          .insert(header.trim().to_string(), value.trim().to_string());
      } else if line.len() == 2 {
        break
      } else {
        return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
      }
    }

    while let Some(line) = lines.next() {
      res.body.append(&mut Vec::from(line))
    }

    Ok(res)
  }

  pub fn vacuum(&mut self, conn: &mut Conn) -> Result<(), HttpError> {
    let mut body = Vec::new();
    if let Err(_) = conn.read(&mut body) {
      return Err(HttpError::new(HttpErrorKind::Reading, "the data cannot be read"))
    }

    *self = Self::parse(match String::from_utf8(body) {
      Ok(d) => d,
      Err(_) => return Err(HttpError::new(HttpErrorKind::Utf8, "the data has not been encoded in utf-8 OR the utf-8 parsing is not correct"))
    })?;

    Ok(())
  }

  pub fn send(&mut self, conn: &mut Conn) -> Result<(), HttpError>{
    let mut data: Vec<u8> = Vec::new();
    data.append(&mut str_vec!("{} {} {}\r\n", self.version, self.status, self.message));
    for (header, value) in &self.headers {
      let mut line = str_vec!("{}: {}\r\n", header, value);
      data.append(&mut line);
    }
    data.append(&mut str_vec!("\r\n"));
    data.append(&mut self.body);

    match conn.write(&data){
      Ok(_) => Ok(()),
      Err(_) => Err(HttpError::new(HttpErrorKind::Sending, "an error has occurred when the response was sent"))
    }
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
