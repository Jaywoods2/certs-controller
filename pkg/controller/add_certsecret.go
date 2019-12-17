package controller

import (
	"github.com/Jaywoods/certs-controller/pkg/controller/certsecret"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, certsecret.Add)
}
