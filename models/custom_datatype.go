package models

import (
	"errors"
)

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

type ClassTicketStatus string

const (
	Booked    ClassTicketStatus = "booked"
	Pending   ClassTicketStatus = "pending"
	Cancelled ClassTicketStatus = "cancelled"
)

func GenerateClassTicketStatus(classTicketStatusString string) (ClassTicketStatus, error) {
	switch classTicketStatusString {
	case string(Booked):
		return Booked, nil
	case string(Pending):
		return Pending, nil
	case string(Cancelled):
		return Cancelled, nil
	}
	return "", errors.New("invalid class ticket status")
}

type TrainingLevel string

const (
	Beginner     TrainingLevel = "beginner"
	Intermediate TrainingLevel = "intermediate"
	Advance      TrainingLevel = "advance"
)

func GenerateTrainingLevel(trainingLevelString string) (TrainingLevel, error) {
	switch trainingLevelString {
	case string(Beginner):
		return Beginner, nil
	case string(Intermediate):
		return Intermediate, nil
	case string(Advance):
		return Advance, nil
	}
	return "", errors.New("invalid class ticket status")
}

type ProductType string

const (
	MembershipProduct ProductType = "membership"
	ClassProduct      ProductType = "class"
)

func GenerateProductType(productTypeString string) (ProductType, error) {
	switch productTypeString {
	case string(MembershipProduct):
		return MembershipProduct, nil
	case string(ClassProduct):
		return ClassProduct, nil
	}
	return "", errors.New("invalid product type")
}
