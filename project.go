package sctl

// Project holds metadata about a project
type Project struct {
	Name        string `json:"name"`
	Folder      string `json:"folder"`
	SwarmToken  string `json:"swarmToken"`
	IsActive    bool   `json:"isActive"`
	Network     string `json:"network"`
	MasterToken string `json:"masterToken"`
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
