# Phase 0 baseline

Fresh local clones from the public upstream defaults, checked 2026-11-24. The supplied Misc directory was absent. The working set is under `C:\Users\gagec\AudioMeta-Unified`.

| Repository | Default and release commit | Latest tag | Other remote branches |
| --- | --- | --- | --- |
| audiometa | `eaf7bad4050011f84867068267d61b5433555f2e` | v3.0.5 | dev, fix/audio-metadata-round-trips |
| mp3meta | `e15861211477035d0fdd95a6344834fd5c838608` | v0.0.5 | fix/audio-metadata-round-trips |
| mp4meta | `fc98bd5d1461b902cf620a8f1f91b3bb5c8d1f` | v0.0.6 | fix/audio-metadata-round-trips |
| oggmeta | `2529922cf68be9f89c3016dbc59f26001d8c1a9b` | v0.0.9 | fix/audio-metadata-round-trips |
| flacmeta | `aa84ede0cc4fedd567cf4046e238d75c7bf0463a` | v0.0.7 | fix/audio-metadata-round-trips |

All clones started clean. The SDK module starts in the `v3` directory. Each baseline module targets Go 1.18. The SDK pins exactly the four adapter tags above. MP3 uses bogem/id3v2 v2.1.4. MP4 uses abema/go-mp4 v1.3.0, bufseekio v0.1.0 and writerseeker v1.1.0. Ogg uses writerseeker v1.1.0. FLAC's runtime uses the standard library.

## Existing controls and gaps

- SDK v3.0.5 recognizes three untagged MPEG Layer III variants and rejects false-positive MPEG headers. Its existing optional FFmpeg integration verifies decoded media following legacy edits. The v3 Tag interface has a small fixed field set. Adding methods would break external implementations. [Release](https://github.com/gcottom/audiometa/releases/tag/v3.0.5)
- MP3 v0.0.5 calculates the declared ID3 span separately from decoder state and preserves bytes after the header, padding and footer. The ID3 dependency normalizes frames and lacks v2.2 support. The legacy writer rewrites managed frames and re-encodes selected artwork.
- MP4 v0.0.6 relocates fast-start chunk offsets when metadata changes. The legacy writer reconstructs the iTunes item list, drops unrecognized items, re-encodes artwork and buffers the output. Existing protections around media offsets and invalid writes must remain covered.
- Ogg v0.0.9 retains unrecognized comments and audio pages. The legacy projection collapses selected repeated fields. Rebuilding the packet can lose OpusTags trailer bytes and the writer stages encoded output in memory.
- FLAC v0.0.7 streams compressed audio and preserves unmanaged blocks. Its legacy display fields use scalar getters. The lower-level tag retains raw comments, while the exported standalone Vorbis comment map drops repeated keys.

## Verification on base revisions

Go 1.27.1, Windows amd64. The first pass ran `go build ./...`, `go test -count=1 -cover ./...` and `go vet ./...` from each module root.

| Module | Build | Tests and coverage | Vet |
| --- | --- | --- | --- |
| AudioMeta v3 | Pass | Pass, 69.0%, unit tests | Pass |
| MP3Meta | Pass | Fail, Windows fixture cleanup while a handle was open, 89.8% | Pass |
| MP4Meta | Pass | Pass, 77.3% | Pass |
| OggMeta | Pass | Pass, 65.2% | Pass |
| FLACMeta | Pass | Pass, 73.9% | Pass |

The initial race invocations were blocked by CGO being disabled. Compiler-enabled, FFmpeg and analyzer baselines use pristine snapshots of the commits listed above.

## Historical design

The August 2023 changes exposed `AdditionalTags`, `SetAdditionalTag`, a map-based Ogg Vorbis/Opus passthrough, and a clear operation with optional unknown preservation. Added tests established the intent to round-trip arbitrary metadata. A single-value map does not represent repeated comments, ordering or ownership. The v3 split removed this API from the unified surface. [Extended-tag change](https://github.com/gcottom/audiometa/commit/cb78d966eeb7481ab8c6f77280c63a6a0b72241d)

With GCC 16.2.0 staged and CGO enabled, `go test -race ./...` passed for AudioMeta, MP4Meta, OggMeta and FLACMeta. MP3Meta reproduced the Win32 fixture failure from a pristine detached copy. Staticcheck 2026.2.2 and govulncheck v1.10.0 passed all five modules at the baseline revisions, with zero reachable vulnerabilities reported and a database updated 2026-11-23. FFmpeg 9.1 passed the tagged MP3, untagged MP3, M4A, M4B, Ogg Vorbis, Opus and FLAC baseline decode and SHA-256 checks. Raw logs are in the workspace sibling `evidence/phase0`.

Baseline FFmpeg reproduction from `audiometa/v3` with FFmpeg and FFprobe on PATH: `go test -tags integration -v -count=1 ./...`.
