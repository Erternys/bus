use crate::{conn::Conn, str_vec};

use super::errors::{HttpError, HttpErrorKind};
#[derive(Default, Debug)]
pub struct Request {
  pub method: String,
  pub version: String,
  pub url: String,

  pub headers: super::Headers,
  pub body: Vec<u8>
}

impl Request {
  pub fn from_conn<'p>(conn: &mut Conn) -> Result<Self, HttpError<'p>> {
    let mut body = Vec::new();
    if let Err(_) = conn.read(&mut body) {
      return Err(HttpError::new(HttpErrorKind::Reading, "the data cannot be read"))
    }

    Self::parse(match String::from_utf8(body) {
      Ok(d) => d,
      Err(_) => return Err(HttpError::new(HttpErrorKind::Utf8, "the data has not been encoded in utf-8 OR the utf-8 parsing is not correct"))
    })
  }

  pub fn parse<'p>(req_str: String) -> Result<Self, HttpError<'p>> {
    let mut req = Self::default();

    let mut lines = req_str.split_inclusive("\r\n");
    match lines.next() {
      Some(line) => {
        let line = line.trim();
        if let [method, url, version] = Vec::from_iter(line.splitn(3, " ")).as_slice() {
          req.method = method.to_string();
          req.url = url.to_string();
          req.version = version.to_string();
        } else {
          return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
        }
      },
      None => ()
    }
    while let Some(line) = lines.next() {
      if let [header, value] = Vec::from_iter(line.splitn(2, ":")).as_slice() {
        req
          .headers
          .insert(header.trim().to_string(), value.trim().to_string());
      } else if line.len() == 2 {
        break
      } else {
        return Err(HttpError::new(HttpErrorKind::Parsing, "the data cannot be parsed"))
      }
    }

    while let Some(line) = lines.next() {
      req.body.append(&mut Vec::from(line))
    }

    Ok(req)
  }

  pub fn send(&mut self, conn: &mut Conn) -> Result<(), HttpError>{
    let mut data: Vec<u8> = Vec::new();
    data.append(&mut str_vec!("{} {} {}\r\n", self.method, self.url, self.version));
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
