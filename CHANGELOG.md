# Changelog

## [3.0.4]
### audiometa
- Use latest dependencies

### MP3
- Minor change to temp writer
- Update dependencies

### MP4
- Update dependencies

### OGG
- Update dependencies

### FLAC
- Read/Write Circle protection
- Update dependencies

## [3.0.1]-[3.0.3] - 2024-12-29
### audiometa
- Use latest dependency versions

### MP3
- Update README
- Merge PR from LaptopCat to remove unecessary logging messages

### MP4
- Update README
- Merge PR from LaptopCat to remove uncessary logging messages

### OGG
- Update README

### FLAC
- Update README


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