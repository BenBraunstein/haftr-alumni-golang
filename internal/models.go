package internal

import (
	gotime "time"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
)

const (
	DefaultPageLimit           = 20
	EmailRecipient             = "lifecycles@haftr.org"
	NoReplyEmailAddress        = "no-reply@haftralumni.org"
	PendingUserStatus          = "PENDING"
	ApprovedUserStatus         = "APPROVED"
	DeniedUserStatus           = "DENIED"
	NewAlumniTemplateName      = "NEW_ALUMNI"
	UpdatedAlumniTemplateName  = "UPDATED_ALUMNI"
	ForgotPasswordTemplateName = "FORGOT_PASSWORD"
	HappyBirthdayTemplateName  = "HAPPY_BIRTHDAY"
)

// User is the internal representation of a user
type User struct {
	ID                   uuid.V4    `bson:"id"`
	Email                string     `bson:"email"`
	Password             []byte     `bson:"password"`
	AlumniID             uuid.V4    `bson:"alumniId"`
	Admin                bool       `bson:"admin"`
	Status               string     `bson:"status"`
	CreatedTimestamp     time.Epoch `bson:"createdTimestamp"`
	LastUpdatedTimestamp time.Epoch `bson:"lastUpdatedTimestamp"`
}

// Alumni is the internal representation of an Alumni
type Alumni struct {
	ID                     uuid.V4       `bson:"id"`
	Title                  string        `bson:"title"`
	Firstname              string        `bson:"firstname"`
	Middlename             string        `bson:"middlename"`
	Lastname               string        `bson:"lastname"`
	MarriedName            string        `bson:"marriedName"`
	MaidenName             string        `bson:"maidenName"`
	MotherName             string        `bson:"motherName"`
	MotherDeceased         bool          `bson:"motherDeceased"`
	FatherName             string        `bson:"fatherName"`
	FatherDeceased         bool          `bson:"fatherDeceased"`
	SpouseName             string        `bson:"spouseName"`
	SpouseMaidenName       string        `bson:"spouseMaidenName"`
	CurrentAddress         Address       `bson:"address"`
	HomePhone              string        `bson:"homePhone"`
	CellPhone              string        `bson:"cellPhone"`
	WorkPhone              string        `bson:"workPhone"`
	EmailAddress           string        `bson:"emailAddress"`
	MiddleSchool           School        `bson:"middleschool"`
	HighSchool             School        `bson:"highschool"`
	IsraelSchool           School        `bson:"israelSchool"`
	CollegeAttended        School        `bson:"collegeAttended"`
	GradSchools            []School      `bson:"gradSchools"`
	Profession             []string      `bson:"profession"`
	Birthday               string        `bson:"birthday"`
	Clubs                  []string      `bson:"clubs"`
	SportsTeams            []string      `bson:"sportsTeams"`
	Awards                 []string      `bson:"awards"`
	Committees             []string      `bson:"committees"`
	OldAddresses           []Address     `bson:"oldAddresses"`
	HillelDayCamp          Camp          `bson:"hillelDayCamp"`
	HillelSleepCamp        Camp          `bson:"hillelSleepCamp"`
	HiliDayCamp            Camp          `bson:"hiliDayCamp"`
	HiliWhiteCamp          Camp          `bson:"hiliWhiteCamp"`
	HiliInternationalCamp  Camp          `bson:"hiliInternationalCamp"`
	HILI                   bool          `bson:"hili"`
	HILLEL                 bool          `bson:"hillel"`
	HAFTR                  bool          `bson:"haftr"`
	ParentOfStudent        bool          `bson:"parentOfStudent"`
	Boards                 []string      `bson:"boards"`
	AlumniPositions        []string      `bson:"alumniPositions"`
	Siblings               []Sibling     `bson:"siblings"`
	Children               []Child       `bson:"children"`
	Grandparents           []Grandparent `bson:"grandparents"`
	ClassPresident         bool          `bson:"classPresident"`
	BoardOfTrustees        bool          `bson:"boardOfTrustees"`
	BoardOfEducation       bool          `bson:"boardOfEducation"`
	BoardsComment          string        `bson:"boardsComment"`
	AlumniNewsletters      bool          `bson:"alumniNewsletters"`
	CommunicationsOutreach bool          `bson:"communicationsOutreach"`
	ClassReunions          bool          `bson:"classReunions"`
	AlumniEvents           bool          `bson:"alumniEvents"`
	FundraisingNetworking  bool          `bson:"fundraisingNetworking"`
	DbResearch             bool          `bson:"dbResearch"`
	AlumniChoir            bool          `bson:"alumniChoir"`
	Comment                string        `bson:"comment"`
	IsPublic               bool          `bson:"isPublic"`
	ProfilePictureKey      string        `bson:"profilePictureKey"`
	CreatedTimestamp       time.Epoch    `bson:"createdTimestamp"`
	LastUpdatedTimestamp   time.Epoch    `bson:"lastUpdatedTimestamp"`
}

type ResetPassword struct {
	Email            string      `bson:"email"`
	Token            string      `bson:"token"`
	CreatedTimestamp gotime.Time `bson:"createdTimestamp"`
}

