package controller

var controllerManager *AllControllers

type AllControllers struct {
	HelloController HelloController
}

func NewControllerManager(HelloController HelloController) AllControllers {
	return AllControllers{HelloController: HelloController}
}

func SetControllerManagerOnce(allControllers *AllControllers) {
	if controllerManager != nil {
		return
	}
	controllerManager = allControllers
	return
}

func GetControllerManager() *AllControllers {
	return controllerManager
}
