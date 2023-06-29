module gagecottom.com/mp3-mp4-tag

go 1.20

require (
	github.com/bogem/id3v2 v1.2.0
	github.com/dhowden/tag v0.0.0
)

require golang.org/x/text v0.3.2 // indirect

replace github.com/dhowden/tag v0.0.0 => github.com/gcottom/tag v0.1.0
