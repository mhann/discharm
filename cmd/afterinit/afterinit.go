package afterinit

import (
	"log"
)

type afterInitializationFunction func()

var (
	afterInitializationFunctions []afterInitializationFunction
)

func RegisterAfterInitializationFunction(function afterInitializationFunction) {
	log.Println("Adding a function to be called after initialization of config and classes have been instanciated")
	afterInitializationFunctions = append(afterInitializationFunctions, function)
}

func RunAfterInitializationFunctions() {
	for _, function := range afterInitializationFunctions {
		function()
	}
}
