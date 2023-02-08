package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WhaleDeployer/service/api/helpers"
	"github.com/notherealmarco/WhaleDeployer/service/api/reqcontext"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (rt *_router) buildComposeProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// we still have no auth

	//// check if the user is changing his own username
	//if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
	//	return
	//}

	id := ps.ByName("id")

	// here we should do some validity checks on the input

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

	var publicKey ssh.AuthMethod = nil

	if p.DeployKey {

		// fix ssh known_hosts
		exec.Command("sh", "-c", "rm $HOME/.ssh/known_hosts").Output()
		exec.Command("sh", "-c", "mkdir $HOME/.ssh").Output()

		ssh_host := strings.Split((strings.Split(p.GitURL, "@")[1]), ":")[0]

		rt.baseLogger.Info("ssh_host: " + ssh_host)

		cmd := "ssh-keyscan " + ssh_host + " >> $HOME/.ssh/known_hosts"
		fo.Write([]byte("\n\n# " + cmd + "\n"))
		c := exec.Command("sh", "-c", cmd)
		c.Stdout = fo
		err = c.Run()
		//

		if err != nil {
			helpers.SendBadRequestError(err, "Error populating known_hosts: "+err.Error(), w, rt.baseLogger)
			rt.db.BuildFail(p.Name)
			return
		}

		publicKey, err = ssh.NewPublicKeysFromFile("git", rt.keysPath+"/"+p.Name+"/key", "")

		if err != nil {
			helpers.SendBadRequestError(err, "Error parsing the key: "+err.Error(), w, rt.baseLogger)
			rt.db.BuildFail(p.Name)
			return
		}

		_, err = gogit.PlainClone(p.Path, false, &gogit.CloneOptions{
			Auth:          publicKey,
			URL:           p.GitURL,
			Progress:      fo,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
		})
	} else {
		_, err = gogit.PlainClone(p.Path, false, &gogit.CloneOptions{
			URL:           p.GitURL,
			Progress:      fo,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
		})
	}

	if errors.Is(err, gogit.ErrRepositoryAlreadyExists) {
		// do a git pull

		r, err := gogit.PlainOpen(p.Path)

		if err != nil {
			helpers.SendBadRequestError(err, "Git error opening repository: "+err.Error(), w, rt.baseLogger)
			rt.db.BuildFail(p.Name)
			return
		}

		wt, err := r.Worktree()

		if err != nil {
			helpers.SendBadRequestError(err, "Git error getting worktree: "+err.Error(), w, rt.baseLogger)
			rt.db.BuildFail(p.Name)
			return
		}

		if p.DeployKey {
			err = wt.Pull(&gogit.PullOptions{
				RemoteName:    "origin",
				Auth:          publicKey,
				Progress:      fo,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
			})
		} else {
			err = wt.Pull(&gogit.PullOptions{
				RemoteName:    "origin",
				Progress:      fo,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", p.GitBranch)),
			})
		}

		if err != nil && err != gogit.NoErrAlreadyUpToDate {
			helpers.SendBadRequestError(err, "Git error: "+err.Error(), w, rt.baseLogger)
			rt.db.BuildFail(p.Name)
			return
		}

		if err != nil {
			fo.Write([]byte("Git: " + err.Error() + "\n"))
		}

	} else if err != nil {
		helpers.SendBadRequestError(err, "Git error: "+err.Error(), w, rt.baseLogger)
		rt.db.BuildFail(p.Name)
		return
	}

	if p.Dockerfile != "" {

		cmd := "docker build -t " + p.ImageName + ":" + p.ImageTag + " -f " + p.Path + "/" + p.Dockerfile + " " + p.Path
		fo.Write([]byte("\n\n# " + cmd + "\n"))

		c := exec.Command("sh", "-c", cmd)
		c.Stdout = fo
		err = c.Run()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fo.Write([]byte("\n\nDocker Build error: " + err.Error()))

			fo.Seek(0, io.SeekStart)
			w.Header().Set("Content-Type", "text/plain")
			_, _ = io.Copy(w, fo)

			rt.db.BuildFail(p.Name)
			return
		}

	}

	cmd := "docker-compose -f " + p.Path + "/docker-compose.yml down"
	fo.Write([]byte("\n\n# " + cmd + "\n"))
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = fo
	err = c.Run()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fo.Write([]byte("\n\nDocker Compose error: " + err.Error()))
	}

	cmd = "docker-compose -f " + p.Path + "/docker-compose.yml up --build -d"
	fo.Write([]byte("\n\n# " + cmd + "\n"))

	c = exec.Command("sh", "-c", cmd)
	c.Stdout = fo
	err = c.Run()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fo.Write([]byte("\n\nDocker Compose error: " + err.Error()))

		fo.Seek(0, io.SeekStart)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.Copy(w, fo)

		rt.db.BuildFail(p.Name)
		return
	}

	// I have written the logs on the file, now I have to send them back to the user

	fo.Seek(0, io.SeekStart)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, _ = io.Copy(w, fo)

	rt.db.BuildSuccess(p.Name)
}

func BasicAuth() {
	panic("unimplemented")
}
