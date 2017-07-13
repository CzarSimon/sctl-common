package sctl

import (
	"fmt"
	"os"

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

// RsyncFolderCMD Creates a command to rsync minion executables to the node
func (node Node) RsyncFolderCMD(execFolder, targetFolder string) Command {
	target := node.rsyncTarget(targetFolder)
	source := execFolder + fmt.Sprintf("%c", os.PathSeparator)
	return Command{
		Main: "rsync",
		Args: []string{"-a", source, target},
	}
}

// RsyncFileCMD Creates a command to rsync a file to the node
func (node Node) RsyncFileCMD(filePath, targetFolder string) Command {
	return Command{
		Main: "rsync",
		Args: []string{filePath, node.rsyncTarget(targetFolder)},
	}
}

// rsyncTarget Prepends user@hostname: to target folder if node is not local
func (node Node) rsyncTarget(targetFolder string) string {
	if node.IsLocal() {
		return targetFolder
	}
	return node.Remote() + ":" + targetFolder
}

// IsLocal Checks if node is localhost
func (node Node) IsLocal() bool {
	localhosts := []string{"localhost", "127.0.0.1", "0.0.0.0"}
	for _, localhost := range localhosts {
		if node.IP == localhost {
			return true
		}
	}
	return false
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
