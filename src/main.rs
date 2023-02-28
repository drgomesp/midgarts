use std::fs;

use crate::fileformat::{grf::GrfFile, FromBytes};

pub mod fileformat;

fn main() {
    let grf = GrfFile::load(
        "/Users/drgomesp/go/src/github.com/drgomesp/midgarts/assets/grf/data.grf".to_string(),
    );

    println!("grf = {:#?}", grf);
}
