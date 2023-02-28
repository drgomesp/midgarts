pub mod grf;

pub trait FromBytes {
    fn from_bytes(bytes: &[u8]) -> Self;
}
