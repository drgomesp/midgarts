#[macro_use]
extern crate log;
extern crate midgarts;

use midgarts::fileformat::{grf::file::GrfFile, Loader};

fn main() {
    let grf = GrfFile::load("assets/data.grf");

    println!("entry = {:#?}", grf.get_entry("data\\sprite\\shadow.spr"));
}
