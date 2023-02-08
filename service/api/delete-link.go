package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"
)

func (rt *_router) deleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//// check if the user is changing his own username
	//if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
	//	return
	//}

	id := ps.ByName("id")

	err := rt.db.DeleteProject(id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: UpdateUsername", w, rt.baseLogger)
		return
	}

	helpers.SendStatus(http.StatusOK, w, "Resource deleted", rt.baseLogger)
}
