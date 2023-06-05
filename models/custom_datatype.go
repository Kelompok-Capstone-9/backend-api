package models

import "errors"

// Gender product type
type Gender string

// const available value for enum
const (
	Pria   Gender = "pria"
	Wanita Gender = "wanita"
)

func GenerateGenderType(genderString string) (Gender, error) {
	if genderString == string(Pria) {
		return Pria, nil
	} else if genderString == string(Wanita) {
		return Wanita, nil
	}
	return "", errors.New("invalid gender")
}

type ClassType string

const (
	Offline ClassType = "offline"
	Online  ClassType = "online"
)

func GenerateClassType(classTypeString string) (ClassType, error) {
	if classTypeString == string(Offline) {
		return Offline, nil
	} else if classTypeString == string(Online) {
		return Online, nil
	}
	return "", errors.New("invalid class type")
}
