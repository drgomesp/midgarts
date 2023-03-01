use std::fmt::Debug;
use std::io::{BufRead, Cursor, Read, Seek, SeekFrom};
use std::{fs, str};

use byteorder::{ByteOrder, LittleEndian, ReadBytesExt};
use bytes::Buf;
use encoding_rs::WINDOWS_1252;
use yazi::*;

use crate::fileformat::FromBytes;

/// The GRF entry header size constant.
pub const ENTRY_HEADER_SIZE: usize = 17;

/// The encryption mode of the GRF entry.
#[derive(Debug, Default)]
pub enum Encryption {
    #[default]
    /// No encryption.
    None = 0x01,
    /// Mixed encryption.
    Mixed = 0x02,
    /// Header-only encryption.
    Header = 0x04,
}

/// GrfEntryHeader is the entry header of a given entry in a GRF file.
#[derive(Debug, Default)]
pub struct GrfEntryHeader {
    /// Compressed size in bytes.
    pub compressed_size: u32,
    /// Compressed size aligned in bytes.
    pub compressed_size_aligned: u32,
    /// Uncompressed size.
    pub uncompressed_size: u32,
    /// Flags
    pub flags: u8,
    /// Offset
    pub offset: u32,
}

/// GrfEntry represents an individual file entry inside a GRF file.
#[derive(Debug, Default)]
pub struct GrfEntry {
    /// The entry raw data.
    pub data: Vec<u8>,
    /// File name.
    pub file_name: String,
    /// Entry header.
    pub header: GrfEntryHeader,
}

impl FromBytes for GrfEntry {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);

        let compressed_size = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to read file compressed size"));

        let compressed_size_aligned = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to read file compressed size aligned"));

        let uncompressed_size = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to read file uncompressed size"));

        let flags = reader
            .read_u8()
            .unwrap_or_else(|_| panic!("failed to read file flags"));

        let offset = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to read file offset"));

        GrfEntry {
            data: bytes.to_vec(),
            file_name: "".to_string(),
            header: GrfEntryHeader {
                compressed_size,
                compressed_size_aligned,
                uncompressed_size,
                flags,
                offset,
            },
        }
    }
}
