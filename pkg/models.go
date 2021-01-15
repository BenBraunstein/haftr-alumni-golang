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

// Alumni is the DTO representation of an Alumni
type Alumni struct {
	UserID                uuid.V4    `json:"id"`
	Firstname             string     `json:"firstname"`
	Middlename            string     `json:"middlename"`
	Lastname              string     `json:"lastname"`
	MarriedName           string     `json:"marriedName"`
	MotherName            string     `json:"motherName"`
	FatherName            string     `json:"fatherName"`
	CurrentAddress        string     `json:"address"`
	HomePhone             string     `json:"homePhone"`
	CellPhone             string     `json:"cellPhone"`
	WorkPhone             string     `json:"workPhone"`
	EmailAddress          string     `json:"emailAddress"`
	LastYearAttended      string     `json:"lastYearAttended"`
	IsraelSchool          string     `json:"israelSchool"`
	CollegeAttended       string     `json:"collegeAttended"`
	GradSchool            string     `json:"gradSchool"`
	Profession            string     `json:"profession"`
	Birthday              string     `json:"birthday"`
	Clubs                 []string   `json:"clubs"`
	SportsTeams           []string   `json:"sportsTeams"`
	Awards                []string   `json:"awards"`
	Committees            []string   `json:"committees"`
	OldAddresses          []string   `json:"oldAddresses"`
	HillelDayCamp         Camp       `json:"hillelDayCamp"`
	HillelSleepCamp       Camp       `json:"hillelSleepCamp"`
	HiliDayCamp           Camp       `json:"hiliDayCamp"`
	HiliWhiteCamp         Camp       `json:"hiliWhiteCamp"`
	HiliInternationalCamp Camp       `json:"hiliInternationalCamp"`
	HILI                  bool       `json:"hili"`
	HILLEL                bool       `json:"hillel"`
	HAFTR                 bool       `json:"haftr"`
	ParentOfStudent       bool       `json:"parentOfStudent"`
	Boards                []string   `json:"boards"`
	AlumniPositions       []string   `json:"alumniPositions"`
	Siblings              []Sibling  `json:"siblings"`
	Children              []Child    `json:"children"`
	ProfilePictureKey     string     `json:"profilePictureKey"`
	CreatedTimestamp      time.Epoch `json:"createdTimestamp"`
	LastUpdatedTimestamp  time.Epoch `json:"lastUpdatedTimestamp"`
}

// Camp is the DTO representation of a camp
type Camp struct {
	Attended  bool   `json:"attended"`
	StartYear string `json:"startYear"`
	EndYear   string `json:"endYear"`
	Specialty string `json:"specialty"`
	Camper    bool   `json:"camper"`
	Counselor bool   `json:"counselor"`
}

// Sibling is the DTO representation of a Sibling
type Sibling struct {
	Name          string `json:"name"`
	YearCompleted string `json:"yearCompleted"`
	School        string `json:"school"`
}

// Child is the DTO representation of a Child
type Child struct {
	Name           string `json:"name"`
	GraduationYear string `json:"graduationYear"`
}
