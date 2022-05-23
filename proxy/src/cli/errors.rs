use crate::http_parser::Response;

static LAYOUT: &str = "<!DOCTYPE html>\
<html lang=\"en\">\
<head>\
  <meta charset=\"UTF-8\">\
  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\
  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\
  <title>Proxy Error</title>\
</head>\
<body style=\"text-align: center; font-family: arial, sans-serif\">\
  <h1>Error %s</h1>\
  <h2>%s</h2>\
  <hr>\
  <p>bus v%s<br>bus proxy %s</p>
</body>\
</html>";

macro_rules! cformat {
  ($f: expr, $($arg: expr),+ $(,)?) => {{
    let mut s = $f.to_string();
    for item in [$($arg),+] {
      s = s.replacen("%s", item, 1);
    }
    s
  }};
}

pub fn e400() -> Response {
  let mut res = Response::default();
  res.status = 400;
  res.message = String::from("Bad Request");
  res.body = cformat!(
    LAYOUT, 
    "400", 
    "Bad Request", 
    "0.1.0-beta", 
    "0.1.0-beta"
  )
    .as_bytes()
    .to_vec();

  res
}

pub fn e404() -> Response {
  let mut res = Response::default();
  res.status = 404;
  res.message = String::from("Not Found");
  res.body = cformat!(
    LAYOUT, 
    "404",
    "Not Found", 
    "0.1.0-beta", 
    "0.1.0-beta"
  )
    .as_bytes()
    .to_vec();

  res
}

pub fn e406() -> Response {
  let mut res = Response::default();
  res.status = 406;
  res.message = String::from("Not Acceptable");
  res.body = cformat!(
    LAYOUT, 
    "406",
    "Not Acceptable",
    "0.1.0-beta", 
    "0.1.0-beta"
  )
    .as_bytes()
    .to_vec();

  res
}

pub fn e500() -> Response {
  let mut res = Response::default();
  res.status = 500;
  res.message = String::from("Internal Server Error");
  res.body = cformat!(
    LAYOUT, 
    "500",
    "Internal Server Error", 
    "0.1.0-beta", 
    "0.1.0-beta"
  )
    .as_bytes()
    .to_vec();

  res
}