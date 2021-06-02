package pkg

import (
	"io"
	"mime/multipart"

	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

// UserRequest is a representation of a request to make a new user
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User is the representation of a DTO User
type User struct {
	ID       uuid.V4 `json:"id"`
	Email    string  `json:"email"`
	AlumniID uuid.V4 `json:"alumniId"`
	Admin    bool    `json:"admin"`
}

type UserResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// AlumniRequest is a representation of a request to make a new alumni
type AlumniRequest struct {
	Title                 string    `json:"title"`
	Firstname             string    `json:"firstname"`
	Middlename            string    `json:"middlename"`
	Lastname              string    `json:"lastname"`
	MarriedName           string    `json:"marriedName"`
	MotherName            string    `json:"motherName"`
	FatherName            string    `json:"fatherName"`
	SpouseName            string    `json:"spouseName"`
	CurrentAddress        Address   `json:"address"`
	HomePhone             string    `json:"homePhone"`
	CellPhone             string    `json:"cellPhone"`
	WorkPhone             string    `json:"workPhone"`
	EmailAddress          string    `json:"emailAddress"`
	MiddleSchool          School    `json:"middleschool"`
	HighSchool            School    `json:"highschool"`
	IsraelSchool          School    `json:"israelSchool"`
	CollegeAttended       School    `json:"collegeAttended"`
	GradSchools           []School  `json:"gradSchools"`
	Profession            []string  `json:"profession"`
	Birthday              string    `json:"birthday"`
	Clubs                 []string  `json:"clubs"`
	SportsTeams           []string  `json:"sportsTeams"`
	Awards                []string  `json:"awards"`
	Committees            []string  `json:"committees"`
	OldAddresses          []Address `json:"oldAddresses"`
	HillelDayCamp         Camp      `json:"hillelDayCamp"`
	HillelSleepCamp       Camp      `json:"hillelSleepCamp"`
	HiliDayCamp           Camp      `json:"hiliDayCamp"`
	HiliWhiteCamp         Camp      `json:"hiliWhiteCamp"`
	HiliInternationalCamp Camp      `json:"hiliInternationalCamp"`
	HILI                  bool      `json:"hili"`
	HILLEL                bool      `json:"hillel"`
	HAFTR                 bool      `json:"haftr"`
	ParentOfStudent       bool      `json:"parentOfStudent"`
	Boards                []string  `json:"boards"`
	AlumniPositions       []string  `json:"alumniPositions"`
	Siblings              []Sibling `json:"siblings"`
	Children              []Child   `json:"children"`
	Comment               string    `json:"comment"`
}

type UpdateAlumniRequest struct {
	Title                 *string   `json:"title,omitempty"`
	Firstname             *string   `json:"firstname,omitempty"`
	Middlename            *string   `json:"middlename,omitempty"`
	Lastname              *string   `json:"lastname,omitempty"`
	MarriedName           *string   `json:"marriedName,omitempty"`
	MotherName            *string   `json:"motherName,omitempty"`
	FatherName            *string   `json:"fatherName,omitempty"`
	SpouseName            *string   `json:"spouseName,omitempty"`
	CurrentAddress        Address   `json:"address,omitempty"`
	HomePhone             *string   `json:"homePhone,omitempty"`
	CellPhone             *string   `json:"cellPhone,omitempty"`
	WorkPhone             *string   `json:"workPhone,omitempty"`
	EmailAddress          *string   `json:"emailAddress,omitempty"`
	MiddleSchool          School    `json:"middleschool,omitempty"`
	HighSchool            School    `json:"highschool,omitempty"`
	IsraelSchool          School    `json:"israelSchool,omitempty"`
	CollegeAttended       School    `json:"collegeAttended,omitempty"`
	GradSchools           []School  `json:"gradSchools,omitempty"`
	Profession            *[]string `json:"profession,omitempty"`
	Birthday              *string   `json:"birthday,omitempty"`
	Clubs                 *[]string `json:"clubs,omitempty"`
	SportsTeams           *[]string `json:"sportsTeams,omitempty"`
	Awards                *[]string `json:"awards,omitempty"`
	Committees            *[]string `json:"committees,omitempty"`
	OldAddresses          []Address `json:"oldAddresses,omitempty"`
	HillelDayCamp         Camp      `json:"hillelDayCamp,omitempty"`
	HillelSleepCamp       Camp      `json:"hillelSleepCamp,omitempty"`
	HiliDayCamp           Camp      `json:"hiliDayCamp,omitempty"`
	HiliWhiteCamp         Camp      `json:"hiliWhiteCamp,omitempty"`
	HiliInternationalCamp Camp      `json:"hiliInternationalCamp,omitempty"`
	HILI                  *bool     `json:"hili,omitempty"`
	HILLEL                *bool     `json:"hillel,omitempty"`
	HAFTR                 *bool     `json:"haftr,omitempty"`
	ParentOfStudent       *bool     `json:"parentOfStudent,omitempty"`
	Boards                *[]string `json:"boards,omitempty"`
	AlumniPositions       *[]string `json:"alumniPositions,omitempty"`
	Siblings              []Sibling `json:"siblings,omitempty"`
	Children              []Child   `json:"children,omitempty"`
	Comment               *string   `json:"comment,omitempty"`
}

type Alumni struct {
	ID                    uuid.V4   `json:"id"`
	Title                 string    `json:"title"`
	Firstname             string    `json:"firstname"`
	Middlename            string    `json:"middlename"`
	Lastname              string    `json:"lastname"`
	MarriedName           string    `json:"marriedName"`
	MotherName            string    `json:"motherName"`
	FatherName            string    `json:"fatherName"`
	SpouseName            string    `json:"spouseName"`
	CurrentAddress        Address   `json:"address"`
	HomePhone             string    `json:"homePhone"`
	CellPhone             string    `json:"cellPhone"`
	WorkPhone             string    `json:"workPhone"`
	EmailAddress          string    `json:"emailAddress"`
	MiddleSchool          School    `json:"middleschool"`
	HighSchool            School    `json:"highschool"`
	IsraelSchool          School    `json:"israelSchool"`
	CollegeAttended       School    `json:"collegeAttended"`
	GradSchools           []School  `json:"gradSchools"`
	Profession            []string  `json:"profession"`
	Birthday              string    `json:"birthday"`
	Clubs                 []string  `json:"clubs"`
	SportsTeams           []string  `json:"sportsTeams"`
	Awards                []string  `json:"awards"`
	Committees            []string  `json:"committees"`
	OldAddresses          []Address `json:"oldAddresses"`
	HillelDayCamp         Camp      `json:"hillelDayCamp"`
	HillelSleepCamp       Camp      `json:"hillelSleepCamp"`
	HiliDayCamp           Camp      `json:"hiliDayCamp"`
	HiliWhiteCamp         Camp      `json:"hiliWhiteCamp"`
	HiliInternationalCamp Camp      `json:"hiliInternationalCamp"`
	HILI                  bool      `json:"hili"`
	HILLEL                bool      `json:"hillel"`
	HAFTR                 bool      `json:"haftr"`
	ParentOfStudent       bool      `json:"parentOfStudent"`
	Boards                []string  `json:"boards"`
	AlumniPositions       []string  `json:"alumniPositions"`
	Siblings              []Sibling `json:"siblings"`
	Children              []Child   `json:"children"`
	Comment               string    `json:"comment"`
	ProfilePictureURL     string    `json:"profilePictureURL"`
}

type FileData struct {
	Content     io.Reader
	Header      *multipart.FileHeader
	ContentType string
}

type QueryParams struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type School struct {
	Name        string `json:"name"`
	YearStarted string `json:"yearStarted"`
	YearEnded   string `json:"yearEnded"`
}

// Camp is the internal representation of a camp
type Camp struct {
	Attended  bool   `json:"attended"`
	StartYear string `json:"startYear"`
	EndYear   string `json:"endYear"`
	Specialty string `json:"specialty"`
	Camper    bool   `json:"camper"`
	Counselor bool   `json:"counselor"`
}

// Sibling is the internal representation of a Sibling
type Sibling struct {
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	YearCompleted string `json:"yearCompleted"`
	MiddleSchool  School `json:"middleSchool"`
	HighSchool    School `json:"highSchool"`
}

// Child is the internal representation of a Child
type Child struct {
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	GraduationYear string `json:"graduationYear"`
}

type Address struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}
