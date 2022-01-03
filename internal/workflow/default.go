package workflow

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/email"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/mapping"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/storage"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/token"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
	"github.com/aymerick/raymond"
	"github.com/gocarina/gocsv"
	"github.com/mazen160/go-random"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AddUser adds a new user
func AddUser(insertUser db.InsertUserFunc,
	retrieveUserByEmail db.RetrieveUserByEmailFunc,
	provideTime time.EpochProviderFunc,
	genUUID uuid.GenV4Func) AddUserFunc {
	return func(req pkg.UserRequest) (pkg.User, string, error) {
		log.Printf("Adding new user with email=%v", req.Email)

		u, err := retrieveUserByEmail(req.Email)
		if err == nil {
			return pkg.User{}, "", errors.Errorf("workflow - user already exists with email=%v", u.Email)
		}

		pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to hash password, email=%v", req.Email)
		}

		user := mapping.ToDbUser(req, pw, genUUID, provideTime)
		if err := insertUser(user); err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to insert user into db, email=%v", req.Email)
		}

		uToken, err := token.CreateUserToken(user, provideTime)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to generate JWT token for userId=%v", user.ID)
		}

		return mapping.ToDTOUser(user), uToken, nil
	}
}

// LoginUser logs in a new user
func LoginUser(retrieveUserByEmail db.RetrieveUserByEmailFunc, provideTime time.EpochProviderFunc) LoginUserFunc {
	return func(req pkg.UserRequest) (pkg.User, string, error) {
		log.Printf("Logging in user with email=%v", req.Email)

		user, err := retrieveUserByEmail(req.Email)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to find user with email=%v", req.Email)
		}

		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
			return pkg.User{}, "", errors.Wrap(err, "workflow - password is invalid")
		}

		uToken, err := token.CreateUserToken(user, provideTime)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to generate JWT token for userId=%v", user.ID)
		}

		return mapping.ToDTOUser(user), uToken, nil
	}
}

// AutoLoginUser auto logs in a user
func AutoLoginUser(retrieveUserById db.RetrieveUserByIDFunc, provideTime time.EpochProviderFunc) AutoLoginUserFunc {
	return func(tokenString string) (pkg.User, string, error) {
		log.Printf("Auto logging in user with token=%v", tokenString)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.User{}, "", errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.User{}, "", errors.Wrap(err, "workflow - unable to find user with given token")
		}

		uToken, err := token.CreateUserToken(user, provideTime)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to generate JWT token for userId=%v", user.ID)
		}

		return mapping.ToDTOUser(user), uToken, nil
	}
}

func ApproveUser(retrieveUserById db.RetrieveUserByIDFunc,
	provideTime time.EpochProviderFunc,
	replaceUser db.ReplaceUserFunc) ApproveUserFunc {
	return func(userId string, tokenString string) (pkg.User, error) {
		log.Printf("Approving user with userId=%v", userId)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to find user with given token")
		}

		if !user.Admin {
			return pkg.User{}, errors.Errorf("workflow - userId=%v is not an admin", user.ID)
		}

		userToApprove, err := retrieveUserById(userId)
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to find user to approve")
		}

		userToApprove.Status = internal.ApprovedUserStatus
		userToApprove.LastUpdatedTimestamp = provideTime()

		if err := replaceUser(userToApprove); err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to update user")
		}

		return mapping.ToDTOUser(userToApprove), nil
	}
}

func DenyUser(retrieveUserById db.RetrieveUserByIDFunc,
	provideTime time.EpochProviderFunc,
	replaceUser db.ReplaceUserFunc) DenyUserFunc {
	return func(userId string, tokenString string) (pkg.User, error) {
		log.Printf("Denying user with userId=%v", userId)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to find user with given token")
		}

		if !user.Admin {
			return pkg.User{}, errors.Errorf("workflow - userId=%v is not an admin", user.ID)
		}

		userToDeny, err := retrieveUserById(userId)
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to find user to deny")
		}

		userToDeny.Status = internal.DeniedUserStatus
		userToDeny.LastUpdatedTimestamp = provideTime()

		if err := replaceUser(userToDeny); err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to update user")
		}

		return mapping.ToDTOUser(userToDeny), nil
	}
}

