package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func main() {
	var supportedLock sync.Mutex
	pdfSupported := supported()

	key := getKey()
	if key == "" {
		fmt.Println("No $KEY set, this is required for the authentication")
		os.Exit(1)
	}
	fmt.Println("Bearer access key:", key)

	app := fiber.New()

	app.Use(compress.New())

	app.Get("", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{"status": "ok"})
	})

	app.Post("/api/preview", requireAuth(), func(c *fiber.Ctx) error {
		file, err := c.FormFile("document")
		if err != nil {
			return err
		}
		height, err := reqFormSizeField(c, "height")
		if err != nil {
			return err
		}
		width, err := reqFormSizeField(c, "width")
		if err != nil {
			return err
		}

		contentType := file.Header.Get("Content-Type")
		if len(contentType) == 0 {
			return errors.New("file content type missing")
		}

		serverFile, err := file.Open()
		if err != nil {
			return errors.New("unable to read form file")
		}

		supportedLock.Lock()
		pdfSupportedCopy := pdfSupported
		supportedLock.Unlock()
		img, cropAlignTop, err := formFileToImage(serverFile, contentType, pdfSupportedCopy)
		if err != nil {
			return err
		}
		img = resizeAndCrop(img, height, width, cropAlignTop)

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 50)
		if err != nil {
			return errors.New("unable to create response")
		}

		resBuf := bytes.NewBuffer(nil)
		err = webp.Encode(resBuf, img, options)
		if err != nil {
			return errors.New("unable to crop file")
		}

		c.Response().Header.Set("Content-Type", "image/webp")
		c.Write(resBuf.Bytes())

		return nil
	})

	app.Listen(":3030")
}

func reqFormSizeField(c *fiber.Ctx, field string) (int, error) {
	fieldVal := c.FormValue(field)
	if len(fieldVal) == 0 {
		return 0, fmt.Errorf("size property %s not set", field)
	}

	res64, err := strconv.ParseInt(fieldVal, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("size property %s not a number", field)
	}
	res := int(res64)

	if res > 10_000 {
		return 0, fmt.Errorf("size property %s cannot be greater than 10_000", field)
	}
	if res < 0 {
		return 0, fmt.Errorf("size property %s cannot be less than 0", field)
	}

	res = res / 20 * 20 // Round the value to steps of 20
	if res == 0 {
		res = 20
	}

	return res, nil
}
