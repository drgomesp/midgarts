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
