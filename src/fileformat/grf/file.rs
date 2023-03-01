use std::collections::HashMap;
use std::fmt::Debug;
use std::fs;
use std::io::prelude::*;
use std::io::{BufReader, Cursor, Read, Seek, SeekFrom};

use byteorder::{LittleEndian, ReadBytesExt};
use encoding_rs::WINDOWS_1252;
use yazi::*;

use crate::fileformat::grf::entry::{GrfEntry, ENTRY_HEADER_SIZE};
use crate::fileformat::grf::header::{GrfHeader, HEADER_SIZE};
use crate::fileformat::grf::Version;
use crate::fileformat::{FromBytes, Loader};

/// The GRF file.
#[derive(Debug, Default)]
pub(crate) struct GrfFile {
    /// GRF file raw data.
    pub(crate) data: Vec<u8>,
    /// The GRF header.
    pub(crate) header: GrfHeader,
    /// The GRF entries table.
    pub(crate) entries: HashMap<String, GrfEntry>,
}

impl GrfFile {
    /// The total file count excluding reserved files.
    pub(crate) fn entry_count(&self) -> usize {
        return (self.header.entry_count - self.header.reserved) as usize - 7;
    }

    /// Get an entry by its path.
    pub(crate) fn get_entry(&self, path: &'static str) -> GrfEntry {
        let entry = self.entries.get(path).unwrap();

        let mut reader = Cursor::new(&self.data);
        reader
            .seek(SeekFrom::Start(
                entry.header._offset as u64 + HEADER_SIZE as u64,
            ))
            .expect("should seek to file table");

        let mut compressed = vec![0u8; entry.header._compressed_size_aligned as usize];
        reader
            .read_exact(&mut compressed)
            .expect("should read entry compressed data");
        let (uncompressed, _checksum) = decompress(&compressed, Format::Zlib).unwrap();

        GrfEntry::from_bytes(&uncompressed)
    }
}

impl Loader for GrfFile {
    fn load(path: &'static str) -> GrfFile {
        GrfFile::from_bytes(&fs::read(path).unwrap())
    }
}

impl FromBytes for GrfFile {
    fn from_bytes(bytes: &[u8]) -> Self {
        let header = GrfHeader::from_bytes(bytes);

        match header.version.try_into() {
            Ok(Version::Version200) => {
                let mut f = GrfFile {
                    data: bytes.to_owned(),
                    header,
                    entries: HashMap::new(),
                };

                let mut reader = Cursor::new(&bytes);
                reader
                    .seek(SeekFrom::Start(f.header.entry_table_offset as u64))
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

                f.entries.reserve(f.header.entry_count as usize);
                let mut reader = BufReader::new(decompressed.as_slice());

                for _i in 0..f.entry_count() {
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
