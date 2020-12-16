package ssh

import (
	"fmt"
	"log"

	"github.com/bernarpa/photo/config"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

// Connect establishes a new SSH connection.
func Connect(target *config.Target) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: target.SSHUser,
		Auth: []ssh.AuthMethod{ssh.Password(target.SSHPassword)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	host := target.SSHHost + ":" + target.SSHPort
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	return client, session, nil
}

// Exec executes a command on an SSH server or exits the program in case of failure.
func Exec(client *ssh.Client, cmd string) []byte {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("SSH connection error: " + err.Error())
	}
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatal(fmt.Sprintf("SSH command execution error: %s\nCommand was %s", err.Error(), cmd))
	}
	return out
}

// Copy copies a file on the SSH server or exits the program in case of failure.
func Copy(client *ssh.Client, localFile string, remoteFile string) {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("SSH connection error: " + err.Error())
	}
	defer session.Close()
	err = scp.CopyPath(localFile, remoteFile, session)
	if err != nil {
		log.Fatal(fmt.Sprintf("SSH error copying %s to %s: %s", localFile, remoteFile, err.Error()))
	}
	session.Close()
}
