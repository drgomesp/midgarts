use std::fmt::{Display, Formatter};
use std::io::{Cursor, Read};

use byteorder::{LittleEndian, ReadBytesExt};
use log::debug;

use crate::fileformat::spr::version::{Version, VersionFormat};
use crate::fileformat::spr::{IndexedImage, Palette, RgbaImage, HEADER_SIGNATURE};
use crate::fileformat::FromBytes;

#[derive(Debug)]
pub enum SpriteType {
    PAL,
    RGBA,
}

impl Display for SpriteType {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            SpriteType::PAL => write!(f, "PAL"),
            SpriteType::RGBA => write!(f, "RGBA"),
        }
    }
}

/// Sprite file format, a compiled form of a texture atlas / sprite sheet.
#[derive(Debug)]
pub struct SprFile<VersionFormat> {
    /// The sprite type (PAL or RGBA).
    pub _type: SpriteType,
    /// The version format.
    pub _version: Version<VersionFormat>,
    /// The number of individual indexed-color images in the atlas
    pub _indexed_image_count: u16,
    /// The number of individual RGBA images in the atlas
    pub _rgba_image_count: Option<u16>,
    /// The indexed images.
    pub _indexed_images: Vec<IndexedImage>,
    /// The RGBA images.
    pub _rgba_images: Vec<RgbaImage>,
    /// The color palette.
    pub _palette: Palette,
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

        let mut buf = [0u8; 2];
        reader
            .read_exact(&mut buf)
            .expect("should read sprite file data");
        let version = Version::from_bytes(vec![buf[0], buf[1]].as_slice());

        let indexed_image_count = reader
            .read_u16::<LittleEndian>()
            .expect("should read palette image count");

        let rgba_image_count = reader
            .read_u16::<LittleEndian>()
            .expect("should read rgba image count");
        let mut indexed_images = Vec::with_capacity(indexed_image_count as usize);

        let mut rgba_images = Vec::with_capacity(rgba_image_count as usize);

        let palette = Palette::default();

        let mut spr_file = SprFile {
            _type: SpriteType::PAL,
            _version: version,
            _indexed_image_count: indexed_image_count,
            _rgba_image_count: Some(rgba_image_count),
            _indexed_images: indexed_images,
            _rgba_images: rgba_images,
            _palette: palette,
        };

        spr_file.read_indexed_images(reader);

        debug!("sprite = {:?}", spr_file);

        spr_file
    }
}

impl SprFile<VersionFormat> {
    fn read_indexed_images(&mut self, mut reader: Cursor<&[u8]>) -> Vec<IndexedImage> {
        (0..self._indexed_image_count)
            .map(|_i| {
                let width = reader.read_u16::<LittleEndian>().unwrap();
                let height = reader.read_u16::<LittleEndian>().unwrap();
                let size = reader.read_u16::<LittleEndian>().unwrap();

                let mut skip = vec![0u8; size as usize];
                let _ = reader.read_exact(&mut skip);

                let size: u32 = (width * height) as u32;
                let mut data = vec![0u8; size as usize];

                reader
                    .read_exact(&mut data)
                    .expect("should read sprite image data");

                debug!("width = {:?}, height = {:?}", width, height);

                IndexedImage {
                    width,
                    height,
                    encoded_data: None,
                    data: None,
                }
            })
            .collect()
    }
}
