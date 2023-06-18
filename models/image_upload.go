package models

import (
	"errors"
	"fmt"
	"gofit-api/constants"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type UploadImage struct {
	Name  string
	Image *multipart.FileHeader
	Extension string
}

func (ui *UploadImage) VerifyImageExtension() error {
	ui.Extension = path.Ext(ui.Image.Filename)
	switch ui.Extension {
	case ".png":
		return nil
	case ".jpg":
		return nil
	case ".jpeg":
		return nil
	case ".jfif":
		return nil
	case ".pjpeg":
		return nil
	case ".pjp":
		return nil
	}
	return errors.New("invalid file type. " + ui.Image.Filename + " is not an image")
}

func (ui *UploadImage) CopyIMGToAssets() (string, error) {
	// source file
	source, err := ui.Image.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	// destination
	copyDestionation := fmt.Sprintf("%s/%s%s", constants.PROFILE_IMG_DST, ui.Name, ui.Extension)
	dst, err := os.Create(copyDestionation)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy source to destination
	if _, err = io.Copy(dst, source); err != nil {
		return "", err
	}

	return copyDestionation, nil
}