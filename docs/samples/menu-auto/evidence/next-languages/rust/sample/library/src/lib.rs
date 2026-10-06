pub fn value() -> i32 { 21 }
#[cfg(test)]
mod tests { #[test] fn value_is_21() { assert_eq!(super::value(), 21); } }
