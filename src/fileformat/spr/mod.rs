use std::io::{Cursor, Read};

use crate::fileformat::FromBytes;

/// Loader submodule for sprite files.
pub mod loader;

/// Sprite file format.
#[derive(Debug, Default)]
pub struct SprFile {
    /// Sprite file header signature (SP).
    pub signature: String,
}

impl SprFile {}

impl FromBytes for SprFile {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 2];
        reader
            .read_exact(&mut buf)
            .expect("should read sprite file data");

        SprFile {
            signature: String::from_utf8_lossy(&buf).to_string(),
        }
    }
}
