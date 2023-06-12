package models

import "errors"

type Location struct {
	ID        uint
	Name      string
	Address   string
	Link      string
	City      string
	Latitude  string
	Longitude string
	Metadata  `gorm:"embedded"`
}

func (l *Location) ToReadableLocation(readableLocation *ReadableLocation) {
	readableMetadata := l.ToReadableMetadata()

	readableLocation.ID = int(l.ID)
	readableLocation.Name = l.Name
	readableLocation.Address = l.Address
	readableLocation.City = l.City
	readableLocation.Latitude = l.Latitude
	readableLocation.Longitude = l.Longitude
	readableLocation.ReadableMetadata = *readableMetadata
}

type ReadableLocation struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	City             string `json:"city"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	ReadableMetadata `json:"metadata"`
}

func (rl *ReadableLocation) Validate() error {
	switch {
	case rl.Name == "":
		return errors.New("invalid name")
	}
	return nil
}

func (rl *ReadableLocation) ToLocationObject(locationObject *Location) {
	locationObject.ID = uint(rl.ID)
	locationObject.Name = rl.Name
	locationObject.Address = rl.Address
	locationObject.City = rl.City
	locationObject.Latitude = rl.Latitude
	locationObject.Longitude = rl.Longitude
}

func ToReadableLocationList(locationModelList []Location, err *CustomError) []ReadableLocation {
	readableLocationList := make([]ReadableLocation, len(locationModelList))

	for i, item := range locationModelList {
		var readableLocation ReadableLocation
		item.ToReadableLocation(&readableLocation)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse location"
			return nil
		}
		readableLocationList[i] = readableLocation
	}

	return readableLocationList
}
