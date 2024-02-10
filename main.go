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
	// Prompt for EPUB path
	fmt.Print("Enter the path to the EPUB file: ")

	epubPath := getUserInput()

	// Prompt for string to delete
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
	// Open the original EPUB file
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

	// Create a zip writer for the new EPUB file
	zipWriter := zip.NewWriter(newEPUB)


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
		numberOfCoincidences := strings.Count(fileContents(originalFile), deleteString)
		log.Info(fmt.Sprintf("The string '%s' was found %d times in the file %s", deleteString, numberOfCoincidences, file.Name))
		// Copy the contents of the original file to the new file, removing the specified string
		_, err = io.Copy(newFile, strings.NewReader(strings.ReplaceAll(fileContents(originalFile), deleteString, "")))
		if err != nil {
			return err
		}
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func fileContents(file io.Reader) string {
	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(content)
}
