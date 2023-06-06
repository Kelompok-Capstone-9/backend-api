package models

type Location struct {
	ID        uint
	Name      string
	Link      string
	City      string
	Latitude  string
	Longitude string
	ClassID   uint
	Metadata  `gorm:"embedded"`
}

func (l *Location) ToReadableLocation(readableLocation *ReadableLocation) {
	readableMetadata := l.ToReadableMetadata()

	readableLocation.ID = int(l.ID)
	readableLocation.Name = l.Name
	readableLocation.Link = l.Link
	readableLocation.City = l.City
	readableLocation.Latitude = l.Latitude
	readableLocation.Longitude = l.Longitude
	readableLocation.ReadableMetadata = *readableMetadata
}

type ReadableLocation struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Link             string `json:"link"`
	City             string `json:"city"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	ReadableMetadata `json:"metadata"`
}

func (rl *ReadableLocation) ToLocationObject(locationObject *Location) {
	locationObject.ID = uint(rl.ID)
	locationObject.Name = rl.Name
	locationObject.Link = rl.Link
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
