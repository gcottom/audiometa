# AudioMeta qualification and limitations

The local candidate adds the shared editor to the four format libraries and to the SDK. The released v3 `Tag` contract, concrete format types and format entry points remain available.

## Verification

- The registry defines 133 fields. The shared suite verifies writable-field round trips and removals using the same typed samples in ID3v2.2, v2.3 and v2.4, fast-start MP4/M4A, M4B, Ogg Vorbis, Ogg Opus and FLAC. It also checks exact no-op bytes, ordered duplicate custom values, repeated saves and read-only errors. Separate cases cover ID3v1.1 reads and coordinated removals, protected data, malformed input and short outputs. [Generated field matrix](compatibility.md)
- FFmpeg 9.1 decoded sample-identical PCM from each source and edited output in the cross-format transfer test. FFprobe 9.1 returned identical audio codec properties, sample rate, channel count, duration and ordered SHA-256 hashes for all original compressed audio packets. The test starts with canonical values from a saved FLAC and transfers them to every target, without copying the FLAC audio. Legacy v3.0.5 FFmpeg regressions run alongside this suite.
- Mutagen 1.48.1 independently recognized ID3 text, ReplayGain, MusicBrainz UFID, picture bytes, synchronized lyrics, chapters and musician credits. FFprobe and Mutagen recognized MP4 Nero chapters, repeated artists, track pairs, booleans, BPM, ISRC freeform and exact artwork. Mutagen verified Vorbis and Opus tags and full PNG payloads after the comment packet grew past one page. It verified FLAC comments, ISRC, ReplayGain, and native picture data and attributes. ExifTool 13.64 provided additional MP3 and Ogg readbacks. External tools saw fallback entries as private string data.
- The clean consumer used the proposed module versions, a fresh module cache, `GOWORK=off`, and the local Go proxy protocol. It compiled and executed legacy and unified client code with Go 1.18.10 and Go 1.27.1. The candidate mirrors retain the current Go 1.18 module declarations. The actual published dependency artifacts still need to be checked after release approval. [Release plan](release-plan.md)

All six modules passed build, unit tests, race detector, vet, Staticcheck and govulncheck. The vulnerability scan reported zero reachable findings. Imported-package or required-module findings outside the calling paths are recorded in the individual logs.

| Package | Statement coverage |
| --- | ---: |
| AudioMeta SDK | 77.5% unit, 83.4% with integration and race enabled |
| AudioMeta Core | 81.4% |
| MP3Meta | 79.8% |
| MP4Meta | 79.9% |
| OggMeta | 71.2% |
| FLACMeta | 75.4% |

The nine fixture suites produced 1,197 canonical field results, including 1,170 editable field round trips and removals and 27 read-only assertions. Fuzz targets ran for the SDK detector, core wire values and picture blocks, ID3, MP4, Ogg comments and FLAC comments.

## Reproduction

From the workspace root with the six modules and local workspace file present, run the following from each module directory. For the SDK use `audiometa/v3`. Enable CGO and install a compatible C compiler for race runs.

```text
go build ./...
go test -count=1 -cover ./...
go test -race -count=1 ./...
go vet ./...
staticcheck ./...
govulncheck ./...
```

Run `go test -tags integration -count=1 -v ./...` in AudioMeta, OggMeta and MP4Meta with FFmpeg and FFprobe on PATH. For the independent Mutagen tests install Mutagen for Python and set `AUDIOMETA_MUTAGEN_PYTHON` to the interpreter executable. Set `AUDIOMETA_EXIFTOOL` to the ExifTool binary for the optional Ogg readback. MP3 generation and readback commands are in the workspace `evidence/mp3meta` folder. Build the clean consumer using `py audiometa/scripts/verify_consumer.py` or run it with `--go` to select a specific Go executable.

The workspace `evidence/final` directory contains raw command results and coverage profiles. The per-adapter subdirectories or files under `evidence` contain targeted fuzz, benchmarks and independent-tool logs. The SDK `docs/canonical-test-results.jsonl` contains the field-specific pass receipts used by the matrix. Run `go generate ./...` from the v3 module to regenerate the aliases and matrix from these receipts and the sibling core repository.

