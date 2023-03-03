/// The special header signature for sprite files (SP).
const HEADER_SIGNATURE: &'static str = "SP";

/// Loader submodule for sprite files.
pub(crate) mod loader;

mod file;
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
