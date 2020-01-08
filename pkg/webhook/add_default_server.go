package webhook

import "github.com/Jaywoods/certs-controller/pkg/webhook/server"

func init() {
	// AddToManagerFuncs is a list of functions to create webhook servers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, server.Add)
}