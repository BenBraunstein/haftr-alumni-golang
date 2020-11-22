package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

// V4 is a type representing a v4 UUID string
type V4 string

// Val returns the underlying value for a UUID V4
func (v V4) Val() string {
	return string(v)
}

// GenV4Func generates a v4 UUID string
type GenV4Func func() V4

//FromStringFunc creates a valid uuid V4 from a string
type FromStringFunc func(uuid string) (V4, error)

// GenV4 returns v4 UUID
func GenV4() V4 {
	return V4(uuid.New().String())
}

// IsValidUUID checks if a given string is a valid v4 UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// GenFromString returns a uuid  from string
func GenFromString(uuid string) (V4, error) {
	ok := IsValidUUID(uuid)
	if !ok {
		return "", fmt.Errorf("Unable to generate UUID V4 from invalid value=%s", uuid)
	}
	return V4(uuid), nil
}