func AddAlumni(retrieveUserById db.RetrieveUserByIDFunc,
	insertAlumni db.InsertAlumniFunc,
	replaceUser db.ReplaceUserFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	provideTime time.EpochProviderFunc,
	genUUID uuid.GenV4Func,
	uploadToS3 storage.UploadImageFunc,
	presignURL storage.GetImageURLFunc,
	sendEmail email.SendEmailFunc) AddAlumniFunc {
	return func(req pkg.AlumniRequest, fileData pkg.FileData, tokenString string, skipFileUpload bool) (pkg.Alumni, error) {
		log.Printf("Adding alumni with details=%+v", req)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.Alumni{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		// If user already has an AlumniID return an error
		if user.AlumniID != "" && !user.Admin {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - userId=%v already has alumniId=%v", user.ID, user.AlumniID)
		}

		// Upload profile picture to S3
		s3Filename := ""
		if !skipFileUpload {
			s3Filename = genUUID().Val()
			if err := uploadToS3(fileData.Content, fileData.ContentType, s3Filename, fileData.Header.Filename); err != nil {
				return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to upload image to S3, userId=%v", user.ID)
			}
		}

		a := mapping.ToDBAlumni(req, s3Filename, provideTime, genUUID)
		if err := insertAlumni(a); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to insert alumni, userId=%v", user.ID)
		}

		if user.Admin {
			return mapping.ToDTOAlumni(a, presignURL, internal.User{}), nil
		}

		// Add the AlumniID to the User
		user.AlumniID = a.ID
		if err := replaceUser(user); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to replace userId=%v", user.ID)
		}

		// Send email
		et, err := getEmailTemplate(internal.NewAlumniTemplateName)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve email template")
		}

		bodyTpl, err := raymond.Parse(et.HTML)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to parse email body template")
		}

		subjectTpl, err := raymond.Parse(et.Subject)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to parse email subject template")
		}

		emailBody, err := bodyTpl.Exec(a)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to exec email body template")
		}

		emailSubject, err := subjectTpl.Exec(a)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to exec email subject template")
		}

		er := email.SendRequest{
			Subject:     emailSubject,
			HTMLContent: emailBody,
			Recipient:   internal.EmailRecipient,
			Sender:      internal.EmailRecipient,
		}

		if err := sendEmail(er); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to send email")
		}

		return mapping.ToDTOAlumni(a, presignURL, internal.User{}), nil
	}
}

func UpdateAlumni(retrieveUserById db.RetrieveUserByIDFunc,
	updateAlumni db.UpdateAlumniFunc,
	retrieveAlumniById db.RetrieveAlumniByIDFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	provideTime time.EpochProviderFunc,
	genUUID uuid.GenV4Func,
	uploadToS3 storage.UploadImageFunc,
	presignURL storage.GetImageURLFunc,
	sendEmail email.SendEmailFunc,
) UpdateAlumniFunc {
	return func(req pkg.UpdateAlumniRequest, alumniId string, fileData pkg.FileData, tokenString string, skipFileUpload bool) (pkg.Alumni, error) {
		log.Printf("Updating alumniId=%v", alumniId)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.Alumni{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		if user.AlumniID != uuid.V4(alumniId) && !user.Admin {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - userId=%v already has alumniId=%v", user.ID, user.AlumniID)
		}

		a, err := retrieveAlumniById(alumniId)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - alumniId=%v does not exist", alumniId)
		}

		s3Filename := a.ProfilePictureKey
		if !skipFileUpload {
			s3Filename = genUUID().Val()
			if err := uploadToS3(fileData.Content, fileData.ContentType, s3Filename, fileData.Header.Filename); err != nil {
				return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to upload image to S3, userId=%v", user.ID)
			}
		}

		updates := mapping.ToAlumniUpdate(req, s3Filename, provideTime)
		if err := updateAlumni(alumniId, updates); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to update alumniId=%v", alumniId)
		}

		alum, err := retrieveAlumniById(alumniId)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve alumniId=%v", alumniId)
		}

		// Send email
		et, err := getEmailTemplate(internal.UpdatedAlumniTemplateName)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve email template")
		}

		bodyTpl, err := raymond.Parse(et.HTML)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to parse email body template")
		}

		subjectTpl, err := raymond.Parse(et.Subject)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to parse email subject template")
		}

		emailBody, err := bodyTpl.Exec(a)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to exec email body template")
		}

		bb, err := json.MarshalIndent(updates, "", "\t")
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to marshal updates")
		}
		emailBody = emailBody + "\n\n" + string(bb)

		emailSubject, err := subjectTpl.Exec(a)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to exec email subject template")
		}

		er := email.SendRequest{
			Subject:     emailSubject,
			HTMLContent: emailBody,
			Recipient:   internal.EmailRecipient,
			Sender:      internal.EmailRecipient,
		}

		if err := sendEmail(er); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to send email")
		}

		return mapping.ToDTOAlumni(alum, presignURL, internal.User{}), nil
	}
}

