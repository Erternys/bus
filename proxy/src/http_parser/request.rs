use std::io::{Error, ErrorKind, Result};

#[derive(Default, Debug)]
pub struct Request {
  pub method: String,
  pub version: String,
  pub url: String,

  pub headers: super::Headers,
  pub body: Vec<u8>
}

impl Request {
  pub fn parse(req_str: &str) -> Result<Self> {
    let mut req = Self::default();
    println!("{:?}", req_str);

    let mut lines = req_str.split_inclusive("\r\n");
    match lines.next() {
      Some(line) => {
        let line = line.trim();
        if let [method, url, version] = Vec::from_iter(line.splitn(3, " ")).as_slice() {
          req.method = method.to_string();
          req.url = url.to_string();
          req.version = version.to_string();
        } else {
          return Err(Error::from(ErrorKind::InvalidData))
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
        return Err(Error::from(ErrorKind::InvalidData))
      }
    }

    while let Some(line) = lines.next() {
      req.body.append(&mut Vec::from(line))
    }

    Ok(req)
  }
}
