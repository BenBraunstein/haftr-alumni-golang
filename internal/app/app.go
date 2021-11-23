package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/email"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/storage"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

// App is a representation of an App
type App struct {
	AddUserHandler            http.HandlerFunc
	LoginUserHandler          http.HandlerFunc
	AutoLoginUserHandler      http.HandlerFunc
	AddAlumniHandler          http.HandlerFunc
	RetrieveAlumniByIDHandler http.HandlerFunc
	RetrieveAllAlumniHandler  http.HandlerFunc
	HappyBirthdayHandler      http.HandlerFunc
	UpdateAlumniHandler       http.HandlerFunc
	MakeAlumniPublicHandler   http.HandlerFunc
	MakeAlumniPrivateHandler  http.HandlerFunc
	CorsHandler               http.HandlerFunc
}

// Handler turns the App into an http hander
func (a *App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.AddUserHandler)
	router.HandlerFunc(http.MethodOptions, "/users", a.CorsHandler)
	router.HandlerFunc(http.MethodPost, "/login", a.LoginUserHandler)
	router.HandlerFunc(http.MethodOptions, "/login", a.CorsHandler)
	router.HandlerFunc(http.MethodGet, "/autologin", a.AutoLoginUserHandler)
	router.HandlerFunc(http.MethodOptions, "/autologin", a.CorsHandler)
	router.HandlerFunc(http.MethodPost, "/alumni", a.AddAlumniHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("/alumni/:%v", alumniIdKey), a.UpdateAlumniHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("/alumni/:%v/gopublic", alumniIdKey), a.MakeAlumniPublicHandler)
	router.HandlerFunc(http.MethodOptions, fmt.Sprintf("/alumni/:%v/gopublic", alumniIdKey), a.CorsHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("/alumni/:%v/goprivate", alumniIdKey), a.MakeAlumniPrivateHandler)
	router.HandlerFunc(http.MethodOptions, fmt.Sprintf("/alumni/:%v/goprivate", alumniIdKey), a.CorsHandler)
	router.HandlerFunc(http.MethodOptions, fmt.Sprintf("/alumni/:%v", alumniIdKey), a.CorsHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/alumni/:%v", alumniIdKey), a.RetrieveAlumniByIDHandler)
	router.HandlerFunc(http.MethodGet, "/alumni", a.RetrieveAllAlumniHandler)
	router.HandlerFunc(http.MethodGet, "/happybirthday", a.HappyBirthdayHandler)
	h := http.HandlerFunc(router.ServeHTTP)
	return h
}

// OptionalArgs is a representation of all the optional arguments for this application
type OptionalArgs struct {
	EpochTimeProvider           time.EpochProviderFunc
	UUIDGenerator               uuid.GenV4Func
	PhotosS3Bucket              string
	AddUser                     db.InsertUserFunc
	RetrieveUserByEmail         db.RetrieveUserByEmailFunc
	RetrieveUserByID            db.RetrieveUserByIDFunc
	ReplaceUser                 db.ReplaceUserFunc
	InsertAlumni                db.InsertAlumniFunc
	RetrieveAlumniByID          db.RetrieveAlumniByIDFunc
	RetrieveAlumnis             db.RetrieveAllAlumniFunc
	UpdateAlumni                db.UpdateAlumniFunc
	ChangeAlumniPrivacyStatus   db.ChangeAlumniPrivacyFunc
	RetrieveEmailTemplateByName db.RetrieveEmailTemplateByNameFunc
	S3Upload                    storage.UploadFunc
	S3Presign                   storage.PresignFunc
	SendEmail                   email.SendEmailFunc
}

// Option is a representation of a function that modifies optional arguments
type Option func(oa *OptionalArgs)

