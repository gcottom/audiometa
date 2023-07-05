# mp3-mp4-tag


ID3 and MP4 tag reader and writer for go


[![Go Reference](https://pkg.go.dev/badge/github.com/gcottom/mp3-mp4-tag.svg)](https://pkg.go.dev/github.com/gcottom/mp3-mp4-tag)


This library allows you to parse and write ID tags for mp3 and mp4 (m4a, m4b, m4p) files.

Simply use the OpenTag() function with a string file path as the argument and the library will return the IDTag.

You can access all of the fields of the IDTag through the accessor functions.


Fields that can be parsed:

MP3: Artist, AlbumArtist, Album, AlbumArt, Comments, Composer, Genre, Title, Year, BPM, ContentType, CopyrightMessage, Date, EncodedBy, Lyricist, FileType, Language, Length, PartOfSet, and Publisher

MP4: Artist, AlbumARtist, Album, AlbumArt, Comments, Composer, Genre, Title, Year, EncodedBy, and CopyrightMessage

Fields that can be written: 

MP3: Artist, AlbumArtist, Album, AlbumArt, Comments, Composer, Genre, Title, Year, BPM, ContentType, CopyrightMessage, Date, EncodedBy, Lyricist, FileType, Language, Length, PartOfSet, and Publisher

MP4: Artist, AlbumArtist, Album, AlbumArt, Comments, Composer, Genre, Title, Year, and CopyrightMessage
