package controller

import (
	"mongodb-operator/pkg/controller/mongodb"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mongodb.Add)
}
