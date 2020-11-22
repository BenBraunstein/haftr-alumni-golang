package mapping

import (
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

func ToDbUser(req pkg.UserRequest, securePw []byte, genUUID uuid.GenV4Func, provideTime time.EpochProviderFunc) internal.User {
	return internal.User{
		ID:                   genUUID(),
		Email:                strings.ToLower(req.Email),
		Password:             securePw,
		Admin:                false,
		Alumni:               internal.Alumni{},
		CreatedTimestamp:     provideTime(),
		LastUpdatedTimestamp: provideTime(),
	}
}

func ToDTOUser(u internal.User) pkg.User {
	return pkg.User{
		ID:                   u.ID,
		Email:                u.Email,
		Admin:                u.Admin,
		Alumni:               ToDTOAlumni(u.Alumni),
		CreatedTimestamp:     u.CreatedTimestamp.ToISO8601(),
		LastUpdatedTimestamp: u.LastUpdatedTimestamp.ToISO8601(),
	}
}

func ToDTOAlumni(a internal.Alumni) pkg.Alumni {
	return pkg.Alumni{}
}
