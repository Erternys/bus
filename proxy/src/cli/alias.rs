use std::{collections::BTreeMap, fs::File, io::Read};

use crate::{http_parser::{Response, Request}, conn::Conn};

use super::errors::{e440, e500};

#[derive(Debug)]
pub enum Alias {
  File(String),
  Http(String),
  Empty
}

impl Alias {
  pub fn from_iter(iter: BTreeMap<String, String>) -> (BTreeMap<String, Alias>, Option<Alias>) {
    let mut map = BTreeMap::new();
    let mut joker = None;
    let joker_key = "*".to_string();
    for (key, alias) in iter{
      if key == joker_key {
          joker = Some(Self::from_str(alias))
      }else{
        map.insert(key, Self::from_str(alias));
      }
    }
    (map, joker)
  }
  pub fn from_str(s: String) -> Alias {
    let splitted = s.splitn(2, ":").collect::<Vec<&str>>();
    if let [schema, value] = *splitted.as_slice() {
      match schema {
        "http" => Alias::Http(value.to_string()),
        "file" => Alias::File(value.to_string()),
        _ => Alias::Empty
      }
    } else if s.len() == 0{
      Alias::Empty
    } else if s.chars().next().unwrap() == '.'{
      Alias::File(s)
    }else{
      Alias::Http(s)
    }
  }

  pub fn to_response(&self, addr: String, key: String, req: &mut Request) -> Response {
    match self {
      Alias::File(path) => {
        let mut file = match File::open(path) {
          Ok(file) => file,
          Err(_) => return e440()
        };
        
        let mut buf = Vec::new();
        let size = match file.read_to_end(&mut buf){
          Ok(size) => size,
          Err(_) => return e500()
        };
        
        let mut res = Response::default();
        if size == 0 {
          res.status = 204;
          res.message = String::from("No Content");
        } else {
          res.status = 200;
          res.message = String::from("OK");
          res.headers.insert(String::from("Content-Length"), size.to_string());
        }

        res.body = buf;
        res
      },
      Alias::Http(path) => {
        let root = "/".to_string();
        let (host, path) = match path.split_once("/") {
          Some((host, path)) => (host.to_string(), root + path),
          None => {
            let path = path.clone();
            if path.chars().next().unwrap() == '/' {
              (addr, path)
            } else {
              (path, root)
            }
          }
        };
        
        req.headers.insert(String::from("Host"), host.clone());
        req.url = req.url.replacen(&key, &path, 1).replace("//", "/");
        
        let mut res = Response::default();
        let mut server_conn = Conn::new(&host).unwrap();
        req.send(&mut server_conn).unwrap();
        res.vacuum(&mut server_conn).unwrap();
        res
      },
      Alias::Empty => {
        e500()
      }
    }
  }
}
