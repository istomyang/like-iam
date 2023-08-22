#![allow(dead_code)]

#[test]
fn test() {
    enum A {
        Nil,
        Str(u32),
    }

    println!("{:?}", std::mem::size_of::<A>());
    println!("{:?}", std::mem::size_of_val(&A::Nil));
    println!("{:?}", std::mem::size_of_val(&A::Str(13)));
}

mod chain;
mod config;
mod rest;
mod tradition;
