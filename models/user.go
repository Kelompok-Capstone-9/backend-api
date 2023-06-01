package models

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User Object for gorm
type User struct {
	ID         uint
	Name       string
	Email      string `gorm:"unique"`
	Password   string
	Gender     Gender `gorm:"type:enum('pria','wanita')"`
	Height     float32
	GoalHeight float32
	Weight     float32
	GoalWeight float32
	IsAdmin    bool
	Metadata   `gorm:"embedded"`
}

func (u *User) InsertID(itemIDString string, err *CustomError) {
	var itemID int
	itemID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater"
	}
	u.ID = uint(itemID)
}

func (u *User) HashingPassword(err *CustomError) {
	var passwordInBytes []byte
	passwordInBytes, err.ErrorMessage = bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if !err.IsError() {
		u.Password = string(passwordInBytes)
	}
}

func (u *User) MatchingPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) ToReadableUser(readableUser *ReadableUser) {
	readableMetadata := u.ToReadableMetadata()
	readableUser.ID = int(u.ID)
	readableUser.Name = u.Name
	readableUser.Email = u.Email
	readableUser.Password = u.Password
	readableUser.Gender = string(u.Gender)
	readableUser.Height = u.Height
	readableUser.GoalHeight = u.GoalHeight
	readableUser.Weight = u.Weight
	readableUser.GoalWeight = u.GoalWeight
	readableUser.ReadableMetadata = *readableMetadata
}

// User Data or Readable data
type ReadableUser struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	Gender           string  `json:"gender"`
	Height           float32 `json:"height"`
	GoalHeight       float32 `json:"goal_height"`
	Weight           float32 `json:"weight"`
	GoalWeight       float32 `json:"goal_weight"`
	ReadableMetadata `json:"metadata"`
}

// convert id string to int
func (ru *ReadableUser) InsertID(itemIDString string, err *CustomError) {
	ru.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater : " + itemIDString
	}
}

func (ru *ReadableUser) HidePassword() {
	ru.Password = "********"
}

func (ru *ReadableUser) ToUserObject(userObject *User, err *CustomError) {
	// metadata := ru.ToMetadata(err)
	// if err.ErrorMessage != nil {
	// 	return User{}
	// }

	var userGender Gender
	userGender, err.ErrorMessage = GenerateGenderType(ru.Gender)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid gender input"
	}

	userObject.ID = uint(ru.ID)
	userObject.Name = ru.Name
	userObject.Email = ru.Email
	userObject.Password = ru.Password
	userObject.Gender = userGender
	userObject.Height = ru.Height
	userObject.GoalHeight = ru.GoalHeight
	userObject.Weight = ru.Weight
	userObject.GoalWeight = ru.GoalWeight
	// userObject.Metadata = *metadata
}

func ToReadableUserList(userModelList []User, err *CustomError) []ReadableUser {
	readableUserList := make([]ReadableUser, len(userModelList))

	for i, item := range userModelList {
		var readableUser ReadableUser
		item.ToReadableUser(&readableUser)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse user"
			return nil
		}
		readableUserList[i] = readableUser
		readableUserList[i].HidePassword()
	}

	return readableUserList
}
