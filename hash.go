package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/corona10/goimagehash"
)

func ImageHash(img image.Image) (*goimagehash.ImageHash, error) {
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		return nil, err
	}
	return hash, err
}

func OpenJpeg(filepath string) (*image.Image, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, err := jpeg.Decode(fd)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func OpenPng(filepath string) (*image.Image, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, err := png.Decode(fd)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
