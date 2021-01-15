package pkg

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

// UserRequest is a representation of a request to make a new user
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AlumniRequest is a representation of a request to make a new alumni
type AlumniRequest struct {
}

// User is the representation of a DTO User
type User struct {
	Token                string       `json:"token"`
	Email                string       `json:"email"`
	AlumniID             uuid.V4      `json:"alumniId"`
	Admin                bool         `json:"admin"`
	CreatedTimestamp     time.ISO8601 `json:"createdTimestamp"`
	LastUpdatedTimestamp time.ISO8601 `json:"lastUpdatedTimestamp"`
}

// Alumni is the representation of a DTO alumni
type Alumni struct {
}
