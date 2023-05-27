package models

type CustomError struct {
	StatusCode   int
	ErrorReason  string
	ErrorMessage error
}

func (ce *CustomError) IsError() bool {
	return ce.ErrorMessage != nil
}

func (ce *CustomError) ErrParseIdParam(err error) {
	ce.StatusCode = 400
	ce.ErrorReason = "invalid id parameter"
	ce.ErrorMessage = err
}

func (ce *CustomError) NoRecordFound(err error) {
	ce.StatusCode = 200
	ce.ErrorReason = "no record found"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailGenerateGenderType(err error) {
	ce.StatusCode = 400
	ce.ErrorReason = "invalid gender: fail to generate gender type"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailLogin(err error) {
	ce.StatusCode = 401
	ce.ErrorReason = "login failed"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailLoginWrongPassword(err error) {
	ce.StatusCode = 401
	ce.ErrorReason = "login failed: wrong password"
	ce.ErrorMessage = err
}

func (ce *CustomError) DuplicateKey(err error) {
	ce.StatusCode = 409
	ce.ErrorReason = "duplicate record found"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailRetrieveDataFromDB(err error) {
	ce.StatusCode = 500
	ce.ErrorReason = "fail to retrieve data from database"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailCreateDataInDB(err error) {
	ce.StatusCode = 500
	ce.ErrorReason = "fail to create new data in database"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailEditDataInDB(err error) {
	ce.StatusCode = 500
	ce.ErrorReason = "fail to edit data in database"
	ce.ErrorMessage = err
}

func (ce *CustomError) FailDeleteDataInDB(err error) {
	ce.StatusCode = 500
	ce.ErrorReason = "fail to delete data in database"
	ce.ErrorMessage = err
}
