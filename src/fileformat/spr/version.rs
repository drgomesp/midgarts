use std::io::{BufReader, Cursor};
use std::marker::PhantomData;

use byteorder::ReadBytesExt;

use crate::fileformat::FromBytes;

/// The version format.
#[derive(Debug, Default)]
pub(crate) enum VersionFormat {
    #[default]
    /// Minor first format.
    MinorFirst,
    /// Major first format.
    MajorFirst,
}

/// The sprite file version.
#[derive(Debug, Default)]
pub(crate) struct Version<VersionFormat> {
    /// The minor version component.
    pub(crate) minor: u8,
    /// The major version component.
    pub(crate) major: u8,

    phantom_data: PhantomData<VersionFormat>,
}

impl FromBytes for Version<VersionFormat> {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = BufReader::new(bytes);

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
