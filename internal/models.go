package internal

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

type User struct {
	ID                   uuid.V4    `bson:"id"`
	Email                string     `bson:"email"`
	Password             []byte     `bson:"password"`
	AlumniID             uuid.V4    `bson:"alumniId"`
	Admin                bool       `bson:"admin"`
	CreatedTimestamp     time.Epoch `bson:"createdTimestamp"`
	LastUpdatedTimestamp time.Epoch `bson:"lastUpdatedTimestamp"`
}

type Alumni struct {
}
