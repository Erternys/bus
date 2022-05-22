use std::fs::File;
use std::process::exit;

use clap::Args;

use crate::config::Config;
use crate::http_parser::Request;
use crate::server::Server;

use super::alias::Alias;

#[derive(Debug, Args)]
pub struct Start;

impl Start {
  pub fn call(&self) {
    let file = match File::open("./.bus.yaml") {
      Ok(f) => f,
      Err(_) => {
        eprintln!("The config file does not exist, please execute the command \"{0} init\" then \"{0} proxy init\" before", "bus");
        exit(1);
      }
    };

    let config: Config = match serde_yaml::from_reader(&file) {
      Ok(c) => c,
      Err(_) => {
        eprintln!("The config file cannot be read, please change the permission");
        exit(1);
      }
    };

    let proxy = if let Some(proxy) = config.proxy {
      proxy
    } else {
      eprintln!(
        "The proxy config does not exist, please execute the command \"{0} proxy init\" before",
        "bus"
      );
      exit(1)
    };

    let addr = format!(
      "{}:{}",
      if proxy.open.unwrap_or(true) {
        "127.0.0.1"
      } else {
        "localhost"
      },
      proxy.port
    );

    let (aliases, joker) = Alias::from_iter(proxy.aliases);
    
    let server = Server::new(addr);
    server.run(move |mut client_conn, addr| {
      let req = Request::from_conn(&mut client_conn);
      
      let mut req = match req {
        Ok(r) => r,
        Err(_) => {
          let mut res = super::errors::e400();
          res.send(&mut client_conn).unwrap();
          return
        }
      };

      let mut iter = aliases.iter();

      let res = loop {
        let (path, alias) = match iter.next(){
          Some((path, alias)) => (path, alias),
          None => break None
        };

        if req.url.starts_with(path) {
          break Some(alias.to_response(addr.clone(), path.clone(), &mut req))
        }
      };
      let mut res = match res {
        Some(r) => r,
        None => {
          if let Some(alias) = &joker {
            alias.to_response(addr.clone(), "/".to_string(), &mut req)
          } else {
            super::errors::e404()
          }
        }
      };
      res.send(&mut client_conn).unwrap();
    });
  }
}
