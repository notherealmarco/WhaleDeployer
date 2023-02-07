package database

import (
	"strings"
	"time"

	"github.com/notherealmarco/WhaleDeployer/service/structures"
)

func (db *appdbimpl) GetProjects() (*[]structures.Project, error) {

	projects := make([]structures.Project, 0)

	rows, err := db.c.Query("SELECT * FROM projects")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p structures.Project
		err = rows.Scan(&p.Name, &p.Path, &p.Description, &p.GitURL, &p.GitBranch, &p.Dockerfile, &p.ImageName, &p.ImageTag, &p.DeployKey, &p.LastBuild, &p.Status)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return &projects, err
}

func (db *appdbimpl) GetProject(name string) (*structures.Project, error) {
	var p structures.Project
	err := db.c.QueryRow("SELECT * FROM projects WHERE name=?", name).Scan(&p.Name, &p.Path, &p.Description, &p.GitURL, &p.GitBranch, &p.Dockerfile, &p.ImageName, &p.ImageTag, &p.DeployKey, &p.LastBuild, &p.Status)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *appdbimpl) AddProject(link *structures.Project) error {
	_, err := db.c.Exec("INSERT INTO projects (name, path, description, git_url, git_branch, dockerfile, image_name, image_tag, deploy_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", link.Name, link.Path, link.Description, link.GitURL, link.GitBranch, link.Dockerfile, link.ImageName, link.ImageTag, link.DeployKey)

	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: projects.name") {
		_, err = db.c.Exec("UPDATE projects SET path=?, description=?, git_url=?, git_branch=?, dockerfile=?, image_name=?, image_tag=?, deploy_key=? WHERE name=?", link.Path, link.Description, link.GitURL, link.GitBranch, link.Dockerfile, link.ImageName, link.ImageTag, link.DeployKey, link.Name)
	}
	return err
}

func (db *appdbimpl) BuildProject(name string) error {
	_, err := db.c.Exec("UPDATE projects SET status=? WHERE name=?", "running", name)
	return err
}

func (db *appdbimpl) BuildSuccess(name string) error {
	_, err := db.c.Exec("UPDATE projects SET last_build=?, status=? WHERE name=?", time.Now().String(), "success", name)
	return err
}

func (db *appdbimpl) BuildFail(name string) error {
	_, err := db.c.Exec("UPDATE projects SET last_build=?, status=? WHERE name=?", time.Now().String(), "fail", name)
	return err
}

func (db *appdbimpl) DeleteProject(name string) error {
	_, err := db.c.Exec("DELETE FROM projects WHERE name=?", name)
	return err
}
