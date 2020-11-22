package app

import (
	"encoding/json"
	"net/http"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/workflow"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

func AddUserHandler(provideTime time.EpochProviderFunc, genUUID uuid.GenV4Func, insertUser db.InsertUserFunc, retrieveUserByEmail db.RetrieveUserByEmailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.UserRequest
		if err := JSONToDTO(&req, w, r); err != nil {
			ServeInternalError(err, w)
			return
		}
		addUser := workflow.AddUser(insertUser, retrieveUserByEmail, provideTime, genUUID)
		user, err := addUser(req)
		if err != nil {
			ServeInternalError(err, w)
			return
		}
		ServeJSON(user, w)
	}
}

func LoginUserHandler(retrieveUserByEmail db.RetrieveUserByEmailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.UserRequest
		if err := JSONToDTO(&req, w, r); err != nil {
			ServeInternalError(err, w)
			return
		}
		loginUser := workflow.LoginUser(retrieveUserByEmail)
		user, err := loginUser(req)
		if err != nil {
			ServeInternalError(err, w)
			return
		}
		ServeJSON(user, w)
	}
}

// JSONToDTO decodes an http request JSON body to a data transfer object
func JSONToDTO(DTO interface{}, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&DTO)
}

func ServeInternalError(err error, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

//ServeJSON returns a JSON response for an http request
func ServeJSON(res interface{}, w http.ResponseWriter) {
	bb, err := json.Marshal(res)
	if err != nil {
		ServeInternalError(err, w)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bb)
}
