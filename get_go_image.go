package main

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
)

func formFileToImage(serverFile multipart.File, contentType string, pdfSupported bool) (img image.Image, cropAlignTop bool, err error) {
	// Just to be sure we have an extra close here
	defer serverFile.Close()

	switch contentType {
	case "image/webp", "image/png", "image/jpeg", "image/gif":
		// Parse the image
		img, _, err = image.Decode(serverFile)
		serverFile.Close()
		if err != nil {
			return nil, false, errors.New("unsupported file format")
		}
	case "application/pdf", "application/x-pdf", "application/x-bzpdf", "application/x-gzpdf":
		if !pdfSupported {
			serverFile.Close()
			return nil, false, fmt.Errorf("unsupported file content type %s", contentType)
		}

		cropAlignTop = true

		inputFileBytes, err := ioutil.ReadAll(serverFile)
		serverFile.Close()
		if err != nil {
			return nil, false, errors.New("unable to read input file")
		}

		// Create temp file with contents of input
		tempFile, err := os.CreateTemp("", "*.pdf")
		if err != nil {
			return nil, false, errors.New("unable to confert pdf file")
		}
		filename := tempFile.Name()
		defer os.Remove(filename)
		tempFile.Write(inputFileBytes)
		tempFile.Close()

		// Convert the pdf file to a png file
		err = exec.Command("pdftocairo", filename, "-singlefile", "-png", filename).Run()
		if err != nil {
			return nil, false, errors.New("unable to read pdf file")
		}

		// Read the generated file
		outFilename := filename + ".png"
		outFileBytes, err := os.Open(outFilename)
		if err != nil {
			return nil, false, errors.New("Cannot")
		}

		// Parse the image
		img, _, err = image.Decode(outFileBytes)
		outFileBytes.Close()
		os.Remove(outFilename)
		os.Remove(filename)
		if err != nil {
			return nil, false, errors.New("unsupported file format")
		}
	default:
		serverFile.Close()
		err = fmt.Errorf("unsupported file content type %s", contentType)
		return
	}

	return
}
