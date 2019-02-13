package compute

import (
	"errors"
	"time"

	"github.com/agoca80/bosh-abiquo-cpi/helpers"
)

const maxTries = 6
const pollWait = time.Second * 10

func retry(l *helpers.Logger, fn func() (bool, error)) (err error) {
	for try := 0; try < maxTries; try++ {
		done, err := fn()
		if err != nil {
			l.Debug("error while trying : %v", err)
		}

		if done {
			return nil
		}

		time.Sleep(pollWait)
	}

	return errors.New("too many tries")
}

func (v *vm) synchronize() (err error) {
	lastSynchronize := v.LastSynchronize
	v.Debug(v.Name+" last synchronize was at %v", lastSynchronize)
	err = retry(v.Logger, func() (synched bool, err error) {
		v.Debug(v.Name + " synchronizing vm")
		err = v.Read(v.VirtualMachine)
		if err != nil {
			return
		}
		synched = lastSynchronize < v.LastSynchronize
		return
	})
	if err != nil {
		return
	}

	v.Debug(v.Name+" vm synchronized at %v", v.LastSynchronize)
	return
}
