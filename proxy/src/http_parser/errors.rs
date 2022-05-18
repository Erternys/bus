#[derive(Debug)]
pub enum ParserErrorKind {
  Parsing,
  Reading,
  Utf8
}

#[derive(Debug)]
pub struct ParserError<'s> {
  kind: ParserErrorKind,
  message: &'s str
}

impl<'s> ParserError<'s> {
  pub fn new(kind: ParserErrorKind, message: &'s str) -> Self {
    Self { kind, message }
  }
}

impl<'s> std::fmt::Display for ParserError<'s> {
  fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
    write!(f, "{}", self.message)
  }
}

impl<'s> std::error::Error for ParserError<'s> {}
