package helpers

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/notherealmarco/WhaleDeployer/service/database"
	"github.com/sirupsen/logrus"
)

func Exec(cmd string, fo *os.File) error {

	fo.Write([]byte("\n\n# " + cmd + "\n"))
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = fo
	return c.Run()
}

func ExecOrWriteError(cmd string, fo *os.File, l logrus.FieldLogger) bool {
	err := Exec(cmd, fo)

	if err != nil {
		fo.Write([]byte("\nError: " + err.Error()))
		l.WithError(err).Error("Error executing command: " + cmd)
		return false
	}
	return true
}

func WriteResponse(status int, w *http.ResponseWriter, fo *os.File) error {
	(*w).WriteHeader(status)
	fo.Seek(0, io.SeekStart)
	(*w).Header().Set("Content-Type", "text/plain")
	_, err := io.Copy(*w, fo)
	return err
}

func ExecOrFail(cmd string, fo *os.File, l logrus.FieldLogger, w *http.ResponseWriter, db *database.AppDatabase, project string) bool {
	if !ExecOrWriteError(cmd, fo, l) {
		WriteResponse(http.StatusInternalServerError, w, fo)
		(*db).BuildFail(project)
		return false
	}
	return true
}
