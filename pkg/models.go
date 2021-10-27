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
	IsPublic              bool      `json:"isPublic"`
	ProfilePictureURL     string    `json:"profilePictureURL"`
}

type CleanAlumni struct {
	ID                 uuid.V4 `json:"id"`
	Firstname          string  `json:"firstname"`
	Lastname           string  `json:"lastname"`
	HighSchoolGradYear string  `json:"highSchoolGradYear"`
	CellPhone          string  `json:"cellPhone"`
	EmailAddress       string  `json:"emailAddress"`
	ProfilePictureURL  string  `json:"profilePictureURL"`
}

type RetrieveCleanAlumniResponse struct {
	Alumni   []CleanAlumni `json:"alumni"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type FileData struct {
	Content     io.Reader
	Header      *multipart.FileHeader
	ContentType string
}

type QueryParams struct {
	Limit     int64  `json:"limit"`
	Page      int64  `json:"page"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
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

// PageInfo returns page info for a result
type PageInfo struct {
	CurrentPage int64 `json:"currentPage"`
	LastPage    int64 `json:"lastPage"`
}
