# Changelog

## [3.0.0] - 2024-11-5
### audiometa
- Now uses an interface for common tags between file types instead of a struct
- Each format is now under its own project
- Majorly simplifies interacting with metadata

### MP3
- New module wrapper
- Concurrency improvements

### MP4
- No longer corrupts files
- Smaller memory footprint
- Now works with album art

### OGG
- CRC32 checksum is now calculated and checked correctly
- Opus and Vorbis now support cover art
- Opus cover art no longer breaks the stream
- Added a large set of vorbis tags
- No longer corrupts files

### FLAC
- No longer corrupts files 
- Added a large set of vorbis tags
- Major memory improvements as the audio stream is now copied directly to the writer instead of reading it all to memory first