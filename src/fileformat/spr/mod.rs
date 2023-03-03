/// The special header signature for sprite files (SP).
const HEADER_SIGNATURE: &'static str = "SP";

/// Loader submodule for sprite files.
pub mod loader;

mod file;
mod version;

#[derive(Clone, Debug)]
/// Encoded data represents an RLE-encoded chunk of data.
pub struct EncodedData(pub Vec<u8>);

/// Indexed image defines images that use the palette.
#[derive(Debug)]
pub struct IndexedImage {
    /// The image width.
    pub width: u16,
    /// The image height.
    pub height: u16,
    /// The RLE-encoded data.
    pub encoded_data: Option<EncodedData>,
    /// The data.
    pub data: Option<Vec<u8>>,
}

/// RGBA image defines images that use RGBA.
#[derive(Debug)]
pub struct RgbaImage {}

/// The color palette color definition (RGBA).
#[derive(Copy, Clone, Debug, Default)]
pub struct PaletteColor {
    _r: u8,
    _g: u8,
    _b: u8,
    _reserved: u8,
}

/// The color palette for indexed images.
#[derive(Copy, Clone, Debug)]
pub struct Palette {
    /// The palette colors.
    pub _colors: [PaletteColor; 256],
}

impl Default for Palette {
    fn default() -> Self {
        Palette {
            _colors: [PaletteColor::default(); 256],
        }
    }
}
