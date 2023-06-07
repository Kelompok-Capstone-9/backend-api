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
