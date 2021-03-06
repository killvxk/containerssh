package kuberun

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (session *kubeRunSession) GetExitCode() int32 {
	if session.exitCode < 0 && session.pod != nil {
		log.Tracef("Fetching exit code...")
		pod, err := session.client.
			CoreV1().
			Pods(session.pod.Namespace).
			Get(session.ctx, session.pod.Name, v1.GetOptions{})
		if err != nil {
			log.Infof("Error while fetching exit code (%s)", err)
			return session.exitCode
		}
		containerStatuses := pod.Status.ContainerStatuses
		if len(containerStatuses) > 0 {
			containerStatus := containerStatuses[session.config.Pod.ConsoleContainerNumber]
			if containerStatus.State.Terminated != nil {
				session.exitCode = containerStatus.State.Terminated.ExitCode
			} else if containerStatus.LastTerminationState.Terminated != nil {
				session.exitCode =
					containerStatuses[session.config.Pod.ConsoleContainerNumber].
						LastTerminationState.
						Terminated.
						ExitCode
			}
		}
	}
	return session.exitCode
}
