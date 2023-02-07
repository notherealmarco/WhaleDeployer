package api

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"
)

func (rt *_router) getLogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) ||
	//	!helpers.SendNotFoundIfBanned(rt.db, ctx.Auth.GetUserID(), uid, w, rt.baseLogger) {
	//	return
	//}

	// ps.ByName("id") returns a string, so we need to convert it to an int64
	id := ps.ByName("id")

	// Get user profile
	p, err := rt.db.GetProject(id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetProject", w, rt.baseLogger)
		return
	}

	logfile := "/tmp/" + p.Name + ".log"

	logs, err := os.ReadFile(logfile)

	if err != nil {
		helpers.SendInternalError(err, "Error reading logs", w, rt.baseLogger)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(logs)
}