func RetrieveAlumniByID(retrieveByID db.RetrieveAlumniByIDFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	retrieveUserByAlumniId db.RetrieveUserByAlumniIDFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc) RetrieveAlumniByIDFunc {
	return func(alumniId string, tokenString string) (pkg.AlumniInterface, error) {
		log.Printf("Retrieving alumni with id=%v", alumniId)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.Alumni{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		a, err := retrieveByID(alumniId)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve alumniId=%v", alumniId)
		}

		if user.AlumniID.Val() != alumniId && !user.Admin && !a.IsPublic || user.AlumniID.Val() != alumniId && !user.IsApproved() && !user.Admin {
			return pkg.Alumni{}, errors.Errorf("workflow - userId=%v does not have access to alumniId=%v", user.ID, alumniId)
		}

		// If a user tried to access another user who is public
		if user.AlumniID.Val() != alumniId && !user.Admin && a.IsPublic {
			ca := mapping.ToCleanAlumni(a, presignURL, internal.User{})
			return ca, nil
		}

		aUser, err := retrieveUserByAlumniId(a.ID.Val())
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve user with alumniId=%v", a.ID)
		}

		return mapping.ToDTOAlumni(a, presignURL, aUser), nil
	}
}

func ChangeAlumniPrivacy(retrieveByID db.RetrieveAlumniByIDFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	changePrivacyStatus db.ChangeAlumniPrivacyFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc,
	isPublic bool) ChangeAlumniPrivacyFunc {
	return func(alumniId, tokenString string) (pkg.Alumni, error) {
		log.Printf("Updating privacy status of alumni with id=%v", alumniId)

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return pkg.Alumni{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		if user.AlumniID.Val() != alumniId && !user.Admin {
			return pkg.Alumni{}, errors.Errorf("workflow - userId=%v does not have access to alumniId=%v", user.ID, alumniId)
		}

		if err := changePrivacyStatus(alumniId, isPublic); err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to update alumniId=%v", alumniId)
		}

		a, err := retrieveByID(alumniId)
		if err != nil {
			return pkg.Alumni{}, errors.Wrapf(err, "workflow - unable to retrieve alumniId=%v", alumniId)
		}

		return mapping.ToDTOAlumni(a, presignURL, internal.User{}), nil
	}
}

func RetrieveAlumni(retrieveAlumnis db.RetrieveAllAlumniFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	retrieveUsersAlumniIDs db.RetrieveUsersAlumniIDsFunc,
	retrieveUserByAlumniId db.RetrieveUserByAlumniIDFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc) RetrieveAlumniFunc {
	return func(params pkg.QueryParams, tokenString string) ([]pkg.CleanAlumni, pkg.PageInfo, error) {
		log.Printf("Retrieving all alumni")

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		if !user.Admin && !user.IsApproved() {
			return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Errorf("workflow - userId=%v does not have access to retrieve alumni until they are approved", user.ID)
		}

		alumniIDs, err := retrieveUsersAlumniIDs(params.Status)
		if err != nil {
			return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Wrapf(err, "workflow - unable to retrieve alumni ids")
		}

		aa, pi, err := retrieveAlumnis(params, user.AlumniID.Val(), user.Admin, alumniIDs...)
		if err != nil {
			return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Wrap(err, "workflow - unable to retrieve all alumnis")
		}

		cleanAlumni := []pkg.CleanAlumni{}
		for _, a := range aa {
			aUser, err := retrieveUserByAlumniId(a.ID.Val())
			if err != nil {
				return []pkg.CleanAlumni{}, pkg.PageInfo{}, errors.Wrapf(err, "workflow - unable to retrieve user with alumniId=%v", a.ID)
			}

			cleanAlumni = append(cleanAlumni, mapping.ToCleanAlumni(a, presignURL, aUser))
		}

		return cleanAlumni, pi, nil
	}
}

