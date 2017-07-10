package sctl

import "database/sql"

// Project holds metadata about a project
type Project struct {
	Name        string `json:"name"`
	Folder      string `json:"folder"`
	SwarmToken  string `json:"swarmToken"`
	IsActive    bool   `json:"isActive"`
	Network     string `json:"network"`
	MasterToken string `json:"masterToken"`
	Master      string `json:"master"`
}

// NewProject generates a new project based on name an folder
func NewProject(name, folder string) Project {
	return Project{
		Name:        name,
		Folder:      folder,
		IsActive:    true,
		Network:     name + "-net",
		MasterToken: GenerateToken(4),
	}
}

// ProjectSchema Returns the database schema for projects
func ProjectSchema() string {
	return `CREATE TABLE PROJECT(
      NAME VARCHAR(50) PRIMARY KEY,
      FOLDER VARCHAR(255),
      SWARM_TOKEN VARCHAR(100),
      IS_ACTIVE BOOLEAN,
      NETWORK VARCHAR(60),
      MASTER_TOKEN VARCHAR(260)
    )`
}

// GetActiveProject Fetches the name of the current active project
func GetActiveProject(db *sql.DB) (string, error) {
	var projectName string
	query := "SELECT NAME FROM PROJECT WHERE IS_ACTIVE=1"
	err := db.QueryRow(query).Scan(&projectName)
	return projectName, err
}

//MakeMasterNode creates node struct for a project master node
func (project Project) MakeMasterNode() Node {
	return Node{
		Project:  project.Name,
		IP:       project.Master,
		OS:       "linux",
		IsMaster: true,
	}
}

// Insert Inserts a new project into the database and sets it to active
func (project Project) Insert(db *sql.DB) error {
	query := "INSERT INTO PROJECT(NAME, FOLDER, SWARM_TOKEN, IS_ACTIVE, NETWORK, MASTER_TOKEN) VALUES ($1,$2,$3,$4,$5,$6)"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Name, project.Folder, project.SwarmToken, true, project.Network, project.MasterToken)
	if err != nil {
		return err
	}
	return project.InactivateOthers(db)
}

// InactivateOthers Sets all other projects besides the supplied one to inactive
func (project Project) InactivateOthers(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE PROJECT SET IS_ACTIVE=0 WHERE NAME!=$1")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Name)
	if err != nil {
		return err
	}
	return nil
}
