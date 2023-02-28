//! I document the current module!

#![deny(
    missing_docs,
    missing_debug_implementations,
    trivial_casts,
    trivial_numeric_casts,
    unsafe_code,
    unstable_features,
    unused_import_braces,
    unused_qualifications
)]

#[macro_use]
extern crate chrono;
extern crate dotenv;
#[macro_use]
extern crate log;
extern crate pretty_env_logger;

use std::env;
use std::fs;

use dotenv::dotenv;

use fileformat::grf::file;

use crate::fileformat::grf::file::GrfFile;
use crate::fileformat::Loader;

/// File format module defines all file formats
pub(crate) mod fileformat;

fn main() {
    dotenv().ok();
    pretty_env_logger::init();

    info!("started {}", chrono::Utc::now());

    let grf = GrfFile::load("assets/data.grf".to_string());

    debug!("{:?}", grf);
}
