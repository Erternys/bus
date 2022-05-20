use std::fs::File;
use std::process::exit;

use clap::Args;
use url::Url;

use crate::config::Config;
use crate::conn::Conn;
use crate::http_parser::{Request, Response};
use crate::server::Server;

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

    let server = Server::new(addr);
    server.run(move |mut client_conn, addr| {
      let mut res = Response::default();
      let req = Request::from_conn(&mut client_conn);

      let mut req = match req {
        Ok(r) => r,
        Err(_) => {
          res.status = 400;
          res.message = String::from("Bad Request");
          res.body = "<h1>Error 400</h1><h3>Bad Request</h3><hr>"
            .as_bytes()
            .to_vec();
          
          res.send(&mut client_conn).unwrap();
          return
        }
      };
      for (path, alias) in &proxy.aliases {
        if req.url.starts_with(&*path) {
          let url = match Url::parse(&alias) {
            Ok(u) => u,
            Err(_) => {
              res.status = 400;
              res.message = String::from("Bad Request");
              res.body = "<h1>Error 400</h1><h3>Bad Request</h3><hr>"
                .as_bytes()
                .to_vec();
              
              res.send(&mut client_conn).unwrap();
              return
            }
          };

          let host = url.host_str().unwrap_or(&addr);
          let addr = if let Some(port) = url.port() {
            format!("{}:{}", &host, port)
          } else {
            String::from(host.clone())
          };

          req.headers.insert(String::from("Host"), host.to_string());

          let mut server_conn = Conn::new(&addr).unwrap();
          req.send(&mut server_conn);
          res.vacuum(&mut server_conn);
          res.send(&mut client_conn);
          return
        }
      }

      res.status = 404;
      res.message = String::from("Not Found");
      res.body = "<h1>Error 404</h1><h3>Not Found</h3><hr>"
        .as_bytes()
        .to_vec();
      
      res.send(&mut client_conn).unwrap();
    });
  }
}
