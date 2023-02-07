package api

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"
	"github.com/notherealmarco/WhaleDeployer/service/structures"
)

func (rt *_router) postProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//// check if the user is changing his own username
	//if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
	//	return
	//}

	// decode request body
	var req structures.Project

	if !helpers.DecodeJsonOrBadRequest(r.Body, w, &req, rt.baseLogger) {
		return
	}

	// here we should do some validity checks on the input

	if req.DeployKey {
		// generate the key using openssh

		// maybe , "-b", "4096"

		err := os.MkdirAll(rt.keysPath+"/"+req.Name, 0700)

		if err != nil {
			helpers.SendInternalError(err, "Error creating directory: "+err.Error(), w, rt.baseLogger)
			return
		}

		//check if file already exists
		_, err = os.Stat(rt.keysPath + "/" + req.Name + "/key")

		if err != nil {
			_, err = exec.Command("ssh-keygen", "-t", "ed25519", "-C", "deploy@overlinks", "-f", rt.keysPath+"/"+req.Name+"/key", "-N", "").Output()

			if err != nil {
				helpers.SendBadRequestError(err, "Error generating key: "+err.Error(), w, rt.baseLogger)
				return
			}
		}

	}

	err := rt.db.AddProject(&req)

	if err != nil {
		helpers.SendInternalError(err, "Database error: UpdateUsername", w, rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if req.DeployKey {
		key, err := os.ReadFile(rt.keysPath + "/" + req.Name + "/key.pub")

		if err != nil {
			helpers.SendInternalError(err, "Error reading key: "+err.Error(), w, rt.baseLogger)
			return
		}

		w.Header().Set("content-type", "text/plain")
		_, _ = w.Write(key)
	}
}
