package mapping

import (
	"fmt"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/storage"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

// ToDbUser maps a UserRequest to an internal User
func ToDbUser(req pkg.UserRequest, securePw []byte, genUUID uuid.GenV4Func, provideTime time.EpochProviderFunc) internal.User {
	currentTime := provideTime()

	return internal.User{
		ID:                   genUUID(),
		Email:                strings.ToLower(req.Email),
		Password:             securePw,
		Admin:                false,
		Status:               internal.PendingUserStatus,
		CreatedTimestamp:     currentTime,
		LastUpdatedTimestamp: currentTime,
	}
}

// ToDTOUser maps an internal User to a pkg User
func ToDTOUser(u internal.User) pkg.User {
	s := u.Status
	if s == "" {
		s = internal.PendingUserStatus
	}
	return pkg.User{
		ID:       u.ID,
		Email:    u.Email,
		Admin:    u.Admin,
		AlumniID: u.AlumniID,
		Status:   s,
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
		MaidenName:            r.MaidenName,
		MotherName:            r.MotherName,
		FatherName:            r.FatherName,
		SpouseName:            r.SpouseName,
		SpouseMaidenName:      r.SpouseMaidenName,
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
		Grandparents:          toDBGrandparents(r.Grandparents),
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
		MaidenName:            r.MaidenName,
		MotherName:            r.MotherName,
		FatherName:            r.FatherName,
		SpouseName:            r.SpouseName,
		SpouseMaidenName:      r.SpouseMaidenName,
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
		Grandparents:          toDBGrandparents(r.Grandparents),
		Comment:               r.Comment,
		ProfilePictureKey:     s3Filename,
		LastUpdatedTimestamp:  provideTime(),
	}
}

func ToDTOAlumni(a internal.Alumni, presignURL storage.GetImageURLFunc, u internal.User) pkg.Alumni {
	url, err := presignURL(a.ProfilePictureKey)
	if err != nil {
		url = ""
	}

	return pkg.Alumni{
		UserID:                u.ID,
		Status:                u.Status,
		ID:                    a.ID,
		Title:                 a.Title,
		Firstname:             a.Firstname,
		Middlename:            a.Middlename,
		Lastname:              a.Lastname,
		HighSchoolGradYear:    a.HighSchool.YearEnded,
		MarriedName:           a.MarriedName,
		MaidenName:            a.MaidenName,
		MotherName:            a.MotherName,
		FatherName:            a.FatherName,
		SpouseName:            a.SpouseName,
		SpouseMaidenName:      a.SpouseMaidenName,
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
		Grandparents:          toDTOGrandparents(a.Grandparents),
		Comment:               a.Comment,
		IsPublic:              a.IsPublic,
		ProfilePictureURL:     url,
	}
}

func ToCleanAlumni(a internal.Alumni, presignURL storage.GetImageURLFunc, u internal.User) pkg.CleanAlumni {
	url, err := presignURL(a.ProfilePictureKey)
	if err != nil {
		url = ""
	}

	return pkg.CleanAlumni{
		ID:                 a.ID,
		Status:             u.Status,
		Firstname:          a.Firstname,
		Lastname:           a.Lastname,
		HighSchoolGradYear: a.HighSchool.YearEnded,
		EmailAddress:       a.EmailAddress,
		ProfilePictureURL:  url,
	}
}

func ToHappyBirthdayAlumni(a internal.Alumni) pkg.HappyBirthdayAlumni {
	aDTO := pkg.HappyBirthdayAlumni{
		Firstname:          a.Firstname,
		Lastname:           a.Lastname,
		HighSchoolGradYear: a.HighSchool.YearEnded,
	}
	if a.Birthday != "" {
		iso, _ := time.NewISO8601(a.Birthday)
		ds := iso.DateString()
		m := strings.Split(ds, "-")[1]
		d := strings.Split(ds, "-")[2]
		bday := fmt.Sprintf("%v-%v", m, d)
		aDTO.Birthday = bday
	}

	return aDTO
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

func toDBGrandparents(gg []pkg.Grandparent) []internal.Grandparent {
	newGG := []internal.Grandparent{}
	for _, g := range gg {
		newGG = append(newGG, internal.Grandparent(g))
	}
	return newGG
}

func toDTOGrandparents(gg []internal.Grandparent) []pkg.Grandparent {
	newGG := []pkg.Grandparent{}
	for _, g := range gg {
		newGG = append(newGG, pkg.Grandparent(g))
	}
	return newGG
}
