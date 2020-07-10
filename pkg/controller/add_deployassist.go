package controller

import (
	"github.com/tthebst/deploy-assist-operator/pkg/controller/deployassist"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, deployassist.Add)
}
