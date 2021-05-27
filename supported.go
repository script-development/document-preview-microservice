package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Supported struct {
	PdfToCairo  bool
	LibreOffice bool
}

func supported() Supported {
	_, err := exec.LookPath("soffice")
	libreoffice := err == nil || len(os.Getenv("SOFFICE_PATH")) > 0
	if libreoffice {
		fmt.Println("OK docx: soffice found (from libreoffice)")
	} else {
		fmt.Println("OK docx: soffice not found (from libreoffice)")
	}

	_, err = exec.LookPath("pdftocairo")
	pdftocairo := err == nil
	if pdftocairo {
		fmt.Println("OK pdf: pdftocairo found (from poppler)")
	} else if libreoffice {
		fmt.Println("OK pdf: soffice found (from libreoffice)")
	} else {
		fmt.Println("NO pdf: pdftocairo not found (from poppler)")
	}

	return Supported{
		PdfToCairo:  pdftocairo,
		LibreOffice: libreoffice,
	}
}
