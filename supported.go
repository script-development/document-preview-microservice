package main

import (
	"fmt"
	"os"
	"os/exec"
)

func supported() (pdf bool, libreoffice bool) {
	_, err := exec.LookPath("pdftocairo")
	pdf = err == nil
	if pdf {
		fmt.Println("OK pdf: pdftocairo found (from poppler)")
	} else {
		fmt.Println("NO pdf: pdftocairo not found (from poppler)")
	}

	_, err = exec.LookPath("soffice")
	libreoffice = err == nil || len(os.Getenv("SOFFICE_PATH")) > 0
	if libreoffice {
		fmt.Println("OK docx: soffice found (from libreoffice)")
	} else {
		fmt.Println("OK docx: soffice not found (from libreoffice)")
	}

	return
}
