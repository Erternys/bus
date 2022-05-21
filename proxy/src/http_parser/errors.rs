#[derive(Debug)]
pub enum HttpErrorKind {
  Closing,
  Parsing,
  Reading,
  Sending,
  Utf8
}

#[derive(Debug)]
pub struct HttpError<'s> {
  kind: HttpErrorKind,
  message: &'s str
}

impl<'s> HttpError<'s> {
  pub fn new(kind: HttpErrorKind, message: &'s str) -> Self {
    Self { kind, message }
  }
}

impl<'s> std::fmt::Display for HttpError<'s> {
  fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
    write!(f, "{}", self.message)
  }
}

impl<'s> std::error::Error for HttpError<'s> {}
