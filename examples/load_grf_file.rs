#[macro_use]
extern crate log;
extern crate midgarts;

use midgarts::fileformat::{grf::file::GrfFile, Loader};

fn main() {
    let grf = GrfFile::load("assets/data.grf");

    println!("entries = {:?}", grf.entry_count());
}
