package assetsmanager

import (
	"fmt"
	"gofit-api/constants"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InitAssetsFile() error {
	isAssetsExists, err := exists(constants.PROFILE_IMG_DST)
	if isAssetsExists {
		fmt.Println("Assets folder exists")
		return nil
	}
	fmt.Println("Assets folder is not exists: ", err)

	fmt.Println("Creating assets folder...")
	if err := os.MkdirAll(constants.PROFILE_IMG_DST, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(constants.CLASS_IMG_DST, os.ModePerm); err != nil {
		return err
	}
	fmt.Println("Assets folder created!")

	return nil
}
