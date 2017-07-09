package sctl

import "database/sql"

// Node contains node metadata
type Node struct {
	Project  string `json:"project"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
	IsMaster bool   `json:"isMaster"`
}

//SetToMinion sets a node struct to hold values of a minion
func (node *Node) SetToMinion(db *sql.DB) error {
	if node.Project == "" {
		projectName, err := GetActiveProject(db)
		if err != nil {
			return err
		}
		node.Project = projectName
	}
	node.IsMaster = false
	node.OS = "linux"
	return nil
}

// NodeSchema Returns the database schema for nodes
func NodeSchema() string {
	return `CREATE TABLE NODE(
      PROJECT VARCHAR(50),
      IP VARCHAR(50),
      OS VARCHAR(10) DEFAULT 'linux',
      IS_MASTER BOOLEAN,
      FOREIGN KEY (PROJECT) REFERENCES PROJECT(NAME),
      PRIMARY KEY (PROJECT, IP)
    )`
}
