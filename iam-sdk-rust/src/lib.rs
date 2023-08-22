pub fn add(left: usize, right: usize) -> usize {
    left + right
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}

pub mod error;

macro_rules! if_iam {
    ($($item:item)*) => {$(
        #[cfg(feature = "iam")]
        $item
    )*}
}

if_iam! {
    mod iam;
    mod api;
}

macro_rules! if_tms {
    ($($item:item)*) => {$(
        #[cfg(feature = "tms")]
        $item
    )*}
}

if_tms! {
    mod tms;
}
