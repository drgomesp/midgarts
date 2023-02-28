pub mod grf;

pub trait FromBytes {
    fn from_bytes(bytes: &[u8]) -> Self;
}

pub trait Loader {
    fn load(path: String) -> Self;
}
