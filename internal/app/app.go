package app

import (
	"net/http"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

// App is a representation of an App
type App struct {
	AddUserHandler   http.HandlerFunc
	LoginUserHandler http.HandlerFunc
}

// Handler turns the App into an http hander
func (a *App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/users", a.AddUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", a.LoginUserHandler)
	h := http.HandlerFunc(router.ServeHTTP)
	return h
}

// OptionalArgs is a representation of all the optional arguments for this application
type OptionalArgs struct {
	EpochTimeProvider   time.EpochProviderFunc
	UUIDGenerator       uuid.GenV4Func
	AddUser             db.InsertUserFunc
	RetrieveUserByEmail db.RetrieveUserByEmailFunc
}

// Option is a representation of a function that modifies optional arguments
type Option func(oa *OptionalArgs)

// New creates a new App
func New(provideDb *mongo.Database, opts ...Option) App {
	oa := OptionalArgs{
		EpochTimeProvider:   time.CurrentEpoch,
		UUIDGenerator:       uuid.GenV4,
		AddUser:             db.InsertUser(provideDb),
		RetrieveUserByEmail: db.RetrieveUserByEmail(provideDb),
	}

	for _, opt := range opts {
		opt(&oa)
	}

	addUserHandler := AddUserHandler(oa.EpochTimeProvider, oa.UUIDGenerator, oa.AddUser, oa.RetrieveUserByEmail)
	loginUserHandler := LoginUserHandler(oa.RetrieveUserByEmail)

	return App{
		AddUserHandler:   addUserHandler,
		LoginUserHandler: loginUserHandler,
	}
}
