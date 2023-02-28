use std::*;

use crate::fileformat::FromBytes;

/// The GRF file header.
pub struct GrfHeader {
    _signature: [u8; 15],
    _encryption: [u8; 15],
    _file_table_offset: u32,
    _reserved_files: u32,
    _file_count: u32,
    _version: u32,
}

impl FromBytes for GrfHeader {
    fn from_bytes(_bytes: &[u8]) -> Self {
        todo!()
    }
}

/// The GRF file.
pub struct GrfFile {
    _header: GrfHeader,
}

impl FromBytes for GrfFile {
    fn from_bytes(_bytes: &[u8]) -> Self {
        todo!()
    }
}
