package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"
)

func (rt *_router) getProjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) ||
	//	!helpers.SendNotFoundIfBanned(rt.db, ctx.Auth.GetUserID(), uid, w, rt.baseLogger) {
	//	return
	//}

	// Get user profile
	links, err := rt.db.GetProjects()

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetLinks", w, rt.baseLogger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(links)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}

func (rt *_router) getProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) ||
	//	!helpers.SendNotFoundIfBanned(rt.db, ctx.Auth.GetUserID(), uid, w, rt.baseLogger) {
	//	return
	//}

	// ps.ByName("id") returns a string, so we need to convert it to an int64
	id := ps.ByName("id")

	// Get user profile
	links, err := rt.db.GetProject(id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetLinks", w, rt.baseLogger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(links)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}
