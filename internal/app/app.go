package app

import (
	"net/http"
	"os"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

// App is a representation of an App
type App struct {
	AddUserHandler       http.HandlerFunc
	LoginUserHandler     http.HandlerFunc
	AutoLoginUserHandler http.HandlerFunc
}

// Handler turns the App into an http hander
func (a *App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.AddUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", a.LoginUserHandler)
	router.HandlerFunc(http.MethodGet, "/autologin", a.AutoLoginUserHandler)
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
}

// Option is a representation of a function that modifies optional arguments
type Option func(oa *OptionalArgs)

// New creates a new App
func New(provideDb *mongo.Database, opts ...Option) App {
	oa := OptionalArgs{
		EpochTimeProvider:   time.CurrentEpoch,
		UUIDGenerator:       uuid.GenV4,
		PhotosS3Bucket:      os.Getenv("PHOTOS_S3_BUCKET"),
		AddUser:             db.InsertUser(provideDb),
		RetrieveUserByEmail: db.RetrieveUserByEmail(provideDb),
		RetrieveUserByID:    db.RetrieveUserByID(provideDb),
	}

	for _, opt := range opts {
		opt(&oa)
	}

	addUserHandler := AddUserHandler(oa.EpochTimeProvider, oa.UUIDGenerator, oa.AddUser, oa.RetrieveUserByEmail)
	loginUserHandler := LoginUserHandler(oa.RetrieveUserByEmail)
	autologinUserHandler := AutoLoginUserHandler(oa.RetrieveUserByID)

	return App{
		AddUserHandler:       addUserHandler,
		LoginUserHandler:     loginUserHandler,
		AutoLoginUserHandler: autologinUserHandler,
	}
}
