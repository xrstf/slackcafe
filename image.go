package main

import (
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
)

func prepareImage(body io.Reader, target string) error {
	i, err := png.Decode(body)
	if err != nil {
		return fmt.Errorf("failed to decode image as PNG: %v", err)
	}

	bounds := i.Bounds()

	x := 35
	y := 170
	width := bounds.Max.X - 35
	height := bounds.Max.Y - 150

	img := imaging.Crop(i, image.Rectangle{
		Min: image.Pt(x, y),
		Max: image.Pt(width, height),
	})

	err = imaging.Save(img, target)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}
