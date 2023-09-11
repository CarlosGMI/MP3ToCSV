package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
)

type CsvFile struct {
	name string
	path string
}

var rootCmd = &cobra.Command{
	Use:   "mp3s",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")

		if err != nil {
			return err
		}

		name, err := cmd.Flags().GetString("name")

		if err != nil {
			return err
		}

		des, err := cmd.Flags().GetString("target")

		if err != nil {
			return err
		}

		if err = walkPath(dir, name, des); err != nil {
			return err
		}

		return nil
	},
}

func walkPath(dir string, name string, des string) error {
	file, err := createCsv(name, des)

	if err != nil {
		return err
	}

	defer file.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileExt := filepath.Ext(info.Name())

		if fileExt == ".mp3" {
			metadata, err := getMetadata(path)
			writer := csv.NewWriter(file)

			defer writer.Flush()

			if err != nil {
				fmt.Println(fmt.Printf("%s had trouble when retrieving the metadata", info.Name()))
				writer.Write([]string{"", "", "", info.Name()})
			} else {
				writer.Write([]string{metadata.Title(), metadata.Artist(), metadata.Album(), info.Name()})
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getMetadata(filepath string) (tag.Metadata, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	metadata, err := tag.ReadFrom(file)

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func createCsv(name string, path string) (*os.File, error) {
	fileData, err := getFileData(name, path)

	if err != nil {
		return nil, err
	}

	file, err := os.Create(fmt.Sprintf("%s.csv", filepath.Join(fileData.path, fileData.name)))

	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	headers := []string{"Title", "Artists", "Album", "Filename"}

	defer writer.Flush()
	writer.Write(headers)

	return file, nil
}

func getFileData(name string, path string) (CsvFile, error) {
	filename := "mp3metadata.csv"
	home, err := os.UserHomeDir()

	if err != nil {
		return CsvFile{}, err
	}

	// We create the default directory to store the CSV file
	filepath, err := getDefaultDestination(home)

	if err != nil {
		return CsvFile{}, err
	}

	if name != "" {
		filename = strings.Split(name, ".")[0]
	}

	if path != "" {
		filepath = path
	}

	return CsvFile{filename, filepath}, nil
}

func getDefaultDestination(home string) (string, error) {
	filepath := filepath.Join(home, ".mp3script")
	if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
		return "", err
	}

	return filepath, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var dir string
	var filename string
	var des string

	rootCmd.Flags().StringVarP(&dir, "dir", "d", "", "The directory to extract the metadata from")
	rootCmd.Flags().StringVarP(&filename, "name", "n", "", "The name of the generated CSV file")
	rootCmd.Flags().StringVarP(&des, "target", "t", "", "The destination path where you want the CSV file to be located at")
	rootCmd.MarkFlagRequired("dir")
}
