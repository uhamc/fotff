package utils

import (
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

func newSSHClient(addr string, user string, passwd string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	return ssh.Dial("tcp", addr, config)
}

func RunCmdViaSSH(addr string, user string, passwd string, cmd string) (string, error) {
	client, err := newSSHClient(addr, user, passwd)
	if err != nil {
		logrus.Errorf("new SSH client to %s err: %v", addr, err)
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	out, err := session.CombinedOutput(cmd)
	return string(out), err
}

func DownloadFileViaSSH(addr string, user string, passwd string, remoteFile string, localFile string) error {
	c, err := newSSHClient(addr, user, passwd)
	if err != nil {
		logrus.Errorf("new SSH client to %s err: %v", addr, err)
		return err
	}
	client, err := sftp.NewClient(c)
	if err != nil {
		logrus.Errorf("new SFTP client to %s err: %v", addr, err)
		return err
	}
	rf, err := client.Open(remoteFile)
	if err != nil {
		logrus.Errorf("open remote file %s at %s err: %v", remoteFile, addr, err)
		return err
	}
	lf, err := os.Create(localFile)
	if err != nil {
		logrus.Errorf("open local file %s at %s err: %v", remoteFile, addr, err)
		return err
	}
	logrus.Infof("copying %s at %s to %s...", remoteFile, addr, localFile)
	t1 := time.Now()
	n, err := io.CopyBuffer(lf, rf, make([]byte, 32*1024*1024))
	if err != nil {
		logrus.Errorf("copy %s at %s to %s err: %v", remoteFile, addr, localFile, err)
		return err
	}
	t2 := time.Now()
	cost := t2.Sub(t1).Seconds()
	logrus.Infof("copy %s at %s to %s done, size: %d cost: %.2fs speed: %.2fMB/s", remoteFile, addr, localFile, n, cost, float64(n)/cost/1024/1024)
	return nil
}