func ExportCSV(retrieveAlumnis db.RetrieveAllAlumniFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	retrieveUsersAlumniIDs db.RetrieveUsersAlumniIDsFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc) ExportCSVFunc {
	return func(params pkg.QueryParams, tokenString string) ([]byte, error) {
		log.Printf("Exporting CSV of alumni")

		id, _, err := token.CheckUserToken(tokenString, provideTime)
		if err != nil {
			return []byte{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(id.Val())
		if err != nil {
			return []byte{}, errors.Wrapf(err, "workflow - unable to find user with given token, userId=%v", user.ID)
		}

		if !user.Admin {
			return []byte{}, errors.Errorf("workflow - user does not have access to export a CSV")
		}

		alumniIDs, err := retrieveUsersAlumniIDs(params.Status)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "workflow - unable to retrieve alumni ids")
		}

		aa, _, err := retrieveAlumnis(params, user.AlumniID.Val(), user.Admin, alumniIDs...)
		if err != nil {
			return []byte{}, errors.Wrap(err, "workflow - unable to retrieve all alumnis")
		}

		bb, err := gocsv.MarshalBytes(aa)
		if err != nil {
			return []byte{}, errors.Wrap(err, "workflow - unable to marshal to csv")
		}

		return bb, nil
	}
}

func ForgotPassword(retrieveUserByEmail db.RetrieveUserByEmailFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	sendEmail email.SendEmailFunc,
	insertResetPassword db.CreateResetPasswordFunc,
	provideTime time.EpochProviderFunc) ForgotPasswordFunc {
	return func(emailAddress string) error {
		log.Printf("Sending reset password email to %v", emailAddress)

		user, err := retrieveUserByEmail(emailAddress)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to retrieve user with email=%v", emailAddress)
		}

		token, err := random.String(16)
		if err != nil {
			return errors.Wrap(err, "workflow - unable to generate token")
		}

		rp := internal.ResetPassword{
			Email:            user.Email,
			Token:            token,
			CreatedTimestamp: provideTime().ToISO8601().Val(),
		}

		if err := insertResetPassword(rp); err != nil {
			return errors.Wrapf(err, "workflow - unable to insert reset password")
		}

		// Send email
		et, err := getEmailTemplate(internal.ForgotPasswordTemplateName)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to retrieve email template")
		}

		bodyTpl, err := raymond.Parse(et.HTML)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to parse email body template")
		}

		subjectTpl, err := raymond.Parse(et.Subject)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to parse email subject template")
		}

		emailBody, err := bodyTpl.Exec(rp)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to exec email body template")
		}

		emailSubject, err := subjectTpl.Exec(rp)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to exec email subject template")
		}

		er := email.SendRequest{
			Subject:     emailSubject,
			HTMLContent: emailBody,
			Recipient:   user.Email,
			Sender:      internal.EmailRecipient,
		}

		if err := sendEmail(er); err != nil {
			return errors.Wrapf(err, "workflow - unable to send email")
		}

		return nil
	}
}

func SetNewPassword(retrieveResetPassword db.FindResetPasswordFunc,
	deleteResetPasswords db.DeleteResetPasswordsFunc,
	retrieveUserByEmail db.RetrieveUserByEmailFunc,
	replaceUser db.ReplaceUserFunc,
	provideTime time.EpochProviderFunc) SetNewPasswordFunc {
	return func(rp pkg.ResetPassword) (pkg.User, string, error) {
		log.Printf("Setting new password for user with email=%v", rp.Email)

		internalRP, err := retrieveResetPassword(rp.Email, rp.Token)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to retrieve reset password")
		}

		user, err := retrieveUserByEmail(internalRP.Email)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to retrieve user with email=%v", internalRP.Email)
		}

		// Use bcrypt to encrypt password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rp.Password), bcrypt.DefaultCost)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to generate hashed password")
		}

		user.Password = hashedPassword

		if err := replaceUser(user); err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to replace user")
		}

		if err := deleteResetPasswords(internalRP.Email); err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to delete reset password")
		}

		uToken, err := token.CreateUserToken(user, provideTime)
		if err != nil {
			return pkg.User{}, "", errors.Wrapf(err, "workflow - unable to generate JWT token for userId=%v", user.ID)
		}

		return mapping.ToDTOUser(user), uToken, nil
	}
}

