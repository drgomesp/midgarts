use std::io::{Cursor, Read};
use std::*;

use byteorder::{LittleEndian, ReadBytesExt};

use crate::fileformat::FromBytes;

/// The GRF file header.
#[derive(Debug, Default)]
pub struct GrfHeader {
    signature: String,
    encryption: [u8; 15],
    _file_table_offset: u32,
    _reserved_files: u32,
    _file_count: u32,
    _version: u32,
}

impl FromBytes for GrfHeader {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        let mut signature = String::from_utf8_lossy(&buf).parse().unwrap();
        assert_eq!(signature, "Master of Magic");

        let encryption = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        GrfHeader {
            signature,
            encryption,
            _file_table_offset: reader.read_u32::<LittleEndian>().unwrap(),
            _reserved_files: reader.read_u32::<LittleEndian>().unwrap(),
            _file_count: reader.read_u32::<LittleEndian>().unwrap(),
            _version: reader.read_u32::<LittleEndian>().unwrap(),
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

        GrfFile { header }
    }
}
