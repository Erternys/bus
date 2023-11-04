use std::fs::File;
use std::process::exit;

use clap::Args;

use crate::config::Config;
use crate::http_parser::{HttpErrorKind, Request};
use crate::server::Server;

use super::alias::Alias;
use super::errors::{e400, e404, e500};

#[derive(Debug, Args)]
pub struct Start {
  #[clap(long)]
  log: bool
}

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
      Err(e) => {
        eprintln!("{e}");
        exit(1);
      }
    };

    let logging = config.proxy.as_ref()
      .and_then(|proxy| proxy.on_script.as_ref())
      .and_then(|on_script| on_script.log)
      .unwrap_or(self.log);

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

    let server = Server::new(addr.clone());
    if logging {
      let time = chrono::Local::now().to_rfc2822();
      println!("[{time} INFO]: Server start on {addr}");
    }
    server.run(move |mut client_conn, addr| {
      let mut req = match Request::from_conn(&mut client_conn) {
        Ok(r) => r,
        Err(error) => {
          let mut res = match error.kind {
            HttpErrorKind::Reading => e500(),
            _ => e400()
          };

          if logging {
            let time = chrono::Local::now().to_rfc2822();
            eprintln!("[{time} ERROR]: {error}");
          }

          res.send(&mut client_conn).unwrap();
          return
        }
      };

      let mut iter = aliases.iter();

      let res = loop {
        let (path, alias) = match iter.next() {
          Some((path, alias)) => (path, alias),
          None => break None
        };

        if req.url.starts_with(path) {
          break Some(alias.to_response(addr.clone(), path.clone(), &mut req))
        }
      };
      let mut res = match res {
        Some(r) => r,
        None =>
          if let Some(alias) = &joker {
            alias.to_response(addr.clone(), "/".to_string(), &mut req)
          } else {
            e404()
          },
      };

      if logging {
        let time = chrono::Local::now().to_rfc2822();
        if res.status >= 400 {
          eprintln!("[{time} ERROR]: {} for {}", res.status, req.url);
        } else {
          println!("[{time} INFO]: {} for {}", res.status, req.url);
        }
      }

      res.send(&mut client_conn).unwrap();
    });
  }
}
