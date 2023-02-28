use std::convert::TryFrom;
use std::io::{Cursor, Read};
use std::*;

use byteorder::{LittleEndian, ReadBytesExt};

use crate::fileformat::FromBytes;

pub(crate) enum Version {
    Version200 = 0x200,
}

impl TryFrom<u32> for Version {
    type Error = ();

    fn try_from(v: u32) -> Result<Self, Self::Error> {
        match v {
            v if v == Version::Version200 as u32 => Ok(Version::Version200),
            _ => todo!("invalid header version {}", v),
        }
    }
}

/// The GRF file.
#[derive(Debug, Default)]
pub struct GrfFile {
    /// The GRF header.
    pub header: GrfHeader,
}

impl GrfFile {
    pub fn load(path: String) -> GrfFile {
        let bytes = fs::read(path).unwrap();

        GrfFile::from_bytes(&bytes)
    }
}

impl FromBytes for GrfFile {
    fn from_bytes(bytes: &[u8]) -> Self {
        let header = GrfHeader::from_bytes(bytes);

        match header.version.try_into() {
            Ok(Version::Version200) => GrfFile { header },
            _ => todo!("unsupported header version"),
        }
    }
}

/// The GRF file header.
#[derive(Debug, Default)]
pub struct GrfHeader {
    encryption: [u8; 15],
    file_table_offset: u32,
    reserved_files: u32,
    file_count: u32,
    version: u32,
}

impl FromBytes for GrfHeader {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        let mut sig: String = String::from_utf8_lossy(&buf).parse().unwrap();
        assert_eq!(sig, "Master of Magic", "invalid file header signature");

        let encryption = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        GrfHeader {
            encryption,
            file_table_offset: reader.read_u32::<LittleEndian>().unwrap(),
            reserved_files: reader.read_u32::<LittleEndian>().unwrap(),
            file_count: reader.read_u32::<LittleEndian>().unwrap(),
            version: reader.read_u32::<LittleEndian>().unwrap(),
        }
    }
}
