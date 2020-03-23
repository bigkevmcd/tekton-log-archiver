package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, archiver.LogArchiver) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, a archiver.LogArchiver) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, a); err != nil {
			return err
		}
	}
	return nil
}
