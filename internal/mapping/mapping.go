package mapping

import (
	"encoding/base64"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

// ToDbUser maps a UserRequest to an internal User
func ToDbUser(req pkg.UserRequest, securePw []byte, genUUID uuid.GenV4Func, provideTime time.EpochProviderFunc) internal.User {
	return internal.User{
		ID:                   genUUID(),
		Email:                strings.ToLower(req.Email),
		Password:             securePw,
		Admin:                false,
		CreatedTimestamp:     provideTime(),
		LastUpdatedTimestamp: provideTime(),
	}
}

// ToDTOUser maps an internal User to a pkg User
func ToDTOUser(u internal.User) pkg.User {
	return pkg.User{
		Token:                base64.StdEncoding.EncodeToString([]byte(u.ID)),
		Email:                u.Email,
		Admin:                u.Admin,
		AlumniID:             u.AlumniID,
		CreatedTimestamp:     u.CreatedTimestamp.ToISO8601(),
		LastUpdatedTimestamp: u.LastUpdatedTimestamp.ToISO8601(),
	}
}

// ToDTOAlumni maps an internal Alumni to a pkg Alumni
func ToDTOAlumni(a internal.Alumni) pkg.Alumni {
	return pkg.Alumni{}
}
