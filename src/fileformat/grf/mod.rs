use std::convert::TryFrom;
use std::*;

/// The GRF entry submodule.
pub(crate) mod entry;
/// The GRF file submodule.
pub mod file;
/// The GRF header submodule.
pub(crate) mod header;

/// The GRF versions.
#[derive(Debug, Default)]
pub(crate) enum Version {
    #[default]
    /// The GRF version 0x200 (512).
    Ox200 = 0x200,
}

impl TryFrom<u32> for Version {
    type Error = ();

    fn try_from(v: u32) -> Result<Self, Self::Error> {
        match v {
            v if v == Version::Ox200 as u32 => Ok(Version::Ox200),
            _ => todo!("invalid header version {}", v),
        }
    }
}
