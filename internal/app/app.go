package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
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
	UpdateAlumniHandler       http.HandlerFunc
}

// Handler turns the App into an http hander
func (a *App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.AddUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", a.LoginUserHandler)
	router.HandlerFunc(http.MethodGet, "/autologin", a.AutoLoginUserHandler)
	router.HandlerFunc(http.MethodPost, "/alumni", a.AddAlumniHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("/alumni/:%v", alumniIdKey), a.UpdateAlumniHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/alumni/:%v", alumniIdKey), a.RetrieveAlumniByIDHandler)
	router.HandlerFunc(http.MethodGet, "/alumni", a.RetrieveAllAlumniHandler)
	h := http.HandlerFunc(router.ServeHTTP)
	return h
}

// OptionalArgs is a representation of all the optional arguments for this application
type OptionalArgs struct {
	EpochTimeProvider   time.EpochProviderFunc
	UUIDGenerator       uuid.GenV4Func
	PhotosS3Bucket      string
	AddUser             db.InsertUserFunc
	RetrieveUserByEmail db.RetrieveUserByEmailFunc
	RetrieveUserByID    db.RetrieveUserByIDFunc
	ReplaceUser         db.ReplaceUserFunc
	InsertAlumni        db.InsertAlumniFunc
	RetrieveAlumniByID  db.RetrieveAlumniByIDFunc
	RetrieveAlumnis     db.RetrieveAllAlumniFunc
	UpdateAlumni        db.UpdateAlumniFunc
	S3Upload            storage.UploadFunc
	S3Presign           storage.PresignFunc
}

// Option is a representation of a function that modifies optional arguments
type Option func(oa *OptionalArgs)

// New creates a new App
func New(provideDb *mongo.Database, opts ...Option) App {
	s3Config := storage.DefaultConfig()

	oa := OptionalArgs{
		EpochTimeProvider:   time.CurrentEpoch,
		UUIDGenerator:       uuid.GenV4,
		PhotosS3Bucket:      os.Getenv("S3_BUCKET"),
		AddUser:             db.InsertUser(provideDb),
		RetrieveUserByEmail: db.RetrieveUserByEmail(provideDb),
		RetrieveUserByID:    db.RetrieveUserByID(provideDb),
		ReplaceUser:         db.ReplaceUser(provideDb),
		InsertAlumni:        db.InsertAlumni(provideDb),
		RetrieveAlumniByID:  db.RetrieveAlumniByID(provideDb),
		RetrieveAlumnis:     db.RetrieveAllAlumni(provideDb),
		UpdateAlumni:        db.UpdateAlumni(provideDb),
		S3Upload:            storage.UploadToS3(s3Config),
		S3Presign:           storage.PresignObject(s3Config),
	}

	for _, opt := range opts {
		opt(&oa)
	}

	uploadImage := storage.UploadImage(oa.S3Upload, oa.PhotosS3Bucket)
	presignURL := storage.GetImageURL(oa.S3Presign, oa.PhotosS3Bucket)

	addUserHandler := AddUserHandler(oa.EpochTimeProvider, oa.UUIDGenerator, oa.AddUser, oa.RetrieveUserByEmail)
	loginUserHandler := LoginUserHandler(oa.RetrieveUserByEmail, oa.EpochTimeProvider)
	autologinUserHandler := AutoLoginUserHandler(oa.RetrieveUserByID, oa.EpochTimeProvider)
	addAlumniHandler := AddAlumniHandler(oa.RetrieveUserByID, oa.InsertAlumni, oa.ReplaceUser, oa.EpochTimeProvider, oa.UUIDGenerator, uploadImage, presignURL)
	updateAlumniHandler := UpdateAlumniHandler(oa.RetrieveUserByID, oa.UpdateAlumni, oa.RetrieveAlumniByID, oa.EpochTimeProvider, oa.UUIDGenerator, uploadImage, presignURL)
	retrieveAlumniByIdHandler := RetrieveAlumniByIDHandler(oa.RetrieveAlumniByID, oa.RetrieveUserByID, oa.EpochTimeProvider, presignURL)
	retrieveAllAlumniHandler := RetrieveAlumniHandler(oa.RetrieveAlumnis, oa.RetrieveUserByID, oa.EpochTimeProvider, presignURL)

	return App{
		AddUserHandler:            addUserHandler,
		LoginUserHandler:          loginUserHandler,
		AutoLoginUserHandler:      autologinUserHandler,
		AddAlumniHandler:          addAlumniHandler,
		RetrieveAlumniByIDHandler: retrieveAlumniByIdHandler,
		RetrieveAllAlumniHandler:  retrieveAllAlumniHandler,
		UpdateAlumniHandler:       updateAlumniHandler,
	}
}
