use std::collections::BTreeMap;

#[derive(Debug)]
pub enum Alias {
  File(String),
  Http(String),
  Empty
}

impl Alias {
  pub fn from_iter(iter: BTreeMap<String, String>) -> (BTreeMap<String, Alias>, Option<Alias>) {
    let mut map = BTreeMap::new();
    let mut joker = None;
    let joker_key = "*".to_string();
    for (key, alias) in iter{
      if key == joker_key {
          joker = Some(Self::from_str(alias))
      }else{
        map.insert(key, Self::from_str(alias));
      }
    }
    (map, joker)
  }
  pub fn from_str(s: String) -> Alias {
    let splitted = s.splitn(2, ":").collect::<Vec<&str>>();
    if let [schema, value] = *splitted.as_slice() {
      match schema {
        "http" => Alias::Http(value.to_string()),
        "file" => Alias::File(value.to_string()),
        _ => Alias::Empty
      }
    } else if s.len() == 0{
      Alias::Empty
    } else if s.chars().next().unwrap() == '.'{
      Alias::File(s)
    }else{
      Alias::Http(s)
    }
  }
}
