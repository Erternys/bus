use std::io::{self, Read, Write};
use std::net::TcpStream;
use std::time::Duration;

const CHUNK_LENGTH_MAX: usize = 1024;

pub struct Conn {
  stream: TcpStream
}

impl Conn {
  pub fn new(addr: &str) -> io::Result<Self> {
    Self::from_stream(TcpStream::connect(addr)?)
  }

  pub fn from_stream(stream: TcpStream) -> io::Result<Self> {
    let mut connection = Self {
      stream
    };

    connection.set_write_timeout(Duration::from_secs(1))?;
    connection.set_read_timeout(Duration::from_secs(1))?;

    Ok(connection)
  }

  pub fn set_read_timeout(&mut self, dur: Duration) -> io::Result<()> {
    self.stream.set_read_timeout(Some(dur))
  }

  pub fn read(&mut self, buf: &mut Vec<u8>) -> io::Result<()> {
    let mut chunk: [u8; CHUNK_LENGTH_MAX] = [0; CHUNK_LENGTH_MAX];
    loop {
      let chunk_length = self.stream.read(&mut chunk)?;
      buf.append(&mut chunk[..chunk_length].to_vec());
      if chunk_length < CHUNK_LENGTH_MAX {
        break
      }
    }

    Ok(())
  }

  pub fn set_write_timeout(&mut self, dur: Duration) -> io::Result<()> {
    self.stream.set_write_timeout(Some(dur))
  }

  pub fn write(&mut self, data: &[u8]) -> io::Result<usize> {
    self.stream.write(data)
  }
}
