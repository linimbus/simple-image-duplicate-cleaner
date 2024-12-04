package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
	"github.com/astaxie/beego/logs"
	"github.com/corona10/goimagehash"
	"golang.org/x/image/bmp"
)

const (
	IMG_PNG  string = "PNG"
	IMG_JPEG string = "JPEG"
	IMG_BMP  string = "BMP"
	IMG_HEIC string = "HEIC"
)

type ImageInfo struct {
	hash *goimagehash.ImageHash
	size ImageSize
	file string
	name string
}

var imageFormatList_ map[string]string

func init() {
	imageFormatList_ = make(map[string]string)
	imageFormatList_[".png"] = IMG_PNG
	imageFormatList_[".jpeg"] = IMG_JPEG
	imageFormatList_[".jpg"] = IMG_JPEG
	imageFormatList_[".bmp"] = IMG_BMP
	imageFormatList_[".heic"] = IMG_HEIC
	imageFormatList_[".heif"] = IMG_HEIC
}

func ImageSimilarity(hash1, hash2 *goimagehash.ImageHash) (float64, error) {
	distance, err := hash1.Distance(hash2)
	if err != nil {
		return 0.0, err
	}
	hashLen := 64
	similarity := (1 - float64(distance)/float64(hashLen)) * 100
	return similarity, nil
}

func ImageHash(img image.Image) (*goimagehash.ImageHash, error) {
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		return nil, err
	}
	logs.Info("image hash %s", hash.ToString())
	return hash, nil
}

func ImageLoad(filepath string, name string) (image.Image, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	switch name {
	case IMG_PNG:
		return png.Decode(fd)
	case IMG_JPEG:
		return jpeg.Decode(fd)
	case IMG_BMP:
		return bmp.Decode(fd)
	case IMG_HEIC:
		return goheif.Decode(fd)
	default:
		return nil, fmt.Errorf("unkown image type")
	}
}

func ImageConfigLoad(filepath string, name string) (*image.Config, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var cfg image.Config

	switch name {
	case IMG_PNG:
		cfg, err = png.DecodeConfig(fd)
	case IMG_JPEG:
		cfg, err = jpeg.DecodeConfig(fd)
	case IMG_BMP:
		cfg, err = bmp.DecodeConfig(fd)
	case IMG_HEIC:
		cfg, err = goheif.DecodeConfig(fd)
	default:
		return nil, fmt.Errorf("unkown image type")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ImageOpen(file string, flags map[string]bool) (*ImageInfo, error) {
	suffix := strings.ToLower(filepath.Ext(file))
	name, b := imageFormatList_[suffix]
	if !b {
		return nil, fmt.Errorf("not image file")
	}
	flag, b := flags[name]
	if !b || !flag {
		return nil, fmt.Errorf("ignore image file")
	}
	img, err := ImageLoad(file, name)
	if err != nil {
		return nil, err
	}
	hash, err := ImageHash(img)
	if err != nil {
		return nil, err
	}
	cfg, err := ImageConfigLoad(file, name)
	if err != nil {
		return nil, err
	}
	return &ImageInfo{file: file, name: name, size: ImageSize{Width: cfg.Width, Height: cfg.Height}, hash: hash}, nil
}
