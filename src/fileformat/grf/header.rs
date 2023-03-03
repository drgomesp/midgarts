use std::io::{Cursor, Read};

use byteorder::{LittleEndian, ReadBytesExt};

use crate::fileformat::FromBytes;

/// The GRF header size constant.
pub const HEADER_SIZE: usize = 46;

/// The GRF file header.
#[derive(Debug, Default)]
pub struct GrfHeader {
    _encryption: [u8; 15],
    /// The offset of the file table
    pub entry_table_offset: u32,
    /// The reserved files.
    pub reserved: u32,
    /// The entry count.
    pub entry_count: u32,
    /// The GRF version.
    pub version: u32,
}

impl FromBytes for GrfHeader {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let mut buf = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        let sig: String = String::from_utf8_lossy(&buf).parse().unwrap();
        assert_eq!(sig, "Master of Magic", "invalid file header signature");

        let encryption = [0u8; 15];
        reader
            .read_exact(&mut buf)
            .expect("should read header signature");

        GrfHeader {
            _encryption: encryption,
            entry_table_offset: reader.read_u32::<LittleEndian>().unwrap() + HEADER_SIZE as u32,
            reserved: reader.read_u32::<LittleEndian>().unwrap(),
            entry_count: reader.read_u32::<LittleEndian>().unwrap(),
            version: reader.read_u32::<LittleEndian>().unwrap(),
        }
    }
}
