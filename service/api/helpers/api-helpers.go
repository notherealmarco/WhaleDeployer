package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/notherealmarco/WhaleDeployer/service/database"
	"github.com/notherealmarco/WhaleDeployer/service/structures"
	"github.com/sirupsen/logrus"
)

// Tries to decode a json, if it fails, it returns Bad Request to the client and the function returns false
// Otherwise it returns true without sending anything to the client
func DecodeJsonOrBadRequest(r io.Reader, w http.ResponseWriter, v interface{}, l logrus.FieldLogger) bool {

	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		SendInternalError(err, "Error decoding json", w, l)
		return false
	}
	return true
}

// This is not needed in overLinks

//// Verifies if a user exists, if it doesn't, it returns Not Found to the client and the function returns false
//// Otherwise it returns true without sending anything to the client
//func VerifyUserOrNotFound(db database.AppDatabase, uid string, w http.ResponseWriter, l logrus.FieldLogger) bool {
//
//	user_exists, err := db.UserExists(uid)
//
//	if err != nil {
//		SendInternalError(err, "Error verifying user existence", w, l)
//		return false
//	}
//
//	if !user_exists {
//		SendNotFound(w, "User not found", l)
//		return false
//	}
//	return true
//}

// Sends a generic status response
// The response is a json object with a "status" field desribing the status of a request
func SendStatus(httpStatus int, w http.ResponseWriter, description string, l logrus.FieldLogger) {
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

// Sends a Not Found error to the client
func SendNotFound(w http.ResponseWriter, description string, l logrus.FieldLogger) {
	SendStatus(http.StatusNotFound, w, description, l)
}

// Sends a Bad Request error to the client
func SendBadRequest(w http.ResponseWriter, description string, l logrus.FieldLogger) {
	SendStatus(http.StatusBadRequest, w, description, l)
}

// Sends a Bad Request error to the client and logs the given error
func SendBadRequestError(err error, description string, w http.ResponseWriter, l logrus.FieldLogger) {
	l.WithError(err).Error(description)
	SendBadRequest(w, description, l)
}

// Sends an Internal Server Error to the client and logs the given error
func SendInternalError(err error, description string, w http.ResponseWriter, l logrus.FieldLogger) {
	l.WithError(err).Error(description)
	SendStatus(http.StatusInternalServerError, w, description, l)
}

// Tries to roll back a transaction, if it fails it logs the error
func RollbackOrLogError(tx database.DBTransaction, l logrus.FieldLogger) {
	err := tx.Rollback()
	if err != nil {
		l.WithError(err).Error("Error rolling back transaction")
	}
}

// This is a WASAPhoto function, but it's not needed in overLinks

//// Checks if a user is banned by another user, then it returns Not Found to the client and the function returns false
//// Otherwise it returns true whithout sending anything to the client
//func SendNotFoundIfBanned(db database.AppDatabase, uid string, banner string, w http.ResponseWriter, l logrus.FieldLogger) bool {
//	banned, err := db.IsBanned(uid, banner)
//	if err != nil {
//		SendInternalError(err, "Database error: IsBanned", w, l)
//		return false
//	}
//	if banned {
//		SendNotFound(w, "User not found", l)
//		return false
//	}
//	return true
//}
