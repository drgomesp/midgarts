use std::collections::HashMap;
use std::fs;

use crate::fileformat::grf::{entry::GrfEntry, header::GrfHeader, Version};
use crate::fileformat::{FromBytes, Loader};

/// The GRF file.
#[derive(Debug, Default)]
pub struct GrfFile {
    /// The GRF header.
    pub header: GrfHeader,
    /// The GRF entries table.
    pub entries: HashMap<String, GrfEntry>,
}

impl Loader for GrfFile {
    fn load(path: String) -> GrfFile {
        let bytes = fs::read(path).unwrap();

        GrfFile::from_bytes(&bytes)
    }
}

impl FromBytes for GrfFile {
    fn from_bytes(bytes: &[u8]) -> Self {
        let header = GrfHeader::from_bytes(bytes);

        match header.version.try_into() {
            Ok(Version::Version200) => GrfFile {
                header,
                entries: Default::default(),
            },
            _ => todo!("unsupported header version"),
        }
    }
}