// New creates a new App
func New(provideDb *mongo.Database, opts ...Option) App {
	s3Config := storage.DefaultConfig()
	sesConfig := email.DefaultConfig()

	oa := OptionalArgs{
		EpochTimeProvider:           time.CurrentEpoch,
		UUIDGenerator:               uuid.GenV4,
		PhotosS3Bucket:              os.Getenv("S3_BUCKET"),
		AddUser:                     db.InsertUser(provideDb),
		RetrieveUserByEmail:         db.RetrieveUserByEmail(provideDb),
		RetrieveUserByID:            db.RetrieveUserByID(provideDb),
		ReplaceUser:                 db.ReplaceUser(provideDb),
		InsertAlumni:                db.InsertAlumni(provideDb),
		RetrieveAlumniByID:          db.RetrieveAlumniByID(provideDb),
		RetrieveAlumnis:             db.RetrieveAllAlumni(provideDb),
		UpdateAlumni:                db.UpdateAlumni(provideDb),
		ChangeAlumniPrivacyStatus:   db.ChangeAlumniPrivacy(provideDb),
		RetrieveEmailTemplateByName: db.RetrieveEmailTemplateByName(provideDb),
		S3Upload:                    storage.UploadToS3(s3Config),
		S3Presign:                   storage.PresignObject(s3Config),
		SendEmail:                   email.SendEmail(sesConfig),
	}

	for _, opt := range opts {
		opt(&oa)
	}

	uploadImage := storage.UploadImage(oa.S3Upload, oa.PhotosS3Bucket)
	presignURL := storage.GetImageURL(oa.S3Presign, oa.PhotosS3Bucket)

	addUserHandler := AddUserHandler(oa.EpochTimeProvider, oa.UUIDGenerator, oa.AddUser, oa.RetrieveUserByEmail)
	loginUserHandler := LoginUserHandler(oa.RetrieveUserByEmail, oa.EpochTimeProvider)
	autologinUserHandler := AutoLoginUserHandler(oa.RetrieveUserByID, oa.EpochTimeProvider)
	addAlumniHandler := AddAlumniHandler(oa.RetrieveUserByID, oa.InsertAlumni, oa.ReplaceUser, oa.RetrieveEmailTemplateByName, oa.EpochTimeProvider, oa.UUIDGenerator, uploadImage, presignURL, oa.SendEmail)
	updateAlumniHandler := UpdateAlumniHandler(oa.RetrieveUserByID, oa.UpdateAlumni, oa.RetrieveAlumniByID, oa.RetrieveEmailTemplateByName, oa.EpochTimeProvider, oa.UUIDGenerator, uploadImage, presignURL, oa.SendEmail)
	corsHandler := CorsHandler()
	retrieveAlumniByIdHandler := RetrieveAlumniByIDHandler(oa.RetrieveAlumniByID, oa.RetrieveUserByID, oa.EpochTimeProvider, presignURL)
	retrieveAllAlumniHandler := RetrieveAlumniHandler(oa.RetrieveAlumnis, oa.RetrieveUserByID, oa.EpochTimeProvider, presignURL)
	happyBirthdayHandler := HappyBirthdayHandler(oa.RetrieveAlumnis, oa.EpochTimeProvider)
	makeAlumniPublicHandler := ChangeAlumniPrivacyHandler(oa.RetrieveAlumniByID, oa.RetrieveUserByID, oa.ChangeAlumniPrivacyStatus, oa.EpochTimeProvider, presignURL, true)
	makeAlumniPrivateHandler := ChangeAlumniPrivacyHandler(oa.RetrieveAlumniByID, oa.RetrieveUserByID, oa.ChangeAlumniPrivacyStatus, oa.EpochTimeProvider, presignURL, false)

	return App{
		AddUserHandler:            addUserHandler,
		LoginUserHandler:          loginUserHandler,
		AutoLoginUserHandler:      autologinUserHandler,
		AddAlumniHandler:          addAlumniHandler,
		RetrieveAlumniByIDHandler: retrieveAlumniByIdHandler,
		RetrieveAllAlumniHandler:  retrieveAllAlumniHandler,
		HappyBirthdayHandler:      happyBirthdayHandler,
		UpdateAlumniHandler:       updateAlumniHandler,
		MakeAlumniPublicHandler:   makeAlumniPublicHandler,
		MakeAlumniPrivateHandler:  makeAlumniPrivateHandler,
		CorsHandler:               corsHandler,
	}
}
