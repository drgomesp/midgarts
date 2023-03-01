/// The GRF module implements GRF file format.
pub mod grf;

/// The SPR module implements the sprite file format (.SPR).
pub mod spr;

/// The FromBytes trait defines a way for decoding structs from byte slices.
pub trait FromBytes {
    /// Decode a slice of bytes into Self.
    fn from_bytes(bytes: &[u8]) -> Self;
}

/// The Loader trait defines the concept of loadable files by path.  
///
/// # Examples
///     
/// Load a GRF file from a given path:
/// ```rust
/// GRF::load("assets/data.grf")
/// ```
///
pub trait Loader {
    /// Loads Self from a given path.
    fn load(path: String) -> Self;
}
