use std::net::TcpListener;

use super::conn::Conn;

pub struct Server<'s> {
  addr: &'s str
}

impl<'s> Server<'s> {
  pub fn new(addr: &'s str) -> Self {
    Self {
      addr
    }
  }

  pub fn run<F>(self, mut callback: F)
  where
    F: FnMut(Conn) + 'static
  {
    let listener = TcpListener::bind(self.addr).unwrap();

    for stream in listener.incoming() {
      callback(Conn::from_stream(stream.unwrap()).unwrap());
    }
  }
}
