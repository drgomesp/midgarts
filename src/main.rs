#![feature(cursor_remaining)]
//! I document the current module!
#![deny(
    missing_docs,
    missing_debug_implementations,
    trivial_casts,
    trivial_numeric_casts,
    unsafe_code,
    unused_imports,
    unused_import_braces,
    unused_qualifications
)]

extern crate chrono;
extern crate core;
#[macro_use]
extern crate log;
extern crate pretty_env_logger;

use crate::fileformat::grf::file::GrfFile;
use crate::fileformat::spr::loader::SpriteLoader;
use crate::fileformat::Loader;

/// File format module defines all file formats
pub(crate) mod fileformat;

fn main() {
    pretty_env_logger::init();

    info!("started {}", chrono::Utc::now());

    let grf = GrfFile::load("assets/data.grf");
    let mut sprite_loader = SpriteLoader::new(&grf);

    debug!(
        "{:?}",
        sprite_loader.load("data\\sprite\\àî°£á·\\¸öåë\\³²\\±â»ç_³².spr")
    );

    info!("finished {}", chrono::Utc::now());
}
