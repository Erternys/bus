use std::{collections::BTreeMap, fs::{File, self}, io::Read};

use crate::{http_parser::{Response, Request}, conn::Conn};

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
    for (key, alias) in iter{
      if key == joker_key {
        joker = Some(Self::from_str(alias))
      }else{
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
    } else if s.len() == 0{
      Self::Empty
    } else {
      match fs::metadata(&s){
        Ok(metadata) =>{
          if metadata.is_file() {
            Self::File(s)
          } else {              
            Self::Folder(s)
          }
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

  pub fn to_response(&self, addr: String, path: String, req: &mut Request) -> Response {
    match self {
      Self::Folder(dest) => {
        if path.contains("./") {
          e406()
        } else {
          let path = req.url.replacen(&path, &dest, 1);
          match fs::metadata(&path){
            Ok(metadata) => {
              let file = if metadata.is_file() {
                Self::File(path)
              } else {
                println!("{}", path.clone() + "/index.html");          
                Self::File(path + "/index.html")
              };
              file.to_response(
                addr, 
                String::new(), 
                req
              )
            },
            Err(_) => return e404()
          }
        }
      },
      Self::File(path) => {
        let mut file = match File::open(path) {
          Ok(file) => file,
          Err(_) => return e404()
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
      Self::Http(path) => {
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
        req.url = req.url.replacen(&path, &path, 1).replace("//", "/");
        
        let mut res = Response::default();
        let mut server_conn = Conn::new(&host).unwrap();
        req.send(&mut server_conn).unwrap();
        res.vacuum(&mut server_conn).unwrap();
        res
      },
      Self::Empty => {
        e500()
      }
    }
  }
}
