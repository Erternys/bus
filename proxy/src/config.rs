use std::collections::HashMap;

use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Config {
  pub name: String,
  pub version: String,
  pub description: String,
  pub repository: String,
  pub manager: String,
  pub packages: Vec<Package>,
  pub proxy: Option<ProxyConfig>
}
#[derive(Debug, Serialize, Deserialize)]
pub struct Package {
  pub path: String,
  pub name: String,
  pub extend: String
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProxyConfig {
  pub port: String,
  pub open: Option<bool>,
  pub aliases: HashMap<String, String>
}
