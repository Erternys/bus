mod alias;
mod errors;
mod init;
mod start;

use std::ffi::OsString;

use clap::{ColorChoice, Parser, Subcommand};

#[derive(Debug, Parser)]
#[clap(
  name = "bus-proxy", 
  bin_name = "bus-proxy",
  version = "0.1.0-beta",
  about = "Server proxy for bus",
  long_about = None,
  color = ColorChoice::Never
)]
pub struct App {
  #[clap(subcommand)]
  pub command: Commands
}

#[derive(Debug, Subcommand)]
pub enum Commands {
  /// Add in the config file the configuration for the proxy
  Init(init::Init),

  /// Starts a proxy server with the info given in the configuration
  Start(start::Start),

  #[clap(external_subcommand)]
  External(Vec<OsString>)
}
