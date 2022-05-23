use crate::config::{Config, ProxyConfig};
use crate::helper;
use clap::Args;
use std::collections::BTreeMap;
use std::fs::{File, OpenOptions};
use std::io::Write;
use std::process::exit;

#[derive(Debug, Args)]
pub struct Init {
  #[clap(long, short)]
  yes: bool
}

impl Init {
  pub fn call(&self) {
    let file = match File::open("./.bus.yaml") {
      Ok(f) => f,
      Err(_) => {
        eprintln!("Please execute the command `bus init` before");
        exit(1);
      }
    };

    let mut config: Config = match serde_yaml::from_reader(&file) {
      Ok(c) => c,
      Err(_) => {
        eprintln!("The config file cannot be read, please change the permission");
        exit(1);
      }
    };
    if let Some(_) = config.proxy {
      return
    }

    let port = helper::input(&format!("port: ({}) ", "80"), "80");
    let open = helper::str_to_bool(helper::input("open: ", "true"));

    let proxy = ProxyConfig {
      port,
      open: Some(open),
      aliases: BTreeMap::new()
    };

    config.proxy = Some(proxy);

    let mut file = OpenOptions::new()
      .write(true)
      .create(true)
      .open("./.bus.yaml")
      .unwrap();

    let content = match serde_yaml::to_string(&config) {
      Ok(s) => s,
      Err(_) => {
        eprintln!("The config cannot be created, an error from a non-utf8 character");
        exit(1);
      }
    };
    let content = content[3..].trim().as_bytes();

    match file.write(content) {
      Ok(_) => (),
      Err(_) => {
        eprintln!("The config file cannot be writed, please change the permission");
        exit(1);
      }
    };
  }
}
