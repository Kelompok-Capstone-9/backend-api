package models

import "errors"

type ClassTicket struct {
	ID             uint
	UserID         uint
	User           User
	ClassPackageID uint
	ClassPackage   ClassPackage
	Status         ClassTicketStatus `gorm:"type:enum('booked','pending','cancelled');default:'pending'"`
	Metadata       `gorm:"embedded"`
}

func (ct *ClassTicket) ToReadableClassTicket(readableClassTicket *ReadableClassTicket) {
	readableMetadata := ct.Metadata.ToReadableMetadata()

	ct.ClassPackage.ToReadableClassPackage(&readableClassTicket.ClassPackage)
	ct.User.ToReadableUser(&readableClassTicket.User)

	readableClassTicket.ID = int(ct.ID)
	readableClassTicket.Status = string(ct.Status)
	readableClassTicket.ReadableMetadata = *readableMetadata
}

type ReadableClassTicket struct {
	ID               int                  `json:"id"`
	User             ReadableUser         `json:"user"`
	ClassPackage     ReadableClassPackage `json:"class_package"`
	Status           string               `json:"status"`
	ReadableMetadata `json:"metadata"`
}

func (rct *ReadableClassTicket) Validate() error {
	switch {
	case rct.User.ID == 0:
		return errors.New("invalid user id")
	case rct.ClassPackage.ID == 0:
		return errors.New("invalid class package id")
	case rct.Status == "":
		return errors.New("invalid status")
	}

	if rct.Status != "booked" && rct.Status != "pending" && rct.Status != "cancelled" {
		return errors.New("invalid status. status must contain booked or pending or cancelled")
	}

	return nil
}

func (rct *ReadableClassTicket) EditValidate() error {
	allFieldBlank := rct.User.ID == 0 && rct.ClassPackage.ID == 0 && rct.Status == ""
	if allFieldBlank {
		return errors.New("all field is blank. nothing to change")
	}

	if rct.Status != "" {
		if rct.Status != "booked" && rct.Status != "pending" && rct.Status != "cancelled" {
			return errors.New("invalid status. status must contain booked or pending or cancelled")
		}
	}

	return nil
}

func (rct *ReadableClassTicket) ToClassTicketObject(classTicketObject *ClassTicket, err *CustomError) {
	var ClassTicketStatus ClassTicketStatus
	ClassTicketStatus, err.ErrorMessage = GenerateClassTicketStatus(rct.Status)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid status input"
	}

	classTicketObject.ID = uint(rct.ID)
	classTicketObject.UserID = uint(rct.User.ID)
	classTicketObject.User.ID = uint(rct.User.ID)
	classTicketObject.ClassPackageID = uint(rct.ClassPackage.ID)
	classTicketObject.ClassPackage.ID = uint(rct.ClassPackage.ID)
	classTicketObject.Status = ClassTicketStatus
}

func ToReadableClassTicketList(classTicketObjectList []ClassTicket, err *CustomError) []ReadableClassTicket {
	readableClassTicketList := make([]ReadableClassTicket, len(classTicketObjectList))

	for i, item := range classTicketObjectList {
		var readableClassTicket ReadableClassTicket
		item.ToReadableClassTicket(&readableClassTicket)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse class ticket"
			return nil
		}
		readableClassTicketList[i] = readableClassTicket
		readableClassTicketList[i].User.HidePassword()
		if readableClassTicketList[i].Status != string(Booked) {
			readableClassTicketList[i].ClassPackage.Class.HideLink()
		}
	}

	return readableClassTicketList
}
