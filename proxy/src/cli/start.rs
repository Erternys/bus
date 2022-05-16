use std::fs::File;
use std::process::exit;

use clap::Args;

use crate::config::Config;
use crate::conn::Conn;
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
    server.run(|mut req| {
      let mut conn = Conn::new("localhost:2002").unwrap();
      let mut body = Vec::new();
      let mut res = Vec::new();

      req.read(&mut body).unwrap();
      conn.write(&body).unwrap();

      conn.read(&mut res).unwrap();
      req.write(&res).unwrap();
      conn.close().unwrap();
      req.close().unwrap();
    });
  }
}
