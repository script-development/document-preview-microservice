package main

import (
	"image"

	"github.com/anthonynsimon/bild/transform"
)

func resizeAndCrop(img image.Image, height, width int, cropAlignTop bool) image.Image {
	imgSize := img.Bounds().Max
	if imgSize.Y < height {
		height = imgSize.Y
	}
	if imgSize.X < width {
		width = imgSize.X
	}

	heightFactor := float64(imgSize.Y) / float64(height)
	widthFactor := float64(imgSize.X) / float64(width)
	if int(heightFactor*100) == int(widthFactor*100) {
		// Do nothing aspect ratio equal
	} else if heightFactor > widthFactor {
		// Crop in height
		newHeight := widthFactor * float64(height)
		offset := float64(0)
		if !cropAlignTop {
			offset = (float64(imgSize.Y) - newHeight) / 2
		}

		img = transform.Crop(img, image.Rect(0, int(offset), imgSize.X, int(offset+newHeight)))
	} else {
		// Crop in width
		newWidth := heightFactor * float64(width)
		offset := (float64(imgSize.X) - newWidth) / 2

		img = transform.Crop(img, image.Rect(int(offset), 0, int(offset+newWidth), imgSize.Y))
	}

	if imgSize.Y > height || imgSize.X > width {
		img = transform.Resize(img, width, height, transform.CatmullRom)
	}

	return img
}
