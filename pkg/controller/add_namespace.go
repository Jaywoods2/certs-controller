package controller

import (
	"github.com/Jaywoods/certs-controller/pkg/controller/namespace"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, namespace.Add)
}
