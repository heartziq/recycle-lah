package handlers

import "errors"

//  type ValidationError error defines error type for input validation
type ValidationError error

var (
	errUserNameFmt = ValidationError(errors.New("Must start/end with letters or digits, may contain . and underscore.  e.g.:user1234, user.1234"))
	// errUserNameFmt      = ValidationError(errors.New("Only letters(a-z), numbers(0-9) and periods(.) are allowed"))
	errNoUserName       = ValidationError(errors.New("Username must not be empty"))
	errUserNameNotAvail = ValidationError(errors.New("Username is not available"))
	errUserNameLength   = ValidationError(errors.New("username length must be between 6 and 30"))
	errNoPassword       = ValidationError(errors.New("missing password"))
	errConfirmPassword  = ValidationError(errors.New("Confirmation password not matched"))
	errPasswordLength   = ValidationError(errors.New("password length must be between 8 and 64"))
	errPasswordFormat   = ValidationError(errors.New("use 8 or more characters with a mix of letters, numbers & symbols"))
	errEmailFormat      = ValidationError(errors.New("invalid email format"))
	// errEmailFormat2     = ValidationError(errors.New("username must be between 6 and 30 characters"))
	errExceedAttempt = ValidationError(errors.New("Exceed three attempts"))
	errNoId          = ValidationError(errors.New("missing id"))
)

// type AppUserError error defines error type for the end users
type AppUserError error

var (
	userErrInvalidApiKey    = AppUserError(errors.New("Invalid application key")) // apikey not found
	userErrApiKeyNotOk      = AppUserError(errors.New("Invalid application key"))
	userErrMissingApiKey    = AppUserError(errors.New("Missing application key"))
	userErrMissingAccount   = AppUserError(errors.New("Missing account name"))
	userErrNoRecord         = AppUserError(errors.New("Record not found"))
	userErrMissingCode      = AppUserError(errors.New("Missing course code"))
	userErrDuplicateCode    = AppUserError(errors.New("Duplicate course code"))
	userErrGeneral          = AppUserError(errors.New("Information currently not available"))
	userErrReadReqBody      = AppUserError(errors.New("Unable to process request"))
	userErrUnmarshalReqBody = AppUserError(errors.New("Unable to process request"))

	userErrQueryValues      = AppUserError(errors.New("Missing query values"))
	userErrNotDeleted       = AppUserError(errors.New("Unable to delete, record not found"))
	userErrNotUpdated       = AppUserError(errors.New("Unable to update"))
	userErrNotUpdatedNotAdd = AppUserError(errors.New("Unable to update/add"))
	userErrNotAdded         = AppUserError(errors.New("Unable to add recod"))
)

// type AppError error defines error type mainly meant for internal use
type AppError error

var (
	errSQLStmt     = AppError(errors.New("Error in sql statement"))
	errRareApiKey  = AppError(errors.New("Application Key not valid"))
	errMoreThanOne = AppError(errors.New("more db records processed for unique code"))
	errUndetected  = AppError(errors.New("undetected condition"))
)

type AuthenticationError error

var errErrAuthenticate = AuthenticationError(errors.New("Invalid username and/or password"))

func appUserError(err error) string {
	if _, ok := err.(AppUserError); ok {
		return err.Error()
	}
	return userErrGeneral.Error()
}
