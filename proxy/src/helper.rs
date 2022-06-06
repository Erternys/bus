use std::io::{stdin, stdout, Write};

pub fn input(prompt: &str, default_value: &str) -> String {
  print!("{}", prompt);
  if let Err(_) = stdout().flush() {
    return default_value.to_string()
  }

  let mut line = String::new();
  if let Ok(_) = stdin().read_line(&mut line) {
    let line = line.trim();
    if line.len() == 0 {
      return default_value.to_string()
    }
    return line.to_string()
  }

  default_value.to_string()
}

pub fn str_to_bool(s: String) -> bool {
  match &*s {
    "true" | "yes" | "y" | "" => true,
    _ => false
  }
}
