package workflow

import (
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

// AddUserFunc signs up a new user and returns a representation of the user and a token for the user
type AddUserFunc func(req pkg.UserRequest) (pkg.User, string, error)

// LoginUserFunc logs in a user and returns a representation of the user and a token for the user
type LoginUserFunc func(req pkg.UserRequest) (pkg.User, string, error)

// AutoLoginUserFunc logs in a user using their JWT tokenand returns a representation of the user and a refreshed token for the user
type AutoLoginUserFunc func(tokenString string) (pkg.User, string, error)

// ApproveUserFunc returns functionaliy for an admin to approve a user
type ApproveUserFunc func(userId string, tokenString string) (pkg.User, error)

// DenyUserFunc returns functionaliy for an admin to deny a user
type DenyUserFunc func(userId string, tokenString string) (pkg.User, error)

// AddAlumniFunc returns functionality to add an alumni
type AddAlumniFunc func(req pkg.AlumniRequest, fileData pkg.FileData, tokenString string, skipFileUpload bool) (pkg.Alumni, error)

// UpdateAlumniFunc returns functionality to update an alumni
type UpdateAlumniFunc func(req pkg.UpdateAlumniRequest, alumniId string, fileData pkg.FileData, tokenString string, skipFileUpload bool) (pkg.Alumni, error)

// RetrieveAlumniByIDFunc returns functionality to retrieve an alumni by ID
type RetrieveAlumniByIDFunc func(alumniId string, tokenString string) (pkg.AlumniInterface, error)

// ChangeAlumniPrivacyFunc returns functionality to change the privacy status of an alumni by their ID
type ChangeAlumniPrivacyFunc func(alumniId string, tokenString string) (pkg.Alumni, error)

// RetrieveAlumniFunc returns functionality to retrieve all alumni
type RetrieveAlumniFunc func(params pkg.QueryParams, tokenString string) ([]pkg.CleanAlumni, pkg.PageInfo, error)

// HappyBirthdayFunc returns functionality to retrieve alumni's with todays birthday
type HappyBirthdayFunc func() (pkg.HappyBirthdayResponse, error)

// HappyBirthdayEmailFunc returns functionality to send an email to all alumni's with todays birthday
type HappyBirthdayEmailFunc func() error

// ExportCSVFunc returns functionality to export a CSV with provided query params
type ExportCSVFunc func(params pkg.QueryParams, tokenString string) ([]byte, error)

// ForgotPasswordFunc returns functionality to send a reset password email
type ForgotPasswordFunc func(email string) error

// SetNewPasswordFunc returns functionality to set a new password
type SetNewPasswordFunc func(rp pkg.ResetPassword) (pkg.User, string, error)
