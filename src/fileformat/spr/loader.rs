use std::io::{Cursor, SeekFrom};

use crate::fileformat::grf::file::GrfFile;
use crate::fileformat::grf::header::HEADER_SIZE;
use crate::fileformat::spr::{SprFile, VersionFormat};
use crate::fileformat::{FromBytes, Loader};

/// Loader of sprite (.SPR) files.
#[derive(Debug)]
pub struct SpriteLoader<'a> {
    grf_file: &'a GrfFile,
}

impl<'a> SpriteLoader<'a> {
    /// Creates a new sprite loader.
    pub fn new(grf_file: &'a GrfFile) -> Self {
        SpriteLoader { grf_file }
    }
}

impl<'a> SpriteLoader<'a> {
    /// Loads a sprite.
    pub fn load(&mut self, path: &'static str) -> Result<SprFile<VersionFormat>, String> {
        let grf_entry = self.grf_file.get_entry(path);

        Ok(SprFile::from_bytes(&grf_entry.data))
    }
}
