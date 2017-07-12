package sctl

import (
	"fmt"

	"database/sql"
)

// Node contains node metadata
type Node struct {
	Project  string `json:"project"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
	IsMaster bool   `json:"isMaster"`
	User     string `json:"user"`
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

// SSHCommand Creates an ssh command to be executed on the node
func (node Node) SSHCommand(args ...string) Command {
	target := []string{node.Remote()}
	return Command{
		Main: "ssh",
		Args: append(target, args...),
	}
}

// RsyncCommand Creates a command to rsync minion executables to the node
func (node Node) RsyncCommand(execFolder string) Command {
	target := fmt.Sprintf("%s:/etc/init.d/", node.Remote())
	return Command{
		Main: "rsync",
		Args: []string{execFolder, target},
	}
}

// Remote Creates a remote node adress such as user@hostname
func (node Node) Remote() string {
	return fmt.Sprintf("%s@%s", node.User, node.IP)
}

// NodeSchema Returns the database schema for nodes
func NodeSchema() string {
	return `CREATE TABLE NODE(
      PROJECT VARCHAR(50),
      IP VARCHAR(50),
      OS VARCHAR(10) DEFAULT 'linux',
      IS_MASTER BOOLEAN,
      USER VARCHAR(40),
      FOREIGN KEY (PROJECT) REFERENCES PROJECT(NAME),
      PRIMARY KEY (PROJECT, IP)
    )`
}
