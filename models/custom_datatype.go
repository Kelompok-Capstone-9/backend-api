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
	switch genderString {
	case string(Pria):
		return Pria, nil
	case string(Wanita):
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
	switch classTypeString {
	case string(Offline):
		return Offline, nil
	case string(Online):
		return Online, nil
	}
	return "", errors.New("invalid class type")
}

type ClassPeriod string

const (
	Daily   ClassPeriod = "daily"
	Weekly  ClassPeriod = "weekly"
	Monthly ClassPeriod = "monthly"
)

func GenerateClassPeriod(classPeriodString string) (ClassPeriod, error) {
	switch classPeriodString {
	case string(Daily):
		return Daily, nil
	case string(Weekly):
		return Weekly, nil
	case string(Monthly):
		return Monthly, nil
	}
	return "", errors.New("invalid class period")
}
