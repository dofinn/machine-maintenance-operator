package controller

import (
	"github.com/dofinn/machine-maintenance-operator/pkg/controller/machinemaintenance"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, machinemaintenance.Add)
}