type UpdateAlumniRequest struct {
	Title                  string        `bson:"title,omitempty"`
	Firstname              string        `bson:"firstname,omitempty"`
	Middlename             string        `bson:"middlename,omitempty"`
	Lastname               string        `bson:"lastname,omitempty"`
	MarriedName            string        `bson:"marriedName,omitempty"`
	MaidenName             string        `bson:"maidenName,omitempty"`
	MotherName             string        `bson:"motherName,omitempty"`
	MotherDeceased         bool          `bson:"motherDeceased,omitempty"`
	FatherName             string        `bson:"fatherName,omitempty"`
	FatherDeceased         bool          `bson:"fatherDeceased,omitempty"`
	SpouseName             string        `bson:"spouseName,omitempty"`
	SpouseMaidenName       string        `bson:"spouseMaidenName,omitempty"`
	CurrentAddress         Address       `bson:"address,omitempty"`
	HomePhone              string        `bson:"homePhone,omitempty"`
	CellPhone              string        `bson:"cellPhone,omitempty"`
	WorkPhone              string        `bson:"workPhone,omitempty"`
	EmailAddress           string        `bson:"emailAddress,omitempty"`
	MiddleSchool           School        `bson:"middleschool,omitempty"`
	HighSchool             School        `bson:"highschool,omitempty"`
	IsraelSchool           School        `bson:"israelSchool,omitempty"`
	CollegeAttended        School        `bson:"collegeAttended,omitempty"`
	GradSchools            []School      `bson:"gradSchools,omitempty"`
	Profession             []string      `bson:"profession,omitempty"`
	Birthday               string        `bson:"birthday,omitempty"`
	Clubs                  []string      `bson:"clubs,omitempty"`
	SportsTeams            []string      `bson:"sportsTeams,omitempty"`
	Awards                 []string      `bson:"awards,omitempty"`
	Committees             []string      `bson:"committees,omitempty"`
	OldAddresses           []Address     `bson:"oldAddresses,omitempty"`
	HillelDayCamp          Camp          `bson:"hillelDayCamp,omitempty"`
	HillelSleepCamp        Camp          `bson:"hillelSleepCamp,omitempty"`
	HiliDayCamp            Camp          `bson:"hiliDayCamp,omitempty"`
	HiliWhiteCamp          Camp          `bson:"hiliWhiteCamp,omitempty"`
	HiliInternationalCamp  Camp          `bson:"hiliInternationalCamp,omitempty"`
	HILI                   bool          `bson:"hili,omitempty"`
	HILLEL                 bool          `bson:"hillel,omitempty"`
	HAFTR                  bool          `bson:"haftr,omitempty"`
	ParentOfStudent        bool          `bson:"parentOfStudent,omitempty"`
	Boards                 []string      `bson:"boards,omitempty"`
	AlumniPositions        []string      `bson:"alumniPositions,omitempty"`
	Siblings               []Sibling     `bson:"siblings,omitempty"`
	Children               []Child       `bson:"children,omitempty"`
	Grandparents           []Grandparent `bson:"grandparents,omitempty"`
	ClassPresident         bool          `bson:"classPresident,omitempty"`
	BoardOfTrustees        bool          `bson:"boardOfTrustees,omitempty"`
	BoardOfEducation       bool          `bson:"boardOfEducation,omitempty"`
	BoardsComment          string        `bson:"boardsComment,omitempty"`
	AlumniNewsletters      bool          `bson:"alumniNewsletters,omitempty"`
	CommunicationsOutreach bool          `bson:"communicationsOutreach,omitempty"`
	ClassReunions          bool          `bson:"classReunions,omitempty"`
	AlumniEvents           bool          `bson:"alumniEvents,omitempty"`
	FundraisingNetworking  bool          `bson:"fundraisingNetworking,omitempty"`
	DbResearch             bool          `bson:"dbResearch,omitempty"`
	AlumniChoir            bool          `bson:"alumniChoir,omitempty"`
	Comment                string        `bson:"comment,omitempty"`
	IsPublic               bool          `bson:"isPublic,omitempty"`
	ProfilePictureKey      string        `bson:"profilePictureKey"`
	LastUpdatedTimestamp   time.Epoch    `bson:"lastUpdatedTimestamp"`
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
	Deceased      bool   `bson:"deceased"`
}

// Child is the internal representation of a Child
type Child struct {
	Firstname      string `bson:"firstname"`
	Lastname       string `bson:"lastname"`
	GraduationYear string `bson:"graduationYear"`
	Deceased       bool   `bson:"deceased"`
}

// Grandparent is the internal representation of a set of Grandparents
type Grandparent struct {
	GrandfatherFirstname string `bson:"grandfatherFirstname"`
	GrandfatherDeceased  bool   `bson:"grandfatherDeceased"`
	GrandmotherFirstname string `bson:"grandmotherFirstname"`
	GrandmotherDeceased  bool   `bson:"grandmotherDeceased"`
	Lastname             string `bson:"lastname"`
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

func (u User) IsApproved() bool {
	if u.Status == ApprovedUserStatus {
		return true
	}
	return false
}
