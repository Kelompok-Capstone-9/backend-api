package models

type ClassPackage struct {
	ID       uint
	Period   ClassPeriod `gorm:"type:enum('daily','weekly','monthly')"`
	Price    float32
	ClassID  uint
	Class    Class
	Metadata `gorm:"embedded"`
}

func (cp *ClassPackage) ToReadableClassPackage(readableClassPackage *ReadableClassPackage) {
	readableMetadata := cp.Metadata.ToReadableMetadata()
	var readableClass ReadableClass
	var readableLocation ReadableLocation

	cp.Class.ToReadableClass(&readableClass)
	cp.Class.Location.ToReadableLocation(&readableLocation)

	readableClassPackage.ID = int(cp.ID)
	readableClassPackage.Period = string(cp.Period)
	readableClassPackage.Price = cp.Price
	readableClassPackage.Class.ID = int(cp.Class.ID)
	readableClassPackage.Class = readableClass
	readableClassPackage.Class.Location = readableLocation
	readableClassPackage.ReadableMetadata = *readableMetadata
}

type ReadableClassPackage struct {
	ID               int           `json:"id"`
	Period           string        `json:"period"`
	Price            float32       `json:"price"`
	Class            ReadableClass `json:"class"`
	ReadableMetadata `json:"metadata"`
}

func (rcp *ReadableClassPackage) ToClassPackageObject(classPackageObject *ClassPackage, err *CustomError) {
	var classPeriod ClassPeriod
	classPeriod, err.ErrorMessage = GenerateClassPeriod(rcp.Period)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid class period input"
	}

	classPackageObject.ID = uint(rcp.ID)
	classPackageObject.Period = classPeriod
	classPackageObject.Price = rcp.Price
	classPackageObject.ClassID = uint(rcp.Class.ID)
}

func ToReadableClassPackageList(classPackageObjectList []ClassPackage, err *CustomError) []ReadableClassPackage {
	readableClassPacakgeList := make([]ReadableClassPackage, len(classPackageObjectList))

	for i, item := range classPackageObjectList {
		var readableClassPackage ReadableClassPackage
		item.ToReadableClassPackage(&readableClassPackage)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse class package"
			return nil
		}
		readableClassPacakgeList[i] = readableClassPackage
	}

	return readableClassPacakgeList
}
