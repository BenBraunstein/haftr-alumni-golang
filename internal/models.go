package internal

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

// User is the internal representation of a user
type User struct {
	ID                   uuid.V4    `bson:"id"`
	Email                string     `bson:"email"`
	Password             []byte     `bson:"password"`
	AlumniID             uuid.V4    `bson:"alumniId"`
	Admin                bool       `bson:"admin"`
	CreatedTimestamp     time.Epoch `bson:"createdTimestamp"`
	LastUpdatedTimestamp time.Epoch `bson:"lastUpdatedTimestamp"`
}

// Alumni is the internal representation of an Alumni
type Alumni struct {
	UserID                uuid.V4    `bson:"id"`
	Firstname             string     `bson:"firstname"`
	Middlename            string     `bson:"middlename"`
	Lastname              string     `bson:"lastname"`
	MarriedName           string     `bson:"marriedName"`
	MotherName            string     `bson:"motherName"`
	FatherName            string     `bson:"fatherName"`
	CurrentAddress        string     `bson:"address"`
	HomePhone             string     `bson:"homePhone"`
	CellPhone             string     `bson:"cellPhone"`
	WorkPhone             string     `bson:"workPhone"`
	EmailAddress          string     `bson:"emailAddress"`
	LastYearAttended      string     `bson:"lastYearAttended"`
	IsraelSchool          string     `bson:"israelSchool"`
	CollegeAttended       string     `bson:"collegeAttended"`
	GradSchool            string     `bson:"gradSchool"`
	Profession            string     `bson:"profession"`
	Birthday              string     `bson:"birthday"`
	Clubs                 []string   `bson:"clubs"`
	SportsTeams           []string   `bson:"sportsTeams"`
	Awards                []string   `bson:"awards"`
	Committees            []string   `bson:"committees"`
	OldAddresses          []string   `bson:"oldAddresses"`
	HillelDayCamp         Camp       `bson:"hillelDayCamp"`
	HillelSleepCamp       Camp       `bson:"hillelSleepCamp"`
	HiliDayCamp           Camp       `bson:"hiliDayCamp"`
	HiliWhiteCamp         Camp       `bson:"hiliWhiteCamp"`
	HiliInternationalCamp Camp       `bson:"hiliInternationalCamp"`
	HILI                  bool       `bson:"hili"`
	HILLEL                bool       `bson:"hillel"`
	HAFTR                 bool       `bson:"haftr"`
	ParentOfStudent       bool       `bson:"parentOfStudent"`
	Boards                []string   `bson:"boards"`
	AlumniPositions       []string   `bson:"alumniPositions"`
	Siblings              []Sibling  `bson:"siblings"`
	Children              []Child    `bson:"children"`
	ProfilePictureKey     string     `bson:"profilePictureKey"`
	CreatedTimestamp      time.Epoch `bson:"createdTimestamp"`
	LastUpdatedTimestamp  time.Epoch `bson:"lastUpdatedTimestamp"`
}

// Camp is the internal representation of a camp
type Camp struct {
	Attended  bool   `bson:"attended"`
	StartYear string `bson:"startYear"`
	EndYear   string `bson:"endYear"`
	Specialty string `bson:"specialty"`
	Camper    bool   `bson:"camper"`
	Counselor bool   `bson:"counselor"`
}

// Sibling is the internal representation of a Sibling
type Sibling struct {
	Name          string `bson:"name"`
	YearCompleted string `bson:"yearCompleted"`
	School        string `bson:"school"`
}

// Child is the internal representation of a Child
type Child struct {
	Name           string `bson:"name"`
	GraduationYear string `bson:"graduationYear"`
}
