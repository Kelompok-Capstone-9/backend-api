package models

import (
	"errors"
	"net/mail"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User Object for gorm
type User struct {
	ID             uint
	Name           string
	Email          string `gorm:"unique"`
	Password       string
	Gender         Gender `gorm:"type:enum('pria','wanita')"`
	Height         float32
	GoalHeight     float32
	Weight         float32
	GoalWeight     float32
	TrainingLevel  TrainingLevel `gorm:"type:enum('beginner','intermediate','advance');default:beginner"`
	ProfilePicture string
	IsAdmin        bool
	OTP            int `gorm:"size:4"`
	ClassTickets   []ClassTicket
	Metadata       `gorm:"embedded"`
}

func (u *User) InsertID(userIDString string, err *CustomError) {
	var userID int
	userID, err.ErrorMessage = strconv.Atoi(userIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater"
	}
	u.ID = uint(userID)
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
	readableUser.TrainingLevel = string(u.TrainingLevel)
	readableUser.ProfilePicture = u.ProfilePicture
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
	TrainingLevel    string  `json:"training_level"`
	ProfilePicture   string  `json:"profile_picture"`
	ReadableMetadata `json:"metadata"`
}

// convert id string to int
func (ru *ReadableUser) InsertID(userIDString string, err *CustomError) {
	ru.ID, err.ErrorMessage = strconv.Atoi(userIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater : " + userIDString
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

	var userTrainingLevel TrainingLevel
	userTrainingLevel, err.ErrorMessage = GenerateTrainingLevel(ru.TrainingLevel)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid training level input"
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
	userObject.TrainingLevel = userTrainingLevel
	userObject.ProfilePicture = ru.ProfilePicture
	// userObject.Metadata = *metadata
}

func (ru *ReadableUser) EditValidate() error {
	allFieldBlank := ru.Name == "" && ru.Email == "" && ru.Password == "" && ru.Gender == "" && ru.Height == 0 && ru.GoalHeight == 0 && ru.Weight == 0 && ru.GoalWeight == 0 && ru.TrainingLevel == "" && ru.ProfilePicture == ""
	if allFieldBlank {
		return errors.New("all field is blank. nothing to change")
	}

	if ru.Email != "" {
		_, emailError := mail.ParseAddress(ru.Email)
		if emailError != nil {
			return errors.New("invalid email " + ru.Email)
		}
	}

	if ru.Gender != "" {
		if ru.Gender != "pria" && ru.Gender != "wanita" {
			return errors.New("invalid gender. must containt pria or wanita")
		}
	}

	if ru.TrainingLevel != "" {
		if ru.TrainingLevel != "beginner" && ru.TrainingLevel != "intermediate" && ru.TrainingLevel != "advance" {
			return errors.New("invalid training level. must containt beginner or intermediate or advance")
		}
	}

	return nil
}

func (ru *ReadableUser) Validate() error {
	switch {
	case ru.Name == "":
		return errors.New("invalid name")
	case ru.Email == "":
		return errors.New("invalid email")
	case ru.Password == "":
		return errors.New("invalid password")
	case ru.Gender == "":
		return errors.New("invalid gender")
	case ru.Height == 0:
		return errors.New("invalid height")
	case ru.Weight == 0:
		return errors.New("invalid weight")
	case ru.TrainingLevel == "":
		return errors.New("invalid training level")
	}

	_, emailError := mail.ParseAddress(ru.Email)
	if emailError != nil {
		return errors.New("invalid email " + ru.Email)
	}

	if ru.Gender != "pria" && ru.Gender != "wanita" {
		return errors.New("invalid gender. must containt pria or wanita")
	}

	if ru.TrainingLevel != "beginner" && ru.TrainingLevel != "intermediate" && ru.TrainingLevel != "advance" {
		return errors.New("invalid training level. must containt beginner or intermediate or advance")
	}

	return nil
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
