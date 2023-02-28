use std::fmt::Debug;
use std::io::{BufRead, Cursor, Read, Seek, SeekFrom};
use std::{fs, str};

use byteorder::{ByteOrder, LittleEndian, ReadBytesExt};

use bytes::Buf;
use encoding_rs::WINDOWS_1252;
use yazi::*;

use crate::fileformat::FromBytes;

/// GrfEntry represents an individual file entry inside a GRF file.
#[derive(Debug, Default)]
pub struct GrfEntry {
    pub file_name: String,
    pub compressed_size: u32,
    pub compressed_size_aligned: u32,
    pub uncompressed_size: u32,
    pub flags: u8,
    pub offset: u32,
}

impl FromBytes for GrfEntry {
    fn from_bytes(bytes: &[u8]) -> Self {
        let mut reader = Cursor::new(bytes);
        let mut string_decoder = encoding_rs_io::DecodeReaderBytesBuilder::new();

        let mut buf = vec![];
        reader
            .read_until(0u8, &mut buf)
            .expect("should read file name as string");

        let mut string_decoder = string_decoder
            .encoding(Some(WINDOWS_1252))
            .build(Cursor::new(&buf[0..buf.len() - 1]));

        let mut file_name = String::new();
        string_decoder
            .read_to_string(&mut file_name)
            .unwrap_or_else(|_| panic!("failed to decode file name"));

        let compressed_size = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to file compressed size"));

        let compressed_size_aligned = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to file compressed size aligned"));

        let uncompressed_size = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to file uncompressed size"));

        let flags = reader
            .read_u8()
            .unwrap_or_else(|_| panic!("failed to file flags"));

        let offset = reader
            .read_u32::<LittleEndian>()
            .unwrap_or_else(|_| panic!("failed to file offset"));

        GrfEntry {
            file_name: file_name.to_lowercase(),
            compressed_size,
            compressed_size_aligned,
            uncompressed_size,
            flags,
            offset,
        }
    }
}
