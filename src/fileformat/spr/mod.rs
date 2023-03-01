use std::io::{Cursor, Read};
use std::marker::PhantomData;
use std::string::ToString;

use byteorder::ReadBytesExt;

use crate::fileformat::FromBytes;

const HEADER_SIGNATURE: &'static str = "SP";

/// Loader submodule for sprite files.
pub mod loader;

/// The version format.
#[derive(Debug, Default)]
pub enum VersionFormat {
    #[default]
    /// Minor first format.
    MinorFirst,
    /// Major first format.
    MajorFirst,
}

/// The sprite file version.
#[derive(Debug, Default)]
pub struct Version<VersionFormat> {
    /// The minor version component.
    pub minor: u8,
    /// The major version component.
    pub major: u8,
    phantom_data: PhantomData<VersionFormat>,
}

impl FromBytes for Version<VersionFormat> {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let minor = reader
            .read_u8()
            .expect("should read version minor component");

        let major = reader
            .read_u8()
            .expect("should read version major component");

        return Version {
            minor,
            major,
            phantom_data: PhantomData,
        };
    }
}

/// Sprite file format.
#[derive(Debug, Default)]
pub struct SprFile<VersionFormat> {
    /// The version format.
    pub version: Version<VersionFormat>,
}

impl SprFile<VersionFormat> {}

impl FromBytes for SprFile<VersionFormat> {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 2];
        reader
            .read_exact(&mut buf)
            .expect("should read sprite file data");
        let mut signature = String::from_utf8_lossy(&buf).to_string();

        assert_eq!(
            signature, HEADER_SIGNATURE,
            "invalid sprite file header signature"
        );

        let version = Version::from_bytes(reader.remaining_slice());

        SprFile { version }
    }
}
