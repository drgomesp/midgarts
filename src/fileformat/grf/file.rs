use std::collections::HashMap;
use std::fs;
use std::io::{Cursor, Read, Seek, SeekFrom};

use byteorder::{LittleEndian, ReadBytesExt};
use yazi::*;

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
            Ok(Version::Version200) => {
                let mut f = GrfFile {
                    header,
                    entries: Default::default(),
                };

                let mut reader = Cursor::new(bytes);
                reader
                    .seek(SeekFrom::Start(f.header.file_table_offset as u64))
                    .expect("should seek to file table");

                let compressed_size = reader
                    .read_u32::<LittleEndian>()
                    .unwrap_or_else(|_| panic!("failed to read compressed size"));

                let uncompressed_size = reader
                    .read_u32::<LittleEndian>()
                    .unwrap_or_else(|_| panic!("failed to read uncompressed size"));

                debug!(
                    "compressed size: {}, uncompressed size: {}",
                    compressed_size, uncompressed_size
                );

                let mut compressed = vec![0u8; compressed_size as usize];
                reader.read_exact(&mut compressed).unwrap();
                let (decompressed, _checksum) = decompress(&compressed, Format::Zlib).unwrap();

                for _i in 0..f.header.file_count {
                    let mut entry = GrfEntry::from_bytes(&decompressed);
                    f.entries.insert(entry.file_name.to_lowercase(), entry);
                }

                f
            }
            _ => todo!("unsupported header version"),
        }
    }
}
