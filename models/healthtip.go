package models

import (
	"errors"
	"time"
)

type Healthtip struct {
	ID          uint
	UserID      uint
	User        User
	Title       string
	Content     string
	Image       string
	PublishDate time.Time `json:"publish_date"`
	Metadata    `gorm:"embedded"`
}

func (ht *Healthtip) ToReadableHealthtip(readableHealthtip *ReadableHealthtip) {
	readableMetadata := ht.Metadata.ToReadableMetadata()

	ht.User.ToReadableUser(&readableHealthtip.User)

	readableHealthtip.ID = int(ht.ID)
	readableHealthtip.ReadableMetadata = *readableMetadata
}

type ReadableHealthtip struct {
	ID               int          `json:"id"`
	User             ReadableUser `json:"user"`
	Title            string       `json:"title"`
	Content          string       `json:"content"`
	Image            string       `json:"image"`
	PublishDate      string       `json:"publish_date"`
	ReadableMetadata `json:"metadata"`
}

func (rht *ReadableHealthtip) Validate() error {
	switch {
	case rht.User.ID == 0:
		return errors.New("invalid user id")
	case rht.Title == "":
		return errors.New("invalid title")
	case rht.Content == "":
		return errors.New("invalid content")
	case rht.Image == "":
		return errors.New("invalid image")
	case rht.PublishDate == "":
		return errors.New("invalid publish date")
	}

	return nil
}

func (rht *ReadableHealthtip) EditValidate() error {
	allFieldBlank := rht.User.ID == 0 && rht.Title == "" && rht.Content == "" && rht.Image == "" && rht.PublishDate == ""
	if allFieldBlank {
		return errors.New("all field is blank. nothing to change")
	}
	return nil
}

func (rht *ReadableHealthtip) ToHealthtipObject(healthtipObject *Healthtip, err *CustomError) {
	healthtipObject.ID = uint(rht.ID)
	healthtipObject.UserID = uint(rht.User.ID)
	healthtipObject.User.ID = uint(rht.User.ID)
	healthtipObject.Title = rht.Title
}

func ToReadableHealthtipList(healthtipObjectList []Healthtip, err *CustomError) []ReadableHealthtip {
	readableHealthtipList := make([]ReadableHealthtip, len(healthtipObjectList))

	for i, item := range healthtipObjectList {
		var readableHealthtip ReadableHealthtip
		item.ToReadableHealthtip(&readableHealthtip)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse healthtip"
			return nil
		}
		readableHealthtipList[i] = readableHealthtip
		readableHealthtipList[i].User.HidePassword()
	}

	return readableHealthtipList
}
