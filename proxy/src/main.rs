mod cli;
mod config;
mod conn;
mod helper;
mod http_parser;
mod server;

use clap::Parser;
use cli::*;

fn main() {
  let app = App::parse();

  match app.command {
    Commands::Init(init) => init.call(),
    Commands::Start(start) => start.call(),
    Commands::External(_) => todo!()
  }
}
