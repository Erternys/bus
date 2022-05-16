mod init;

pub use init::Init;

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
  Init(Init),

  #[clap(external_subcommand)]
  External(Vec<OsString>)
}
