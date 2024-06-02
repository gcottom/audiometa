# audiometa


MP3/MP4/FLAC/OGG tag reader and writer for go


[![Go Reference](https://pkg.go.dev/badge/github.com/gcottom/audiometa/v2.svg)](https://pkg.go.dev/github.com/gcottom/audiometa/v2)   [![Go Report Card](https://goreportcard.com/badge/github.com/gcottom/audiometa/v2)](https://goreportcard.com/report/github.com/gcottom/audiometa/v2)   [![Coverage Status](https://coveralls.io/repos/github/gcottom/audiometa/badge.svg?branch=main)](https://coveralls.io/github/gcottom/audiometa?branch=main)


This package enables parsing and writing of ID tags for mp3, mp4 (m4a, m4b, m4p), FLAC, and ogg (Vorbis, OPUS) files.

All fields can be accessed via the provided accessor methods. Only supported fields per file type will return non zero data. The exception to this is that ogg files support a passthrough map. 
By setting kv pairs in the passthrough map, non-standard vorbis comment tags can be written to both ogg Vorbis and ogg OPUS files. 


# Parsable Fields Per FileType 

## MP3
Artist, AlbumArtist, Album, AlbumArt, BPM, ContentType, Comments, Composer, CopyrightMessage, Date, EncodedBy, FileType, Genre, Language, Length, Lyricist, PartOfSet, Publisher, Title, Year

## MP4 & MP4 Types (m4a, m4b, m4p, mp4)
Artist, AlbumArtist, Album, AlbumArt, Comments, Composer, CopyrightMessage, EncodedBy, Genre, Title, Year

## FLAC
Artist, Album, AlbumArt, Date, Genre, Title

## OGG (Vorbis and OPUS within an ogg container)
Artist, AlbumArtist, Album, AlbumArt, Comment, Composer, Copyright, Date, Genre, Title, Publisher, (extended support for custom fields via passthrough map)


# Writable Fields Per FileType

## MP3
Artist, AlbumArtist, Album, AlbumArt, BPM, ContentType, Comments, Composer, CopyrightMessage, Date, EncodedBy, FileType, Genre, Language, Length, Lyricist, PartOfSet, Publisher, Title, Year

## MP4 & MP4 Types (m4a, m4b, m4p, mp4)
Artist, AlbumArtist, Album, Comments, CopyrightMessage, Composer, Genre, Title, Year (AlbumArt is currently not writeable)

## FLAC
Artist, Album, AlbumArt, Date, Genre, Title

## OGG (Vorbis and OPUS within an ogg container)
Artist, AlbumArtist, Album, AlbumArt, Comment, Composer, Copyright, Date, Genre, Title, Publisher, (extended support for custom fields via passthrough map)

# Usage

## Open Tag For Reading From Path
```
tag, err := audiometa.OpenTagFromPath("./my-audio-file.mp3")
if err != nil{
    panic(err)
}

artist := tag.Artist()
album := tag.Album()
title := tag.Title()
```

## Open Tag For Reading From io.ReadSeeker
```
f, err := os.Open("./my-audio-file.mp3")
tag, err := audiometa.Open(f, audiometa.ParseOptions{Format: audiometa.MP3}) //ParseOptions is only required if the file does not have an extension
if err != nil{
    panic(err)
}

artist := tag.Artist()
album := tag.Album()
title := tag.Title()
```

## Update Tag
```
f, err := os.Open("./my-audio-file.mp3")
if err != nil{
    panic(err)
}
tag, err := audiometa.Open(f)
if err != nil{
    panic(err)
}
tag.SetArtist("Beyonce")
err = tag.Save(f)
if err != nil{
    panic(err)
}
```

## Clear All Tags And Save
```
f, err := os.Open("./my-audio-file.mp3")
if err != nil{
    panic(err)
}
tag, err := audiometa.Open(f)
if err != nil{
    panic(err)
}
tag.ClearAllTags()
err = tag.Save(f)
if err != nil{
    panic(err)
}
```