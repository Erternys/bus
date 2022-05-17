use std::io::{Error, ErrorKind, Result};

#[derive(Default, Debug)]
pub struct Response {
  pub status: String,
  pub version: String,
  pub message: String,

  pub headers: super::Headers,
  pub body: Vec<u8>
}

impl Response {
  pub fn parse(res_str: &str) -> Result<Self> {
    let mut res = Self::default();
    println!("{:?}", res_str);

    let mut lines = res_str.split_inclusive("\r\n");
    match lines.next() {
      Some(line) => {
        let line = line.trim();
        if let [version, status, message] = Vec::from_iter(line.splitn(3, " ")).as_slice() {
          res.status = status.to_string();
          res.message = message.to_string();
          res.version = version.to_string();
        } else {
          return Err(Error::from(ErrorKind::InvalidData))
        }
      },
      None => return Err(Error::from(ErrorKind::InvalidData))
    }
    while let Some(line) = lines.next() {
      if let [header, value] = Vec::from_iter(line.splitn(2, ":")).as_slice() {
        res
          .headers
          .insert(header.trim().to_string(), value.trim().to_string());
      } else if line.len() == 2 {
        break
      } else {
        return Err(Error::from(ErrorKind::InvalidData))
      }
    }

    while let Some(line) = lines.next() {
      res.body.append(&mut Vec::from(line))
    }

    Ok(res)
  }
}
