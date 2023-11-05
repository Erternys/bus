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

    loop {
      let stream = listener.accept().unwrap().0;
      match Conn::from_stream(stream) {
        Ok(conn) => callback(conn, &self.addr),
        Err(_) => {
          let time = chrono::Local::now().to_rfc2822();
          eprintln!("[{time} ERROR SERVER]: The connection die.");
        }
      }
    }
  }
}
