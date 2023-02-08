package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (rt *_router) buildComposeProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//// check if the user is changing his own username
	//if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
	//	return
	//}

	id := ps.ByName("id")

	p, err := rt.db.GetProject(id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetProjects", w, rt.baseLogger)
		return
	}

	err = rt.db.BuildProject(p.Name)

	if err != nil {
		helpers.SendInternalError(err, "Database error: BuildProject", w, rt.baseLogger)
		return
	}

	fo, err := os.Create("/tmp/" + p.Name + ".log")

	if err != nil {
		helpers.SendInternalError(err, "Error creating log file: "+err.Error(), w, rt.baseLogger)
		rt.db.BuildFail(p.Name)
		return
	}

	defer fo.Close()

	fo.Write([]byte("Project: " + p.Name + "\n\n"))

	// git

	var publicKey ssh.AuthMethod = nil

	if p.DeployKey {

		// populate ~/.ssh/known_hosts

		helpers.ExecOrWriteError("rm $HOME/.ssh/known_hosts", fo, rt.baseLogger)
		helpers.ExecOrWriteError("mkdir $HOME/.ssh", fo, rt.baseLogger)

		ssh_host := strings.Split((strings.Split(p.GitURL, "@")[1]), ":")[0]
		rt.baseLogger.Info("ssh_host: " + ssh_host)
		helpers.ExecOrWriteError("ssh-keyscan "+ssh_host+" >> $HOME/.ssh/known_hosts", fo, rt.baseLogger)

		publicKey, err = ssh.NewPublicKeysFromFile("git", rt.keysPath+"/"+p.Name+"/key", "")

		// clone or pull

		if err = helpers.CloneOrPullSSH(p, publicKey, fo); err != nil {
			rt.baseLogger.Error("Error cloning or pulling: " + err.Error())
			rt.db.BuildFail(p.Name)
			return
		}

	} else {
		// clone or pull

		if err = helpers.CloneOrPull(p, fo); err != nil {
			rt.baseLogger.Error("Error cloning or pulling: " + err.Error())
			rt.db.BuildFail(p.Name)
			return
		}
	}

	// build the docker image if requested
	if p.Dockerfile != "" {
		cmd := "docker build -t " + p.ImageName + ":" + p.ImageTag + " -f " + p.Path + "/" + p.Dockerfile + " " + p.Path

		if !helpers.ExecOrFail(cmd, fo, rt.baseLogger, &w, &rt.db, p.Name) {
			return
		}
	}

	cmd := "docker-compose -f " + p.Path + "/docker-compose.yml down"

	if !helpers.ExecOrFail(cmd, fo, rt.baseLogger, &w, &rt.db, p.Name) {
		return
	}

	cmd = "docker-compose -f " + p.Path + "/docker-compose.yml up --build -d"

	if !helpers.ExecOrFail(cmd, fo, rt.baseLogger, &w, &rt.db, p.Name) {
		return
	}

	// I have written the logs on the file, now I have to send them back to the user

	_ = helpers.WriteResponse(http.StatusOK, &w, fo)
	rt.db.BuildSuccess(p.Name)
}
