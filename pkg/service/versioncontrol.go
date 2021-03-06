package service

import (
	"github.com/abhinavk1/remote-config-server/pkg/versioncontrol"
	"log"
	"time"
)

type AbstractVersionControl interface {
	Checkout(branch string) error

	Pull() error

	PollRepository(interval time.Duration) error
}

type VersionControl struct {
	versionControlClient versioncontrol.Abstract
}

func NewVersionControl(versionControlClient versioncontrol.Abstract) *VersionControl {
	return &VersionControl{
		versionControlClient: versionControlClient,
	}
}

func (svc *VersionControl) Checkout(branch string) error {
	return svc.versionControlClient.Checkout(branch)
}

func (svc *VersionControl) Pull() error {
	return svc.versionControlClient.Pull()
}

func (svc *VersionControl) PollRepository(interval time.Duration) error {

	go func() {
		for true {
			err := svc.versionControlClient.Pull()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(interval)
		}
	}()

	return nil
}
