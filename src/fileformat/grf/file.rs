use std::collections::HashMap;
use std::fs;
use std::io::prelude::*;
use std::io::{BufReader, Cursor, Read, Seek, SeekFrom};

use byteorder::{LittleEndian, ReadBytesExt};
use encoding_rs::WINDOWS_1252;
use yazi::*;

use crate::fileformat::grf::entry::ENTRY_HEADER_SIZE;
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

impl GrfFile {
    pub fn file_count(&self) -> usize {
        return (self.header.file_count - self.header.reserved_files) as usize - 7;
    }
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
                reader
                    .read_exact(&mut compressed)
                    .expect("should read compressed data");
                let (decompressed, _checksum) = decompress(&compressed, Format::Zlib).unwrap();

                f.entries = HashMap::with_capacity(f.header.file_count as usize);
                // let mut reader = BufReader::new(decompressed.as_slice());
                let mut reader = Cursor::new(&decompressed);

                for _i in 0..f.file_count() {
                    let mut buf = vec![];
                    let mut string_decoder = encoding_rs_io::DecodeReaderBytesBuilder::new();
                    reader
                        .read_until(b'\0', &mut buf)
                        .expect("should read file name as string");

                    let mut string_decoder = string_decoder
                        .encoding(Some(WINDOWS_1252))
                        .build(Cursor::new(&buf[0..buf.len() - 1]));

                    let mut file_name = String::new();
                    string_decoder
                        .read_to_string(&mut file_name)
                        .unwrap_or_else(|_| panic!("failed to read file name"));

                    let mut buf = vec![0u8; ENTRY_HEADER_SIZE];
                    reader
                        .read_exact(&mut buf)
                        .expect("should read entry header");

                    let mut entry = GrfEntry::from_bytes(&buf);
                    entry.file_name = file_name.to_lowercase();

                    f.entries.insert(entry.file_name.to_lowercase(), entry);
                }

                f
            }
            _ => todo!("unsupported header version"),
        }
    }
}