## Measured resource costs

Windows amd64, Intel Core Ultra 9 275HX, Go 1.27.1. These are Go benchmark allocations per operation, not peak memory. Disk speed, data shapes and the destination writer affect timing.

| Adapter and scenario | Finding |
| --- | --- |
| FLAC, 16 KiB synthetic transport | About 22 μs and 33.6 KB allocated per save |
| FLAC, 64 MiB synthetic transport | About 37.8 ms and 33.6 KB allocated per save |
| FLAC, title edit with 4 MiB untouched synthetic picture payload | About 16 μs and 33.6 KB allocated to a discard writer |
| Ogg, 50 KB fixture | About 3.8 ms and 133 KB allocated per save |
| Ogg, three-minute 3.7 MB generated Opus | About 33 ms and 121 KB allocated per title save |
| MP3, 32 KiB through 40 MiB synthetic transport | About 82 KiB allocated per metadata rewrite |
| MP3, open editor saving 1 MiB and 8 MiB artwork | About 1.1 MB and 8.4 MB allocated |
| MP4, 1–32 MiB synthetic transport | About 29 KB allocated per save across those sizes |
| MP4, 4 MiB native artwork and artwork with supplemental attributes | About 12.7 MB and 60.5 MB allocated, respectively |

Ogg stages the output on disk before committing bytes to the destination. MP3 stages file destinations and streams ordinary writers. FLAC uses a fixed-size frame copy buffer. MP4 canonical saves stream media and retain bounded metadata and offset tables. The legacy MP4 writer still buffers output. Structured private artwork costs more than native artwork. The canvas, encoded artwork and serialization buffers have separate resource costs.

## Format limits

| Format | Scope and protections |
| --- | --- |
| ID3 | Canonical writing supports versions 2.2, 2.3, 2.4 and v2.4 creation on untagged input. v1/v1.1 supplies lower-priority text and track data. Unsupported revisions, unsupported tag flags, selected flagged frames, meaningful absolute chapter offsets, and ambiguous paired fields can reject edits. Unknown frames and non-ID3 tail formats remain in the source. APEv2 is outside the clear scope. |
| MP4 | Edits apply to movie-scope iTunes items. The canonical writer moves validated local 32- and 64-bit chunk offsets. Fragmented files, competing handlers, external item/index references or unrebasable auxiliary offsets can reject a change. Existing QuickTime chapter tracks remain in place and reject chapter edits. Nero chapters and conventional cover images can be paired with exact canonical companion data. The legacy API still uses the single-media-box buffered writer. |
| FLAC | Native PICTURE and comment blocks must fit 24-bit block lengths. Duplicate comment blocks and malformed lengths can reject the operation. STREAMINFO, CUESHEET, SEEKTABLE, APPLICATION and reserved blocks remain intact. If rewriting a malformed comment would drop undeclared trailing data, the editor returns an error. |
| Ogg | Canonical save requires a complete unchained, unmultiplexed Vorbis or Opus logical stream with valid pages. Its comment packet and count limits are 64 MiB and 1,048,576. Later structural faults can cause a save to fail before output even when the leading headers opened. Clears retain OpusTags binary extensions marked for preservation and channel-layout data. The legacy API keeps its single-comment-page and single-rendered-cover write model. |

Conventions can differ across tools. For example, the AudioMeta field registry reads Xiph `DATE` as a recording date and keeps a separate release date, while some taggers write the release date into `DATE`. Dependent values such as R128 gain cannot be translated across playback gain schemes without contextual information. Conventional Xiph chapters, older cover-art aliases and instrument-qualified performer comments can remain opaque when their exact semantics are unavailable. Structured canonical values transfer via the lossless fallback when the physical convention cannot express them.

Each matrix row identifies a preferred carrier and the pass status for a representative semantic value. Native recognition was checked independently for the subset described above. The matrix does not establish behavior for every arbitrary atom, malformed tag, proprietary box or media player.
