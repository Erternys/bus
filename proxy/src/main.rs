mod conn;
mod server;

use server::Server;

fn main() {
  let server = Server::new("127.0.0.1:2001");
  server.run(|mut req| {
    let mut conn = conn::Conn::new("localhost:2002").unwrap();
    let mut body = Vec::new();
    let mut res = Vec::new();

    req.read(&mut body).unwrap();
    conn.write(&body).unwrap();

    conn.read(&mut res).unwrap();
    req.write(&res).unwrap();
  });
}
