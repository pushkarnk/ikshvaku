package controller

import (
	"github.com/ikshvaku/ikshvaku/pkg/controller/ikshvaku"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ikshvaku.Add)
}
