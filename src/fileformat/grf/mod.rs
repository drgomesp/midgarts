use std::collections::HashMap;
use std::convert::TryFrom;
use std::io::{Cursor, Read};
use std::*;

use byteorder::{LittleEndian, ReadBytesExt};

use crate::fileformat::FromBytes;

pub mod entry;
pub mod file;
pub mod header;

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
