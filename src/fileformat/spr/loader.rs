use crate::fileformat::grf::file::GrfFile;
use crate::fileformat::spr::file::SprFile;
use crate::fileformat::spr::version::VersionFormat;
use crate::fileformat::FromBytes;

/// Loader of sprite (.SPR) files.
#[derive(Debug)]
pub(crate) struct SpriteLoader<'a> {
    grf_file: &'a GrfFile,
}

impl<'a> SpriteLoader<'a> {
    /// Creates a new sprite loader.
    pub(crate) fn new(grf_file: &'a GrfFile) -> Self {
        SpriteLoader { grf_file }
    }
}

impl<'a> SpriteLoader<'a> {
    /// Loads a sprite.
    pub(crate) fn load(&mut self, path: &'static str) -> Result<SprFile<VersionFormat>, String> {
        let grf_entry = self.grf_file.get_entry(path);

        Ok(SprFile::from_bytes(&grf_entry.data))
    }
}
