package cmd

import (
	"fmet/utils"
	"fmt"
	"github.com/bogem/id3v2"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var mediaRegex string
var mediaFile string

func init() {
	rootCmd.AddCommand(mediaCmd)
	initShow()
	initEdit()

	mediaCmd.PersistentFlags().StringVarP(&mediaRegex, "regex", "r", "", "Regex usage in format '^(.+)\\.pdf$/'. Note: Do not use double quotes.")
	mediaCmd.PersistentFlags().StringVarP(&mediaFile, "file", "f", "", "Name of file or filepath.")
}

var mediaCmd = &cobra.Command{
	Use:              "media",
	Short:            "Working with media file metadata ",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

		// fill file path
		if !utils.IsPath(mediaFile) {
			path, err := utils.CurrentPath()
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			mediaFile = filepath.Join(path, mediaFile)
		}
	},
}

// =====================================================================================================================
// Show subcommand
// =====================================================================================================================

var showMediaCmd = &cobra.Command{
	Use:   "show",
	Short: "Show metadata",
	Run: func(cmd *cobra.Command, args []string) {
		show()
	},
}

func initShow() {
	mediaCmd.AddCommand(showMediaCmd)
}

func show() {
	if len(mediaRegex) > 0 {
		showRegex(mediaRegex)
	} else {
		showFileMetadata(mediaFile)
	}
}

func printFileMetadata(file *id3v2.Tag) {
	fmt.Printf(`
Title -> %s
Artist -> %s
Album -> %s
Year -> %s
Genre -> %s
`,
		file.Title(),
		file.Artist(),
		file.Album(),
		file.Year(),
		file.Genre())

	fmt.Println(file.AllFrames())
}

func showFileMetadata(filePath string) {
	if !utils.FileExists(filePath) {
		fmt.Printf("%s\n", "Resource not found.")
		return
	}
	tagFile, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer func() { _ = tagFile.Close() }()

	printFileMetadata(tagFile)
}

func showRegex(regex string) {
	if verbose {
		fmt.Printf("Regex -> %s\n", mediaRegex)
	}
	fRegex, err := regexp.Compile(regex)
	if err != nil {
		fmt.Printf("%s\n", "Current filenames regex is incorrect.")
		return
	}

	if verbose {
		fmt.Printf("Current working directory -> %s\n", workingDirectory)
	}
	if len(workingDirectory) == 0 {
		workingDirectory, err = utils.CurrentPath()
		if err != nil {
			fmt.Printf("%s\n", "Unknown working directory.")
			os.Exit(1)
		}
	}

	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		fmt.Printf("%s\n", "Unknown working directory.")
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file != nil && fRegex.MatchString(file.Name()) {
			mediaFilePath := path.Join(workingDirectory, file.Name())

			if !utils.FileExists(mediaFilePath) {
				fmt.Printf("%s\n", "Resource not found.")
				return
			}

			if verbose {
				fmt.Printf("File path -> %s\n", mediaFilePath)
			}

			mediaFile, err := id3v2.Open(mediaFilePath, id3v2.Options{
				Parse: true,
			})
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			printFileMetadata(mediaFile)
			_ = mediaFile.Close()
		}
	}
}

// =====================================================================================================================

// =====================================================================================================================
// Edit subcommand
// =====================================================================================================================

type Media struct {
	Title  string
	Artist string
	Album  string
	Year   string
	Genre  string
}

var editMedia Media
var emptyData bool

var editMediaCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit metadata",
	Run: func(cmd *cobra.Command, args []string) {
		edit()
	},
}

func initEdit() {
	mediaCmd.AddCommand(editMediaCmd)
	editMediaCmd.Flags().StringVarP(&editMedia.Title, "title", "t", "", "Title of media. Note: If you use regex param, you may set group(s) from regex.")
	editMediaCmd.Flags().StringVarP(&editMedia.Artist, "artist", "a", "", "Artist of media. Note: If you use regex param, you may set group(s) from regex.")
	editMediaCmd.Flags().StringVarP(&editMedia.Album, "album", "l", "", "Album of media. Note: If you use regex param, you may set group(s) from regex.")
	editMediaCmd.Flags().StringVarP(&editMedia.Year, "year", "y", "", "Year of media. Note: If you use regex param, you may set group(s) from regex.")
	editMediaCmd.Flags().StringVarP(&editMedia.Genre, "genre", "g", "", "Genre of media. Note: If you use regex param, you may set group(s) from regex.")
	editMediaCmd.Flags().BoolVarP(&emptyData, "empty", "e", false, "Empty all data.")
}

func edit() {
	if len(mediaRegex) > 0 {
		editRegexFileMetadata(mediaRegex, &editMedia)
	} else {
		editFileMetadata(mediaFile, &editMedia)
	}
}

func editFileMetadata(filePath string, media *Media) {
	if !utils.FileExists(filePath) {
		fmt.Printf("%s\n", "Resource not found.")
		return
	}
	mediaFile, err := id3v2.Open(filePath, id3v2.Options{
		Parse: true,
	})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if len(media.Title) > 0 {
		mediaFile.SetTitle(media.Title)
	}
	if len(media.Artist) > 0 {
		mediaFile.SetArtist(media.Artist)
	}
	if len(media.Album) > 0 {
		mediaFile.SetAlbum(media.Album)
	}
	if len(media.Year) > 0 {
		mediaFile.SetYear(media.Year)
	}
	if len(media.Genre) > 0 {
		mediaFile.SetGenre(media.Genre)
	}
	if emptyData {
		mediaFile.DeleteAllFrames()
	}

	if !dryRun {
		if err = mediaFile.Save(); err != nil {
			log.Fatal("Error while saving a tag: ", err)
		}
	}

	printFileMetadata(mediaFile)
	defer func() { _ = mediaFile.Close() }()
}

func editRegexFileMetadata(regex string, media *Media) {
	fRegex, err := regexp.Compile(regex)
	if err != nil {
		fmt.Printf("%s\n", "Current filenames regex is incorrect.")
	}

	if len(workingDirectory) == 0 {
		workingDirectory, err = utils.CurrentPath()
		if err != nil {
			fmt.Printf("%s\n", "Unknown working directory.")
			os.Exit(1)
		}
	}

	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		fmt.Printf("%s\n", "Unknown working directory.")
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file != nil && fRegex.MatchString(file.Name()) {
			data := Media{}
			switch true {
			case len(media.Title) > 0:
				data.Title = fRegex.ReplaceAllString(file.Name(), media.Title)
				fallthrough
			case len(media.Artist) > 0:
				data.Artist = fRegex.ReplaceAllString(file.Name(), media.Artist)
				fallthrough
			case len(media.Album) > 0:
				data.Album = fRegex.ReplaceAllString(file.Name(), media.Album)
				fallthrough
			case len(media.Year) > 0:
				data.Year = fRegex.ReplaceAllString(file.Name(), media.Year)
				fallthrough
			case len(media.Genre) > 0:
				data.Genre = fRegex.ReplaceAllString(file.Name(), media.Genre)
			}

			mediaFilePath := path.Join(workingDirectory, file.Name())
			editFileMetadata(mediaFilePath, &data)
		}
	}
}