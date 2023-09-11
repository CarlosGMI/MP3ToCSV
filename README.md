# MP3ToCSV

This is a simple CLI application made with [Cobra](https://github.com/spf13/cobra) to extract the metadata of MP3 files in a specific directory and add them to a CSV file. The metadata of the mp3 files is being extracted using [tag](https://github.com/dhowden/tag).

## Usage

`go run main.go`

_Windows OS is not supported_

### Required flags

-   `--dir`, `-d` | The path to the directory with the mp3 files

### Optional flags

-   `--name`, `-n` | The desired name of the CSV file
-   `--target`, `-t` | The path to the directory where the CSV file is going to be stored
