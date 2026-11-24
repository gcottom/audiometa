# audiometa

audiometa is the final piece that ties several projects together with a simple interface. While the code looks short for this specific project, the underlying modules are extensive. 

## File Types
- MP3 (files with an ID3v2 header, or untagged MPEG Layer III files with two compatible frame headers at the start; free-format streams are not auto-detected)
- MP4 (files containing an ftyp header)
- OGG (files containing an oggs header, with vorbis or opus encoding)
- FLAC (files containing a fLaC header)

## Improvements
This version ties together bug fixes in almost every module. Improvements in performance and memory usage along with tons of optimization and testing in each package. Check the CHANGLOG for a comprehensive list of changes. 

## Related Projects
[mp3meta](https://github.com/gcottom/mp3meta)

[mp4meta](https://github.com/gcottom/mp4meta)

[oggmeta](https://github.com/gcottom/oggmeta)

[flacmeta](https://github.com/gcottom/flacmeta)

## Reading and saving

Keep the seekable source open until saving is complete, and use a separate
destination. A tag can hold format-specific state; do not share one tag or
source reader across concurrent operations. Save is not transactional: if it
returns an error, discard the partial output instead of replacing the original.

```go
func retitle(input io.ReadSeeker, output io.Writer, title string) error {
    tag, err := audiometa.OpenTag(input)
    if err != nil {
        return err
    }
    tag.SetTitle(title)
    return tag.Save(output)
}
```

Import `github.com/gcottom/audiometa/v3` and `io` for this example. The Go
module lives in `v3/`. Run the unit tests there with `go test ./...`. For
decoded-audio round trips, install FFmpeg with AAC, libmp3lame, libvorbis and
libopus support, then run `go test -tags=integration ./...`. The integration
tests generate their own media and compare decoded samples before and after
metadata edits.
