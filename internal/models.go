package internal

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

const (
	DefaultPageLimit          = 20
	EmailRecipient            = "benzbraunstein@gmail.com"
	NewAlumniTemplateName     = "NEW_ALUMNI"
	UpdatedAlumniTemplateName = "UPDATED_ALUMNI"
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
	ID                    uuid.V4    `bson:"id"`
	Title                 string     `bson:"title"`
	Firstname             string     `bson:"firstname"`
	Middlename            string     `bson:"middlename"`
	Lastname              string     `bson:"lastname"`
	MarriedName           string     `bson:"marriedName"`
	MotherName            string     `bson:"motherName"`
	FatherName            string     `bson:"fatherName"`
	SpouseName            string     `bson:"spouseName"`
	CurrentAddress        Address    `bson:"address"`
	HomePhone             string     `bson:"homePhone"`
	CellPhone             string     `bson:"cellPhone"`
	WorkPhone             string     `bson:"workPhone"`
	EmailAddress          string     `bson:"emailAddress"`
	MiddleSchool          School     `bson:"middleschool"`
	HighSchool            School     `bson:"highschool"`
	IsraelSchool          School     `bson:"israelSchool"`
	CollegeAttended       School     `bson:"collegeAttended"`
	GradSchools           []School   `bson:"gradSchools"`
	Profession            []string   `bson:"profession"`
	Birthday              string     `bson:"birthday"`
	Clubs                 []string   `bson:"clubs"`
	SportsTeams           []string   `bson:"sportsTeams"`
	Awards                []string   `bson:"awards"`
	Committees            []string   `bson:"committees"`
	OldAddresses          []Address  `bson:"oldAddresses"`
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
	Comment               string     `bson:"comment"`
	IsPublic              bool       `bson:"isPublic"`
	ProfilePictureKey     string     `bson:"profilePictureKey"`
	CreatedTimestamp      time.Epoch `bson:"createdTimestamp"`
	LastUpdatedTimestamp  time.Epoch `bson:"lastUpdatedTimestamp"`
}

type UpdateAlumniRequest struct {
	Title                 string     `bson:"title,omitempty"`
	Firstname             string     `bson:"firstname,omitempty"`
	Middlename            string     `bson:"middlename,omitempty"`
	Lastname              string     `bson:"lastname,omitempty"`
	MarriedName           string     `bson:"marriedName,omitempty"`
	MotherName            string     `bson:"motherName,omitempty"`
	FatherName            string     `bson:"fatherName,omitempty"`
	SpouseName            string     `bson:"spouseName,omitempty"`
	CurrentAddress        Address    `bson:"address,omitempty"`
	HomePhone             string     `bson:"homePhone,omitempty"`
	CellPhone             string     `bson:"cellPhone,omitempty"`
	WorkPhone             string     `bson:"workPhone,omitempty"`
	EmailAddress          string     `bson:"emailAddress,omitempty"`
	MiddleSchool          School     `bson:"middleschool,omitempty"`
	HighSchool            School     `bson:"highschool,omitempty"`
	IsraelSchool          School     `bson:"israelSchool,omitempty"`
	CollegeAttended       School     `bson:"collegeAttended,omitempty"`
	GradSchools           []School   `bson:"gradSchools,omitempty"`
	Profession            []string   `bson:"profession,omitempty"`
	Birthday              string     `bson:"birthday,omitempty"`
	Clubs                 []string   `bson:"clubs,omitempty"`
	SportsTeams           []string   `bson:"sportsTeams,omitempty"`
	Awards                []string   `bson:"awards,omitempty"`
	Committees            []string   `bson:"committees,omitempty"`
	OldAddresses          []Address  `bson:"oldAddresses,omitempty"`
	HillelDayCamp         Camp       `bson:"hillelDayCamp,omitempty"`
	HillelSleepCamp       Camp       `bson:"hillelSleepCamp,omitempty"`
	HiliDayCamp           Camp       `bson:"hiliDayCamp,omitempty"`
	HiliWhiteCamp         Camp       `bson:"hiliWhiteCamp,omitempty"`
	HiliInternationalCamp Camp       `bson:"hiliInternationalCamp,omitempty"`
	HILI                  bool       `bson:"hili,omitempty"`
	HILLEL                bool       `bson:"hillel,omitempty"`
	HAFTR                 bool       `bson:"haftr,omitempty"`
	ParentOfStudent       bool       `bson:"parentOfStudent,omitempty"`
	Boards                []string   `bson:"boards,omitempty"`
	AlumniPositions       []string   `bson:"alumniPositions,omitempty"`
	Siblings              []Sibling  `bson:"siblings,omitempty"`
	Children              []Child    `bson:"children,omitempty"`
	Comment               string     `bson:"comment,omitempty"`
	IsPublic              bool       `bson:"isPublic,omitempty"`
	ProfilePictureKey     string     `bson:"profilePictureKey"`
	LastUpdatedTimestamp  time.Epoch `bson:"lastUpdatedTimestamp"`
}

type School struct {
	Name        string `bson:"name"`
	YearStarted string `bson:"yearStarted"`
	YearEnded   string `bson:"yearEnded"`
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
	Firstname     string `bson:"firstname"`
	Lastname      string `bson:"lastname"`
	YearCompleted string `bson:"yearCompleted"`
	MiddleSchool  School `bson:"middleSchool"`
	HighSchool    School `bson:"highSchool"`
}

// Child is the internal representation of a Child
type Child struct {
	Firstname      string `bson:"firstname"`
	Lastname       string `bson:"lastname"`
	GraduationYear string `bson:"graduationYear"`
}

type Address struct {
	Line1   string `bson:"line1"`
	Line2   string `bson:"line2"`
	City    string `bson:"city"`
	State   string `bson:"state"`
	Zip     string `bson:"zip"`
	Country string `bson:"country"`
}

type EmailTemplate struct {
	Name    string `bson:"name"`
	Subject string `bson:"subject"`
	HTML    string `bson:"html"`
}
