use std::net::TcpListener;

use super::conn::Conn;

pub struct Server {
  addr: String
}

impl Server {
  pub fn new(addr: String) -> Self {
    Self {
      addr
    }
  }

  pub fn run<F>(self, mut callback: F)
  where
    F: FnMut(Conn, &String) + 'static
  {
    let listener = TcpListener::bind(&self.addr).unwrap();

    for stream in listener.incoming() {
      callback(Conn::from_stream(stream.unwrap()).unwrap(), &self.addr);
    }
  }
}
