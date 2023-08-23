package main

import (
	"fmt"
	"os"
	"strings"

	mp3mp4tag "github.com/gcottom/mp3-mp4-tag"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 && args[0] == "." {
		args = args[1:]
	}
	if len(args) >= 1 {
		file := args[0]
		mode := "p"
		if len(args) >= 2 {
			mode = strings.ToLower(args[1])
		}

		if mode == "p" || mode == "parse" || mode == "r" || mode == "read" {
			tag, err := mp3mp4tag.OpenTag(file)
			if err != nil {
				panic(err)
			}
			fmt.Printf("File: %s\n", args[0])
			if tag.Artist() != "" {
				fmt.Printf("Artist: %s\n", tag.Artist())
			}
			if tag.AlbumArtist() != "" {
				fmt.Printf("AlbumArtist: %s\n", tag.AlbumArtist())
			}
			if tag.Album() != "" {
				fmt.Printf("Album: %s\n", tag.Album())
			}
			if tag.BPM() != "" {
				fmt.Printf("BPM: %s\n", tag.BPM())
			}
			if tag.Comments() != "" {
				fmt.Printf("Comment: %s\n", tag.Comments())
			}
			if tag.Composer() != "" {
				fmt.Printf("Composer: %s\n", tag.Composer())
			}
			if tag.CopyrightMsg() != "" {
				fmt.Printf("Copyright: %s\n", tag.CopyrightMsg())
			}
			if tag.Date() != "" {
				fmt.Printf("Date: %s\n", tag.Date())
			}
			if tag.EncodedBy() != "" {
				fmt.Printf("EncodedBy: %s\n", tag.EncodedBy())
			}
			if tag.Genre() != "" {
				fmt.Printf("Genre: %s\n", tag.Genre())
			}
			if tag.Language() != "" {
				fmt.Printf("Language: %s\n", tag.Language())
			}
			if tag.Length() != "" {
				fmt.Printf("Length: %s\n", tag.Length())
			}
			if tag.Lyricist() != "" {
				fmt.Printf("Lyricist: %s\n", tag.Lyricist())
			}
			if tag.PartOfSet() != "" {
				fmt.Printf("PartOfSet: %s\n", tag.PartOfSet())
			}
			if tag.Publisher() != "" {
				fmt.Printf("Publisher: %s\n", tag.Publisher())
			}
			if tag.Title() != "" {
				fmt.Printf("Title: %s\n", tag.Title())
			}
			if tag.Year() != "" {
				fmt.Printf("Year: %s\n", tag.Year())
			}
			if len(tag.AdditionalTags()) > 0 {
				fmt.Print("\nUnmapped Tags:\n")
				for key, value := range tag.AdditionalTags() {
					fmt.Printf("%s: %s\n", key, value)
				}
			}
		}
	}
}
