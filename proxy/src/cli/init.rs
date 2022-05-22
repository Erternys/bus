use crate::config::{Config, ProxyConfig};
use crate::helper;
use clap::Args;
use std::collections::BTreeMap;
use std::fs::{File, OpenOptions};
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

    let file = OpenOptions::new()
      .write(true)
      .create(true)
      .open("./.bus.yaml")
      .unwrap();

    match serde_yaml::to_writer(file, &config) {
      Ok(_) => (),
      Err(_) => {
        eprintln!("The config file cannot be write, please change the permission");
        exit(1);
      }
    };
  }
}