func HappyBirthday(retrieveAlumnis db.RetrieveAllAlumniFunc, provideTime time.EpochProviderFunc) HappyBirthdayFunc {
	return func() (pkg.HappyBirthdayResponse, error) {
		iso, err := time.New(provideTime().ToISO8601().Val().Local())
		if err != nil {
			return pkg.HappyBirthdayResponse{}, errors.Wrap(err, "workflow - unable to create iso time")
		}
		ds := iso.DateString()
		m := strings.Split(ds, "-")[1]
		d := strings.Split(ds, "-")[2]

		bday := fmt.Sprintf("%v-%v", m, d)

		log.Printf("Retrieving Alumnis with Birthday=%v", bday)

		qp := pkg.QueryParams{Limit: -1, Birthday: bday}

		aa, _, err := retrieveAlumnis(qp, "", true)
		if err != nil {
			return pkg.HappyBirthdayResponse{}, errors.Wrapf(err, "workflow - unable to retrieve alumnis")
		}

		hbdAA := []pkg.HappyBirthdayAlumni{}
		for _, a := range aa {
			hbdAA = append(hbdAA, mapping.ToHappyBirthdayAlumni(a))
		}

		upcoming := []pkg.HappyBirthdayAlumni{}
		currentDay := provideTime().ToISO8601().Val()
		nextDay := currentDay
		for i := 0; i < 6; i++ {
			nextDay = nextDay.AddDate(0, 0, 1)
			iso, err := time.New(nextDay)
			if err != nil {
				return pkg.HappyBirthdayResponse{}, errors.Wrapf(err, "workflow - unable to create time")
			}
			ds := iso.DateString()
			m := strings.Split(ds, "-")[1]
			d := strings.Split(ds, "-")[2]
			bday := fmt.Sprintf("%v-%v", m, d)
			qp := pkg.QueryParams{Limit: -1, Birthday: bday}

			aa, _, err := retrieveAlumnis(qp, "", true)
			if err != nil {
				return pkg.HappyBirthdayResponse{}, errors.Wrapf(err, "workflow - unable to retrieve alumnis")
			}
			for _, a := range aa {
				fmt.Println(a.ID)
				upcoming = append(upcoming, mapping.ToHappyBirthdayAlumni(a))
			}
		}

		hbdResponse := pkg.HappyBirthdayResponse{
			Today:    hbdAA,
			Upcoming: upcoming,
		}

		return hbdResponse, nil
	}
}

func HappyBirthdayEmail(retrieveAlumnis db.RetrieveAllAlumniFunc,
	provideTime time.EpochProviderFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	retrieveUserByAlumniId db.RetrieveUserByAlumniIDFunc,
	sendEmail email.SendEmailFunc) HappyBirthdayEmailFunc {
	return func() error {
		ds := provideTime().ToISO8601().DateString()
		m := strings.Split(ds, "-")[1]
		d := strings.Split(ds, "-")[2]

		bday := fmt.Sprintf("%v-%v", m, d)
		qp := pkg.QueryParams{Limit: -1, Birthday: bday}

		log.Printf("Sending emails to Alumnis with Birthday=%v", bday)

		aa, _, err := retrieveAlumnis(qp, "", true)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to retrieve alumnis")
		}

		// Send email
		et, err := getEmailTemplate(internal.HappyBirthdayTemplateName)
		if err != nil {
			return errors.Wrapf(err, "workflow - unable to retrieve email template")
		}
		for _, a := range aa {
			user, err := retrieveUserByAlumniId(a.ID.Val())
			if err != nil {
				return errors.Wrapf(err, "workflow - unable to retrieve user with alumniId=%v", a.ID.Val())
			}

			bodyTpl, err := raymond.Parse(et.HTML)
			if err != nil {
				return errors.Wrapf(err, "workflow - unable to parse email body template")
			}

			subjectTpl, err := raymond.Parse(et.Subject)
			if err != nil {
				return errors.Wrapf(err, "workflow - unable to parse email subject template")
			}

			emailBody, err := bodyTpl.Exec(a)
			if err != nil {
				return errors.Wrapf(err, "workflow - unable to exec email body template")
			}

			emailSubject, err := subjectTpl.Exec(a)
			if err != nil {
				return errors.Wrapf(err, "workflow - unable to exec email subject template")
			}

			er := email.SendRequest{
				Subject:     emailSubject,
				HTMLContent: emailBody,
				Recipient:   user.Email,
				Sender:      internal.EmailRecipient,
			}

			if err := sendEmail(er); err != nil {
				return errors.Wrapf(err, "workflow - unable to send email")
			}
		}
		return nil
	}
}
