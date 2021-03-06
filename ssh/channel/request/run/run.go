package run

import (
	"github.com/janoszen/containerssh/backend"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"sync"
)

type responseMsg struct {
	exitStatus uint32
}

func run(program string, channel ssh.Channel, session backend.Session) error {
	var mutex = &sync.Mutex{}
	closeSession := func() {
		mutex.Lock()
		session.Close()
		exitCode := session.GetExitCode()
		mutex.Unlock()

		if exitCode < 0 {
			log.Printf("invalid exit code (%d)", exitCode)
		}

		//Send the exit status before closing the session. No reply is sent.
		_, _ = channel.SendRequest("exit-status", false, ssh.Marshal(responseMsg{
			exitStatus: uint32(exitCode),
		}))
		//Close the channel as described by the RFC
		_ = channel.Close()
	}
	err := session.RequestProgram(program, channel, channel, channel.Stderr(), closeSession)
	if err != nil {
		return err
	}
	return nil
}
