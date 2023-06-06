package models

type ClassPackage struct {
	ID       uint
	Name     string
	Price    float32
	ClassID  uint
	Metadata `gorm:"embedded"`
}

func (cp *ClassPackage) ToReadableClassPackage(readableClassPackage *ReadableClassPackage) {
	readableMetadata := cp.Metadata.ToReadableMetadata()

	readableClassPackage.ID = int(cp.ID)
	readableClassPackage.Name = cp.Name
	readableClassPackage.Price = int(cp.Price)
	readableClassPackage.ClassID = int(cp.ClassID)
	readableClassPackage.ReadableMetadata = *readableMetadata
}

type ReadableClassPackage struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
	ClassID          int    `json:"class_id"`
	ReadableMetadata `json:"metadata"`
}

func (rcp *ReadableClassPackage) ToClassPackageObject(classPackageObject *ClassPackage) {
	classPackageObject.ID = uint(rcp.ID)
	classPackageObject.Name = rcp.Name
	classPackageObject.Price = float32(rcp.Price)
	classPackageObject.ClassID = uint(rcp.ClassID)
}
