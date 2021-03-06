package mapping

import (
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/storage"
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
		ID:       u.ID,
		Email:    u.Email,
		Admin:    u.Admin,
		AlumniID: u.AlumniID,
	}
}

func ToDBAlumni(r pkg.AlumniRequest, s3Filename string, provideTime time.EpochProviderFunc, genUUID uuid.GenV4Func) internal.Alumni {
	id := genUUID()
	currentTime := provideTime()
	bday := ""
	if r.Birthday != "" {
		iso, err := time.NewISO8601(r.Birthday)
		if err == nil {
			bday = iso.DateString()
		}
	}

	return internal.Alumni{
		ID:                    id,
		Title:                 r.Title,
		Firstname:             r.Firstname,
		Middlename:            r.Middlename,
		Lastname:              r.Lastname,
		MarriedName:           r.MarriedName,
		MotherName:            r.MotherName,
		FatherName:            r.FatherName,
		SpouseName:            r.SpouseName,
		CurrentAddress:        internal.Address(r.CurrentAddress),
		HomePhone:             r.HomePhone,
		CellPhone:             r.CellPhone,
		WorkPhone:             r.WorkPhone,
		EmailAddress:          r.EmailAddress,
		MiddleSchool:          internal.School(r.MiddleSchool),
		HighSchool:            internal.School(r.HighSchool),
		IsraelSchool:          internal.School(r.IsraelSchool),
		CollegeAttended:       internal.School(r.CollegeAttended),
		GradSchools:           toDBSchools(r.GradSchools),
		Profession:            r.Profession,
		Birthday:              bday,
		Clubs:                 r.Clubs,
		SportsTeams:           r.SportsTeams,
		Awards:                r.AlumniPositions,
		Committees:            r.Committees,
		OldAddresses:          toDBAddresses(r.OldAddresses),
		HillelDayCamp:         internal.Camp(r.HillelDayCamp),
		HillelSleepCamp:       internal.Camp(r.HillelSleepCamp),
		HiliDayCamp:           internal.Camp(r.HiliDayCamp),
		HiliWhiteCamp:         internal.Camp(r.HiliWhiteCamp),
		HiliInternationalCamp: internal.Camp(r.HiliInternationalCamp),
		HILI:                  r.HILI,
		HILLEL:                r.HILLEL,
		HAFTR:                 r.HAFTR,
		ParentOfStudent:       r.ParentOfStudent,
		Boards:                r.Boards,
		AlumniPositions:       r.AlumniPositions,
		Siblings:              toDBSiblings(r.Siblings),
		Children:              toDBChildren(r.Children),
		Comment:               r.Comment,
		ProfilePictureKey:     s3Filename,
		CreatedTimestamp:      currentTime,
		LastUpdatedTimestamp:  currentTime,
	}
}

func ToAlumniUpdate(r pkg.UpdateAlumniRequest, s3Filename string, provideTime time.EpochProviderFunc) internal.UpdateAlumniRequest {
	bday := ""
	if r.Birthday != "" {
		iso, err := time.NewISO8601(r.Birthday)
		if err == nil {
			bday = iso.DateString()
		}
	}

	return internal.UpdateAlumniRequest{
		Title:                 r.Title,
		Firstname:             r.Firstname,
		Middlename:            r.Middlename,
		Lastname:              r.Lastname,
		MarriedName:           r.MarriedName,
		MotherName:            r.MotherName,
		FatherName:            r.FatherName,
		SpouseName:            r.SpouseName,
		CurrentAddress:        internal.Address(r.CurrentAddress),
		HomePhone:             r.HomePhone,
		CellPhone:             r.CellPhone,
		WorkPhone:             r.WorkPhone,
		EmailAddress:          r.EmailAddress,
		MiddleSchool:          internal.School(r.MiddleSchool),
		HighSchool:            internal.School(r.HighSchool),
		IsraelSchool:          internal.School(r.IsraelSchool),
		CollegeAttended:       internal.School(r.CollegeAttended),
		GradSchools:           toDBSchools(r.GradSchools),
		Profession:            r.Profession,
		Birthday:              bday,
		Clubs:                 r.Clubs,
		SportsTeams:           r.SportsTeams,
		Awards:                r.AlumniPositions,
		Committees:            r.Committees,
		OldAddresses:          toDBAddresses(r.OldAddresses),
		HillelDayCamp:         internal.Camp(r.HillelDayCamp),
		HillelSleepCamp:       internal.Camp(r.HillelSleepCamp),
		HiliDayCamp:           internal.Camp(r.HiliDayCamp),
		HiliWhiteCamp:         internal.Camp(r.HiliWhiteCamp),
		HiliInternationalCamp: internal.Camp(r.HiliInternationalCamp),
		HILI:                  r.HILI,
		HILLEL:                r.HILLEL,
		HAFTR:                 r.HAFTR,
		ParentOfStudent:       r.ParentOfStudent,
		Boards:                r.Boards,
		AlumniPositions:       r.AlumniPositions,
		Siblings:              toDBSiblings(r.Siblings),
		Children:              toDBChildren(r.Children),
		Comment:               r.Comment,
		ProfilePictureKey:     s3Filename,
		LastUpdatedTimestamp:  provideTime(),
	}
}

