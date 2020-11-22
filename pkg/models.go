package pkg

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID                   uuid.V4      `json:"id"`
	Email                string       `json:"email"`
	AlumniID             uuid.V4      `json:"alumniId,omitempty"`
	Admin                bool         `json:"admin"`
	CreatedTimestamp     time.ISO8601 `json:"createdTimestamp"`
	LastUpdatedTimestamp time.ISO8601 `json:"lastUpdatedTimestamp"`
}

type Alumni struct {
}
