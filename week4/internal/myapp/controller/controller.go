package controller

var ControllerManager = AllControllers{}

type AllControllers struct {
	HelloController HelloController
}

func NewControllerManager(HelloController HelloController) AllControllers {
	return AllControllers{HelloController: HelloController}
}
