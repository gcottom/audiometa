package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gcottom/audiometa"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 && args[0] == "." {
		args = args[1:]
	}
	if len(args) >= 1 {
		file := ""
		if len(args) >= 2 {
			file = args[1]
		}
		mode := "help"
		mode = strings.ToLower(args[0])
		if len(args)%2 == 0 && len(args) != 1 {
			if mode == "p" || mode == "parse" || mode == "r" || mode == "read" || mode == "-p" || mode == "-parse" || mode == "-r" || mode == "-read" {
				tag, err := audiometa.OpenTagFromPath(file)
				if err != nil {
					panic(err)
				}
				printedTags := 0
				fmt.Printf("File: %s\n", file)
				if tag.Artist() != "" {
					fmt.Printf("Artist: %s\n", tag.Artist())
					printedTags++
				}
				if tag.AlbumArtist() != "" {
					fmt.Printf("AlbumArtist: %s\n", tag.AlbumArtist())
					printedTags++
				}
				if tag.Album() != "" {
					fmt.Printf("Album: %s\n", tag.Album())
					printedTags++
				}
				if tag.BPM() != "" {
					fmt.Printf("BPM: %s\n", tag.BPM())
					printedTags++
				}
				if tag.Comments() != "" {
					fmt.Printf("Comment: %s\n", tag.Comments())
					printedTags++
				}
				if tag.Composer() != "" {
					fmt.Printf("Composer: %s\n", tag.Composer())
					printedTags++
				}
				if tag.CopyrightMsg() != "" {
					fmt.Printf("Copyright: %s\n", tag.CopyrightMsg())
					printedTags++
				}
				if tag.Date() != "" {
					fmt.Printf("Date: %s\n", tag.Date())
					printedTags++
				}
				if tag.EncodedBy() != "" {
					fmt.Printf("EncodedBy: %s\n", tag.EncodedBy())
					printedTags++
				}
				if tag.Genre() != "" {
					fmt.Printf("Genre: %s\n", tag.Genre())
					printedTags++
				}
				if tag.Language() != "" {
					fmt.Printf("Language: %s\n", tag.Language())
					printedTags++
				}
				if tag.Length() != "" {
					fmt.Printf("Length: %s\n", tag.Length())
					printedTags++
				}
				if tag.Lyricist() != "" {
					fmt.Printf("Lyricist: %s\n", tag.Lyricist())
					printedTags++
				}
				if tag.PartOfSet() != "" {
					fmt.Printf("PartOfSet: %s\n", tag.PartOfSet())
					printedTags++
				}
				if tag.Publisher() != "" {
					fmt.Printf("Publisher: %s\n", tag.Publisher())
					printedTags++
				}
				if tag.Title() != "" {
					fmt.Printf("Title: %s\n", tag.Title())
					printedTags++
				}
				if tag.Year() != "" {
					fmt.Printf("Year: %s\n", tag.Year())
					printedTags++
				}
				if len(tag.AdditionalTags()) > 0 {
					fmt.Print("\nUnmapped Tags:\n")
					for key, value := range tag.AdditionalTags() {
						fmt.Printf("%s: %s\n", key, value)
						printedTags++
					}
				}
				if printedTags == 0 {
					fmt.Println("No tags found for this file!")
				}
			} else if mode == "s" || mode == "w" || mode == "save" || mode == "write" || mode == "-s" || mode == "-w" || mode == "-save" || mode == "-write" {
				if len(args) > 2 {
					tag, err := audiometa.OpenTagFromPath(file)
					if err != nil {
						panic(err)
					}
					for i := 2; i < len(args); i += 2 {
						cmdTag := strings.ToLower(args[i])
						writeTag := strings.ToLower(args[i+1])
						if cmdTag == "art" || cmdTag == "-art" || cmdTag == "artist" || cmdTag == "-artist" {
							tag.SetArtist(writeTag)
						} else if cmdTag == "aa" || cmdTag == "-aa" || cmdTag == "-albumartist" || cmdTag == "albumartist" {
							tag.SetAlbumArtist(writeTag)
						} else if cmdTag == "alb" || cmdTag == "-alb" || cmdTag == "album" || cmdTag == "-album" {
							tag.SetAlbum(writeTag)
						} else if cmdTag == "c" || cmdTag == "-c" || cmdTag == "cover" || cmdTag == "-cover" {
							tag.SetAlbumArtFromFilePath(writeTag)
						} else if cmdTag == "comment" || cmdTag == "-comment" || cmdTag == "comments" || cmdTag == "-comments" {
							tag.SetComments(writeTag)
						} else if cmdTag == "composer" || cmdTag == "-composer" {
							tag.SetComposer(writeTag)
						} else if cmdTag == "g" || cmdTag == "-g" || cmdTag == "genre" || cmdTag == "-genre" {
							tag.SetGenre(writeTag)
						} else if cmdTag == "t" || cmdTag == "-t" || cmdTag == "title" || cmdTag == "-title" {
							tag.SetTitle(writeTag)
						} else if cmdTag == "y" || cmdTag == "-y" || cmdTag == "-year" || cmdTag == "year" {
							tag.SetYear(writeTag)
						} else if cmdTag == "b" || cmdTag == "-b" || cmdTag == "bpm" || cmdTag == "-bpm" {
							tag.SetBPM(writeTag)
						} else {
							fileType, err := audiometa.GetFileType(args[1])
							if err != nil {
								fmt.Println(err)
							}
							if fileType == "ogg" {
								cmdTag = strings.TrimPrefix(cmdTag, "-")
								tag.SetAdditionalTag(strings.ToUpper(cmdTag), writeTag)
							} else {
								fmt.Printf("Unsupported tag: %s\n", cmdTag)
							}
						}
					}
					if err = tag.Save(); err != nil {
						fmt.Println(err)
					}
				}

			} else if mode == "c" || mode == "clear" || mode == "e" || mode == "empty" || mode == "-c" || mode == "-clear" || mode == "-e" || mode == "-empty" {
				tag, err := audiometa.OpenTagFromPath(file)
				if err != nil {
					panic(err)
				}
				tag.ClearAllTags(false)
				err = tag.Save()
				if err != nil {
					fmt.Println(err)
				}

			} else if mode == "h" || mode == "-h" || mode == "-help" || mode == "help" {
				if args[1] == "r" || args[1] == "p" || args[1] == "parse" || args[1] == "read" {
					fmt.Println("mp3-mp4-tag-cmd-help: Read(Parse) mode\nThe parse mode shows all tags that are attached to the file. If the tag is not mapped with this application it will be shown in the \"unmapped tags\" section at the bottom of the output.\nex usage: mp3-mp4-tag parse filepath.mp3")
				} else if args[1] == "w" || args[1] == "s" || args[1] == "write" || args[1] == "save" || args[1] == "-w" || args[1] == "-s" || args[1] == "-write" || args[1] == "-save" {
					fmt.Println("mp3-mp4-tag-cmd-help: Write(Save) mode\nThe save mode allows you to save specified tags to a file. The following flags can be used:\n\"art\"= artist\n\"aa\"= albumartist\n\"alb\"= album\n\"c\"= cover (this sets the cover art from a filepath)\n\"comment\"= comments\n\"composer\"= composer\n\"g\"= genre\n\"t\"= title\n\"y\"= year\n\"b\"= bpm\n\nIf the filetype is ogg you can specify additional tags by specifying the tag name as a flag\nex usage: mp3-mp4-tag write filepath.ogg mycustomflag value\nAdditionally every tag can be specified by its full name\nex usage: mp3-mp4-tag write filepath.mp3 artist artistname\nstandard ex usage: mp3-mp4-tag write filepath.mp3 art artistname\nYou can specify 1 or more tag pairs as long as they are complete pairs. Make sure to enclose values with spaces in them within \"quotes\"")
				} else if args[1] == "c" || args[1] == "e" || args[1] == "clear" || args[1] == "empty" || args[1] == "-c" || args[1] == "-e" || args[1] == "-clear" || args[1] == "-empty" {
					fmt.Println("mp3-mp4-tag-cmd-help: Clear(Empty) mode\nThe clear mode clears all of the tags in the file and saves the file with empty tags. You can use this mode to clear all tags of a file and then use the write mode to write all new tags to the file.\nex usage: mp3-mp4-tag clear filepath.mp3")
				}
			} else {
				if len(args) == 1 && args[0] == "h" || args[0] == "-h" || args[0] == "-help" || args[0] == "help" {
					fmt.Println("mp3-mp4-tag-cmd-help: The application offers 3 modes: read, write, and clear. Use the command \"help\" followed by a mode to learn more about its usage.")
				} else {
					fmt.Println("Invalid number of arguments!\nmp3-mp4-tag-cmd-help: The application offers 3 modes: read, write, and clear. Use the command \"help\" followed by a mode to learn more about its usage.")
				}

			}
		}
	}
}