func ToDTOAlumni(a internal.Alumni, presignURL storage.GetImageURLFunc) pkg.Alumni {
	url, err := presignURL(a.ProfilePictureKey)
	if err != nil {
		url = ""
	}

	return pkg.Alumni{
		ID:                    a.ID,
		Title:                 a.Title,
		Firstname:             a.Firstname,
		Middlename:            a.Middlename,
		Lastname:              a.Lastname,
		MarriedName:           a.MarriedName,
		MotherName:            a.MotherName,
		FatherName:            a.FatherName,
		SpouseName:            a.SpouseName,
		CurrentAddress:        pkg.Address(a.CurrentAddress),
		HomePhone:             a.HomePhone,
		CellPhone:             a.CellPhone,
		WorkPhone:             a.WorkPhone,
		EmailAddress:          a.EmailAddress,
		MiddleSchool:          pkg.School(a.MiddleSchool),
		HighSchool:            pkg.School(a.HighSchool),
		IsraelSchool:          pkg.School(a.IsraelSchool),
		CollegeAttended:       pkg.School(a.CollegeAttended),
		GradSchools:           toDTOSchools(a.GradSchools),
		Profession:            a.Profession,
		Birthday:              a.Birthday,
		Clubs:                 a.Clubs,
		SportsTeams:           a.SportsTeams,
		Awards:                a.AlumniPositions,
		Committees:            a.Committees,
		OldAddresses:          toDTOAddresses(a.OldAddresses),
		HillelDayCamp:         pkg.Camp(a.HillelDayCamp),
		HillelSleepCamp:       pkg.Camp(a.HillelSleepCamp),
		HiliDayCamp:           pkg.Camp(a.HiliDayCamp),
		HiliWhiteCamp:         pkg.Camp(a.HiliWhiteCamp),
		HiliInternationalCamp: pkg.Camp(a.HiliInternationalCamp),
		HILI:                  a.HILI,
		HILLEL:                a.HILLEL,
		HAFTR:                 a.HAFTR,
		ParentOfStudent:       a.ParentOfStudent,
		Boards:                a.Boards,
		AlumniPositions:       a.AlumniPositions,
		Siblings:              toDTOSiblings(a.Siblings),
		Children:              toDTOChildren(a.Children),
		Comment:               a.Comment,
		ProfilePictureURL:     url,
	}
}

func toDBSchools(ss []pkg.School) []internal.School {
	newSS := []internal.School{}
	for _, s := range ss {
		newSS = append(newSS, internal.School(s))
	}
	return newSS
}

func toDTOSchools(ss []internal.School) []pkg.School {
	newSS := []pkg.School{}
	for _, s := range ss {
		newSS = append(newSS, pkg.School(s))
	}
	return newSS
}

func toDBAddresses(aa []pkg.Address) []internal.Address {
	newAA := []internal.Address{}
	for _, a := range aa {
		newAA = append(newAA, internal.Address(a))
	}
	return newAA
}

func toDTOAddresses(aa []internal.Address) []pkg.Address {
	newAA := []pkg.Address{}
	for _, a := range aa {
		newAA = append(newAA, pkg.Address(a))
	}
	return newAA
}

func toDBSiblings(ss []pkg.Sibling) []internal.Sibling {
	newSS := []internal.Sibling{}
	for _, s := range ss {
		newSS = append(newSS, internal.Sibling{
			Firstname:     s.Firstname,
			Lastname:      s.Lastname,
			YearCompleted: s.YearCompleted,
			MiddleSchool:  internal.School(s.MiddleSchool),
			HighSchool:    internal.School(s.HighSchool),
		})
	}
	return newSS
}

func toDTOSiblings(ss []internal.Sibling) []pkg.Sibling {
	newSS := []pkg.Sibling{}
	for _, s := range ss {
		newSS = append(newSS, pkg.Sibling{
			Firstname:     s.Firstname,
			Lastname:      s.Lastname,
			YearCompleted: s.YearCompleted,
			MiddleSchool:  pkg.School(s.MiddleSchool),
			HighSchool:    pkg.School(s.HighSchool),
		})
	}
	return newSS
}

func toDBChildren(cc []pkg.Child) []internal.Child {
	newCC := []internal.Child{}
	for _, c := range cc {
		newCC = append(newCC, internal.Child(c))
	}
	return newCC
}

func toDTOChildren(cc []internal.Child) []pkg.Child {
	newCC := []pkg.Child{}
	for _, c := range cc {
		newCC = append(newCC, pkg.Child(c))
	}
	return newCC
}
