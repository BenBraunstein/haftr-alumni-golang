package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/email"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/storage"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/workflow"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
	"github.com/gabriel-vasile/mimetype"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const (
	profilePictureKey = "profile"
	authTokenKey      = "Authorization"
	jsonDataKey       = "json"
	alumniIdKey       = "alumniId"
	limitKey          = "limit"
	pageKey           = "page"
	firstnameKey      = "firstname"
	lastnameKey       = "lastname"
)

var (
	allowedMimeTypes = []string{"image/jpeg", "image/png"}
)

// AddUserHandler handles an http request to add a new User
func AddUserHandler(provideTime time.EpochProviderFunc, genUUID uuid.GenV4Func, insertUser db.InsertUserFunc, retrieveUserByEmail db.RetrieveUserByEmailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.UserRequest
		if err := JSONToDTO(&req, w, r); err != nil {
			ServeInternalError(err, w)
			return
		}

		addUser := workflow.AddUser(insertUser, retrieveUserByEmail, provideTime, genUUID)
		user, uToken, err := addUser(req)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		resp := pkg.UserResponse{
			User:  user,
			Token: uToken,
		}

		ServeJSON(resp, w)
	}
}

// LoginUserHandler handles an http request to log in a User
func LoginUserHandler(retrieveUserByEmail db.RetrieveUserByEmailFunc, provideTime time.EpochProviderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.UserRequest
		if err := JSONToDTO(&req, w, r); err != nil {
			ServeInternalError(err, w)
			return
		}
		loginUser := workflow.LoginUser(retrieveUserByEmail, provideTime)
		user, uToken, err := loginUser(req)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		resp := pkg.UserResponse{
			User:  user,
			Token: uToken,
		}

		ServeJSON(resp, w)
	}
}

// AutoLoginUserHandler handles an http request to auto login a user
func AutoLoginUserHandler(retrieveUserById db.RetrieveUserByIDFunc, provideTime time.EpochProviderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r)

		autologin := workflow.AutoLoginUser(retrieveUserById, provideTime)
		user, uToken, err := autologin(token)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		resp := pkg.UserResponse{
			User:  user,
			Token: uToken,
		}

		ServeJSON(resp, w)
	}
}

func AddAlumniHandler(retrieveUserById db.RetrieveUserByIDFunc,
	insertAlumni db.InsertAlumniFunc,
	replaceUser db.ReplaceUserFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	provideTime time.EpochProviderFunc,
	genUUID uuid.GenV4Func,
	uploadToS3 storage.UploadImageFunc,
	presignURL storage.GetImageURLFunc,
	sendEmail email.SendEmailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			ServeInternalError(err, w)
			return
		}

		var req pkg.AlumniRequest
		if err := json.Unmarshal([]byte(r.Form[jsonDataKey][0]), &req); err != nil {
			ServeInternalError(err, w)
			return
		}

		f, fh, fileErr := r.FormFile(profilePictureKey)

		var buf bytes.Buffer
		var tee io.Reader
		var fileData pkg.FileData
		if fileErr == nil {
			tee = io.TeeReader(f, &buf)

			fileData = pkg.FileData{
				Content:     &buf,
				Header:      fh,
				ContentType: getFileContentType(tee, fh.Filename),
			}
		}

		if !isAllowedMimeType(fileData.ContentType) && fileErr == nil {
			ServeInternalError(errors.Errorf("handler - mime type=%v is unsupported", fileData.ContentType), w)
			return
		}

		token := getAuthToken(r)

		addAlum := workflow.AddAlumni(retrieveUserById, insertAlumni, replaceUser, getEmailTemplate, provideTime, genUUID, uploadToS3, presignURL, sendEmail)
		alumni, err := addAlum(req, fileData, token, fileErr != nil)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		ServeJSON(alumni, w)
	}
}

func UpdateAlumniHandler(retrieveUserById db.RetrieveUserByIDFunc,
	updateAlumni db.UpdateAlumniFunc,
	retrieveAlumniById db.RetrieveAlumniByIDFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	provideTime time.EpochProviderFunc,
	genUUID uuid.GenV4Func,
	uploadToS3 storage.UploadImageFunc,
	presignURL storage.GetImageURLFunc,
	sendEmail email.SendEmailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			ServeInternalError(err, w)
			return
		}

		alumId, err := retrieveResourceID(alumniIdKey, r)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		var req pkg.UpdateAlumniRequest
		if err := json.Unmarshal([]byte(r.Form[jsonDataKey][0]), &req); err != nil {
			ServeInternalError(err, w)
			return
		}

		f, fh, fileErr := r.FormFile(profilePictureKey)
		var buf bytes.Buffer
		var tee io.Reader
		var fileData pkg.FileData
		if fileErr == nil {
			tee = io.TeeReader(f, &buf)

			fileData = pkg.FileData{
				Content:     &buf,
				Header:      fh,
				ContentType: getFileContentType(tee, fh.Filename),
			}
		}

		if !isAllowedMimeType(fileData.ContentType) && fileErr == nil {
			ServeInternalError(errors.Errorf("handler - mime type=%v is unsupported", fileData.ContentType), w)
			return
		}

		token := getAuthToken(r)

		updateAlum := workflow.UpdateAlumni(retrieveUserById, updateAlumni, retrieveAlumniById, getEmailTemplate, provideTime, genUUID, uploadToS3, presignURL, sendEmail)
		alumni, err := updateAlum(req, alumId, fileData, token, fileErr != nil)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		ServeJSON(alumni, w)
	}
}

func CorsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Request-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	}
}

func RetrieveAlumniByIDHandler(retrieveByID db.RetrieveAlumniByIDFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r)

		alumId, err := retrieveResourceID(alumniIdKey, r)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		retrieveAlum := workflow.RetrieveAlumniByID(retrieveByID, retrieveUserById, provideTime, presignURL)
		alum, err := retrieveAlum(alumId, token)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		ServeJSON(alum, w)
	}
}

func RetrieveAlumniHandler(retrieveAlumnis db.RetrieveAllAlumniFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r)

		params, err := getQueryParams(r)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		retrieveAlumnis := workflow.RetrieveAlumni(retrieveAlumnis, retrieveUserById, provideTime, presignURL)
		aa, pi, err := retrieveAlumnis(params, token)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		res := pkg.RetrieveCleanAlumniResponse{
			Alumni:   aa,
			PageInfo: pi,
		}

		ServeJSON(res, w)
	}
}

func HappyBirthdayHandler(retrieveAlumnis db.RetrieveAllAlumniFunc, provideTime time.EpochProviderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		happyBirthday := workflow.HappyBirthday(retrieveAlumnis, provideTime)
		aa, err := happyBirthday()
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		ServeJSON(aa, w)
	}
}

func ChangeAlumniPrivacyHandler(retrieveByID db.RetrieveAlumniByIDFunc,
	retrieveUserById db.RetrieveUserByIDFunc,
	changePrivacyStatus db.ChangeAlumniPrivacyFunc,
	provideTime time.EpochProviderFunc,
	presignURL storage.GetImageURLFunc,
	isPublic bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r)

		alumId, err := retrieveResourceID(alumniIdKey, r)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		changeStatus := workflow.ChangeAlumniPrivacy(retrieveByID, retrieveUserById, changePrivacyStatus, provideTime, presignURL, isPublic)
		a, err := changeStatus(alumId, token)
		if err != nil {
			ServeInternalError(err, w)
			return
		}

		ServeJSON(a, w)
	}
}

// JSONToDTO decodes an http request JSON body to a data transfer object
func JSONToDTO(DTO interface{}, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&DTO)
}

// ServeInternalError serves a 500 error
func ServeInternalError(err error, w http.ResponseWriter) {
	var newError struct {
		Message string `json:"message"`
	}
	newError.Message = err.Error()
	bb, err := json.MarshalIndent(newError, "", "\t")
	if err != nil {
		ServeInternalError(err, w)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(bb)
}

//ServeJSON returns a JSON response for an http request
func ServeJSON(res interface{}, w http.ResponseWriter) {
	bb, err := json.Marshal(res)
	if err != nil {
		ServeInternalError(err, w)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bb)
}

// retrieveResourceID retrieves a resource id from an incoming http request
func retrieveResourceID(idKey string, r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	tID := params.ByName(idKey)
	if tID == "" {
		return tID, fmt.Errorf("%s not found in request context", idKey)
	}
	return tID, nil
}

func getQueryParams(r *http.Request) (pkg.QueryParams, error) {
	var lim, page int
	var err error

	if r.URL.Query().Get(limitKey) != "" {
		lim, err = strconv.Atoi(r.URL.Query().Get(limitKey))
		if err != nil {
			return pkg.QueryParams{}, errors.Wrapf(err, "handler - error parsing limit param=%v", r.URL.Query().Get(limitKey))
		}
	}

	if r.URL.Query().Get(pageKey) != "" {
		page, err = strconv.Atoi(r.URL.Query().Get(pageKey))
		if err != nil {
			return pkg.QueryParams{}, errors.Wrapf(err, "handler - error parsing page param=%v", r.URL.Query().Get(pageKey))
		}
	}

	params := pkg.QueryParams{
		Limit:     int64(lim),
		Page:      int64(page),
		Firstname: r.URL.Query().Get(firstnameKey),
		Lastname:  r.URL.Query().Get(lastnameKey),
	}

	if params.Limit == 0 {
		params.Limit = internal.DefaultPageLimit
	}
	if params.Page == 0 {
		params.Page = 1
	}

	return params, nil
}

func getAuthToken(r *http.Request) string {
	authHeader := r.Header.Get(authTokenKey)
	if authHeader != "" {
		return strings.Split(authHeader, " ")[1]
	}
	return ""
}

func getFileContentType(r io.Reader, fn string) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	mimeType, _ := mimetype.DetectReader(buf)
	if mimeType.String() != "application/octet-stream" {
		return mimeType.String()
	}
	extArr := strings.Split(fn, ".")
	ext := extArr[len(extArr)-1]
	switch ext {
	case "jpg":
		return "image/jpeg"
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	}
	return mimeType.String()
}

func isAllowedMimeType(t string) bool {
	for _, m := range allowedMimeTypes {
		if t == m {
			return true
		}
	}
	return false
}
