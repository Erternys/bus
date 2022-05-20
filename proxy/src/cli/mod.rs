mod init;
mod start;
// mod stop;

use std::ffi::OsString;

use clap::{ColorChoice, Parser, Subcommand};

#[derive(Debug, Parser)]
#[clap(
  name = "bus-proxy", 
  bin_name = "bus-proxy",
  version = "1.0.0",
  about = "A fictional versioning CLI", 
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