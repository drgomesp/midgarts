use std::io::{Cursor, Read};

use byteorder::{LittleEndian, ReadBytesExt};

use crate::fileformat::spr::version::{Version, VersionFormat};
use crate::fileformat::FromBytes;

/// The special header signature for sprite files (SP).
const HEADER_SIGNATURE: &'static str = "SP";

/// Loader submodule for sprite files.
pub(crate) mod loader;

mod version;

/// Indexed image defines images that use the palette.
#[derive(Debug)]
pub(crate) struct IndexedImage {}

/// RGBA image defines images that use RGBA.
#[derive(Debug)]
pub(crate) struct RgbaImage {}

/// The color palette color definition (RGBA).
#[derive(Copy, Clone, Debug, Default)]
pub(crate) struct PaletteColor {
    _r: u8,
    _g: u8,
    _b: u8,
    _reserved: u8,
}

/// The color palette for indexed images.
#[derive(Copy, Clone, Debug)]
pub(crate) struct Palette {
    pub _colors: [PaletteColor; 256],
}

impl Default for Palette {
    fn default() -> Self {
        Palette {
            _colors: [PaletteColor::default(); 256],
        }
    }
}

/// Sprite file format, a compiled form of a texture atlas / sprite sheet.
#[derive(Debug)]
pub(crate) struct SprFile<VersionFormat> {
    /// The version format.
    pub(crate) _version: Version<VersionFormat>,
    /// The number of individual indexed-color images in the atlas
    pub(crate) _indexed_image_count: u16,
    /// The number of individual RGBA images in the atlas
    pub(crate) _rgba_image_count: Option<u16>,
    /// The indexed images.
    pub(crate) _indexed_images: Vec<IndexedImage>,
    /// The RGBA images.
    pub(crate) _rgba_images: Vec<RgbaImage>,
    /// The color palette.
    pub(crate) _palette: Palette,
}

impl FromBytes for SprFile<VersionFormat> {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 2];
        reader
            .read_exact(&mut buf)
            .expect("should read sprite file data");
        let signature = String::from_utf8_lossy(&buf).to_string();

        assert_eq!(
            signature, HEADER_SIGNATURE,
            "invalid sprite file header signature"
        );

        let version = Version::from_bytes(reader.remaining_slice());

        let indexed_image_count = reader
            .read_u16::<LittleEndian>()
            .expect("should read palette image count");

        let rgba_image_count = reader
            .read_u16::<LittleEndian>()
            .expect("should read rgba image count");

        let indexed_images = Vec::with_capacity(indexed_image_count as usize);
        let rgba_images = Vec::with_capacity(rgba_image_count as usize);
        let palette = Palette::default();

        SprFile {
            _version: version,
            _indexed_image_count: indexed_image_count,
            _rgba_image_count: Some(rgba_image_count),
            _indexed_images: indexed_images,
            _rgba_images: rgba_images,
            _palette: palette,
        }
    }
}
