use std::fmt;
use std::fmt::{Display, Formatter};
use std::io::BufReader;
use std::marker::PhantomData;

use byteorder::ReadBytesExt;

use crate::fileformat::FromBytes;

/// The version format.
#[derive(Debug, Default)]
pub enum VersionFormat {
    #[default]
    /// Minor first format.
    MinorFirst,
    // /// Major first format.
    // MajorFirst,
}

/// The sprite file version.
#[derive(Copy, Clone, Debug)]
pub struct Version<VersionFormat> {
    /// The minor version component.
    pub _minor: u8,
    /// The major version component.
    pub _major: u8,

    _phantom_data: PhantomData<VersionFormat>,
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
            _minor: minor,
            _major: major,
            _phantom_data: PhantomData,
        };
    }
}

impl<T> Display for Version<T> {
    fn fmt(&self, formatter: &mut Formatter<'_>) -> fmt::Result {
        write!(formatter, "{}.{}", self._major, self._minor)
    }
}
