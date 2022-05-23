use std::collections::BTreeMap;
use std::fs::{self, File};
use std::io::Read;

use crate::conn::Conn;
use crate::http_parser::{Request, Response};

use super::errors::{e404, e406, e500};

#[derive(Debug)]
pub enum Alias {
  Folder(String),
  File(String),
  Http(String),
  Empty
}

impl Alias {
  pub fn from_iter(iter: BTreeMap<String, String>) -> (BTreeMap<String, Self>, Option<Self>) {
    let mut map = BTreeMap::new();
    let mut joker = None;
    let joker_key = "*".to_string();
    for (key, alias) in iter {
      if key == joker_key {
        joker = Some(Self::from_str(alias))
      } else {
        map.insert(key, Self::from_str(alias));
      }
    }
    (map, joker)
  }

  pub fn from_str(s: String) -> Self {
    let splitted = s.splitn(2, ":").collect::<Vec<&str>>();
    if let [schema, value] = *splitted.as_slice() {
      match schema {
        "http" => Self::Http(value.to_string()),
        "file" => Self::File(value.to_string()),
        "folder" => Self::Folder(value.to_string()),
        _ => Self::Empty
      }
    } else if s.len() == 0 {
      Self::Empty
    } else {
      match fs::metadata(&s) {
        Ok(metadata) =>
          if metadata.is_file() {
            Self::File(s)
          } else {
            Self::Folder(s)
          },
        Err(_) => {
          let first_char = s.chars().next().unwrap();
          if first_char == '.' || first_char == '/' {
            Self::Empty
          } else {
            Self::Http(s)
          }
        }
      }
    }
  }

  pub fn to_response(&self, addr: String, req_path: String, req: &mut Request) -> Response {
    match self {
      Self::Folder(dest) =>
        if req_path.contains("./") {
          e406()
        } else {
          let path = req.url.replacen(&req_path, &dest, 1);
          match fs::metadata(&path) {
            Ok(metadata) => {
              let file = if metadata.is_file() {
                Self::File(path)
              } else {
                Self::File(path + "/index.html")
              };
              file.to_response(addr, String::new(), req)
            },
            Err(_) => return e404()
          }
        },
      Self::File(dest) => {
        let mut file = match File::open(dest) {
          Ok(file) => file,
          Err(_) => return e404()
        };

        let mut buf = Vec::new();
        let size = match file.read_to_end(&mut buf) {
          Ok(size) => size,
          Err(_) => return e500()
        };

        let mut res = Response::default();
        if size == 0 {
          res.status = 204;
          res.message = "No Content".to_string();
        } else {
          res.status = 200;
          res.message = "OK".to_string();
          res
            .headers
            .insert("Content-Length".to_string(), size.to_string());
        }

        res.body = buf;
        res
      },
      Self::Http(urn) => {
        let root = "/".to_string();
        let (host, path) = match urn.split_once("/") {
          Some((host, path)) => (host.to_string(), root + path),
          None => {
            let path = urn.clone();
            if path.chars().next().unwrap() == '/' {
              (addr.clone(), path)
            } else {
              (path, root)
            }
          }
        };

        let host_str = "Host".to_string();
        let req_host = req
          .headers
          .insert(host_str.clone(), host.clone())
          .unwrap_or(addr);

        req.url = req.url.replacen(&req_path, &path, 1).replace("//", "/");

        let mut res = Response::default();
        let mut server_conn = Conn::new(&host).unwrap();
        req.send(&mut server_conn).unwrap();
        res.vacuum(&mut server_conn).unwrap();
        res.headers.insert(host_str.clone(), req_host);
        res
      },
      Self::Empty => e500()
    }
  }
}
