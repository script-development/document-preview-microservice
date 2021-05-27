package main

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"strings"
)

func formFileToImage(serverFile multipart.File, contentType string, supported Supported) (img image.Image, cropAlignTop bool, err error) {
	// Just to be sure we have an extra close here
	defer serverFile.Close()

	switch contentType {
	case "image/webp", "image/png", "image/jpeg", "image/gif":
		// Parse the image
		img, _, err = image.Decode(serverFile)
		return img, false, err
	case "application/pdf", "application/x-pdf", "application/x-bzpdf", "application/x-gzpdf":
		if !supported.PdfToCairo {
			if supported.LibreOffice {
				// libre office can also convert PDFs to images :)
				img, err = convertOfficeToImageFromBytes(serverFile)
				return img, true, err
			}

			return nil, false, fmt.Errorf("unsupported file content type %s", contentType)
		}

		// Create temp file with contents of input
		filename, err := createTempFile(serverFile, "pdf")
		if err != nil {
			return nil, false, err
		}
		defer os.Remove(filename)

		// Convert the file to a image
		img, err = convertPdfToImage(filename)
		return img, true, err
	case
		"application/vnd.oasis.opendocument.text", // .odt
		"text/csv",           // .csv
		"application/msword", // .doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",   // .docx
		"application/vnd.openxmlformats-officedocument.presentationml.presentation", // .pptx
		"application/mspowerpoint", // .ppt
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", // .xlsx
		"application/msexcel": // .xls

		if !supported.LibreOffice {
			return nil, false, fmt.Errorf("unsupported file content type %s", contentType)
		}

		cropAlignTop = true

		img, err = convertOfficeToImageFromBytes(serverFile)
		return img, true, err
	default:
		return nil, false, fmt.Errorf("unsupported file content type %s", contentType)
	}
}

func convertOfficeToImageFromBytes(serverFile multipart.File) (image.Image, error) {
	// Create temp file with contents of input
	filename, err := createTempFile(serverFile, "document")
	if err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	return convertOfficeToImage(filename)
}

func convertOfficeToImage(filename string) (image.Image, error) {
	// convert the office file to a png
	command := "soffice"
	customSofficePath := os.Getenv("SOFFICE_PATH")
	if len(customSofficePath) > 0 {
		command = customSofficePath
	}

	err := exec.Command(command, "--headless", "--invisible", "--convert-to", "png", "--outdir", path.Dir(filename), filename).Run()
	if err != nil {
		return nil, errors.New("unable to read the document")
	}

	// Read the generated file
	outFilename := strings.TrimSuffix(filename, ".document") + ".png"
	defer os.Remove(outFilename)
	outFileBytes, err := os.Open(outFilename)
	if err != nil {
		return nil, errors.New("unable to convert document to image")
	}

	// Parse the image
	img, _, err := image.Decode(outFileBytes)
	outFileBytes.Close()
	if err != nil {
		return nil, errors.New("unsupported file format")
	}

	return img, nil
}

func convertPdfToImage(filename string) (image.Image, error) {
	// Convert the pdf file to a png file
	err := exec.Command("pdftocairo", filename, "-singlefile", "-png", filename).Run()
	if err != nil {
		return nil, errors.New("unable to read pdf file")
	}

	// Read the generated file
	outFilename := filename + ".png"
	defer os.Remove(outFilename)
	outFileBytes, err := os.Open(outFilename)
	if err != nil {
		return nil, errors.New("unable to convert document to image")
	}

	// Parse the image
	img, _, err := image.Decode(outFileBytes)
	outFileBytes.Close()
	if err != nil {
		return nil, errors.New("unsupported file format")
	}

	return img, nil
}

func createTempFile(file multipart.File, extension string) (string, error) {
	inputFileBytes, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		return "", errors.New("unable to read input file")
	}

	tempFile, err := os.CreateTemp("", "*."+extension)
	if err != nil {
		return "", errors.New("unable to convert document to image")
	}

	filename := tempFile.Name()
	tempFile.Write(inputFileBytes)
	tempFile.Close()
	return filename, nil
}
