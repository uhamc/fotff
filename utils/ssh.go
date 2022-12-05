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

type Direct string

const (
	Download Direct = "download"
	Upload   Direct = "upload"
)

func TransFileViaSSH(verb Direct, addr string, user string, passwd string, remoteFile string, localFile string) error {
	c, err := newSSHClient(addr, user, passwd)
	if err != nil {
		logrus.Errorf("new SSH client to %s err: %v", addr, err)
		return err
	}
	defer c.Close()
	client, err := sftp.NewClient(c)
	if err != nil {
		logrus.Errorf("new SFTP client to %s err: %v", addr, err)
		return err
	}
	defer client.Close()
	var prep string
	var src, dst io.ReadWriteCloser
	if verb == Download {
		src, _ = client.Open(remoteFile)
		os.Remove(localFile)
		dst, _ = os.Create(localFile)
		prep = "to"
	} else {
		src, _ = os.Open(localFile)
		client.Remove(remoteFile)
		dst, _ = client.Create(remoteFile)
		prep = "from"
	}
	defer src.Close()
	defer dst.Close()
	logrus.Infof("%sing %s at %s %s %s...", verb, remoteFile, addr, prep, localFile)
	t1 := time.Now()
	n, err := io.CopyBuffer(dst, src, make([]byte, 32*1024*1024))
	if err != nil {
		logrus.Errorf("%s %s at %s %s %s err: %v", verb, remoteFile, addr, prep, localFile, err)
		return err
	}
	t2 := time.Now()
	cost := t2.Sub(t1).Seconds()
	logrus.Infof("%s %s at %s %s %s done, size: %d cost: %.2fs speed: %.2fMB/s", verb, remoteFile, addr, prep, localFile, n, cost, float64(n)/cost/1024/1024)
	return nil
}
