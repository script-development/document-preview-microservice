package main

import (
	"fmt"
	"os/exec"
)

func supported() (pdf bool) {
	_, err := exec.LookPath("pdftocairo")
	if err == nil {
		pdf = true
		fmt.Println("OK pdf: pdftocairo found (poppler)")
	} else {
		fmt.Println("NO pdf: pdftocairo not found (poppler)")
	}

	return
}
