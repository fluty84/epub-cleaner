package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	log	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	fmt.Print("Enter the path to the EPUB file: ")

	epubPath := getUserInput()

	fmt.Print("Enter the string to delete: ")
	deleteString := getUserInput()

	// Create a new EPUB file with the string removed
	err := removeStringFromEPUB(epubPath, deleteString)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("EPUB file created successfully!")
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func removeStringFromEPUB(epubPath, deleteString string) error {
	originalEPUB, err := zip.OpenReader(epubPath)
	if err != nil {
		return err
	}
	defer originalEPUB.Close()

	// Create a new EPUB file
	newEPUBPath := strings.TrimSuffix(epubPath, ".epub") + "_cleaned.epub"
	newEPUB, err := os.Create(newEPUBPath)
	if err != nil {
		return err
	}
	defer newEPUB.Close()

	zipWriter := zip.NewWriter(newEPUB)

	numberOfCoincidences := 0
	// Copy each file from the original EPUB to the new EPUB, removing the specified string
	for _, file := range originalEPUB.File {
		// Open the file from the original EPUB
		originalFile, err := file.Open()
		if err != nil {
			return err
		}
		defer originalFile.Close()

		// Create a new file in the new EPUB
		newFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}
		numberOfCoincidences += strings.Count(fileContents(originalFile), deleteString)
		// Copy the contents of the original file to the new file, removing the specified string
		_, err = io.Copy(newFile, strings.NewReader(strings.ReplaceAll(fileContents(originalFile), deleteString, "")))
		if err != nil {
			return err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("The string '%s' was found %d times", deleteString, numberOfCoincidences))

	return nil
}

func fileContents(file io.Reader) string {
	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(content)
}
