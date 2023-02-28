use std::fs;

use crate::fileformat::{grf::GrfFile, FromBytes};

pub mod fileformat;

fn main() {
    let path = "/Users/drgomesp/go/src/github.com/drgomesp/midgarts/assets/grf/data.grf";
    let bytes = fs::read(path).unwrap_or_else(|_| panic!("failed to load archive from {path}"));

    let grf = GrfFile::from_bytes(&bytes);

    println!("grf = {:#?}", grf);
}
