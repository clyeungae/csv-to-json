package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rawData := readCSV(getFilePath())
	fileDataArr := processData(rawData)
	saveAsJson(fileDataArr)
}

func getFilePath() string {
	var pwd string
	pwd, _ = os.Getwd()

	var defaultPath = filepath.Join(pwd, "/csv/data.csv")

	_, err := os.Stat(defaultPath)
	if errors.Is(err, os.ErrNotExist) {
		filename := getFileNameFromUser()
		return filepath.Join(pwd, filename)
	}

	return defaultPath
}

func getFileNameFromUser() string {
	var filepath string

	fmt.Println("Enter the file path")
	fmt.Scan(&filepath)

	return filepath
}

func readCSV(path string) [][]string {
	// open csv
	fd, error := os.Open(path)

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println("Open csv successfully")
	defer fd.Close()

	// read csv
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}
	return records
}

type fileData struct {
	fileName  string
	dataEntry []string
}

func processData(rawData [][]string) []fileData {
	var fileDataArr []fileData
	for index, row := range rawData {
		var key string

		for colIndex, col := range row {
			// handle first row as file name
			if index == 0 {
				if colIndex > 0 {
					fileDataArr = append(fileDataArr, fileData{fileName: col, dataEntry: make([]string, len(rawData)-1)})
				}
				continue
			}

			if colIndex == 0 {
				key = quoteString(col) + ":"
				continue
			}

			fileDataArr[colIndex-1].dataEntry[index-1] = key + quoteString(col)
		}
	}

	return fileDataArr
}

func quoteString(rawString string) string {
	return "\"" + rawString + "\""
}

func saveAsJson(fileDataArr []fileData) {
	pwd, _ := os.Getwd()
	var outputDir = "output"

	for _, d := range fileDataArr {
		f, error := os.Create(filepath.Join(pwd, outputDir) + "/" + d.fileName + ".json")
		checkError(error)
		defer f.Close()

		_, writeErr := f.WriteString("{" + strings.Join(d.dataEntry, ", ") + "}")
		checkError(writeErr)
		f.Sync()
	}

	fmt.Println("Finish Output")
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
