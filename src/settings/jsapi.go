package settings

func JsHandle(args ...string) interface{} {
	switch args[0] {
	case "ML_enabled":
		return Enable_ML
	case "enable_ML":
		Enable_ML = true

	case "disable_ML":
		Enable_ML = false

	case "app_port":
		return App_port

	case ""
	}
	return nil
}
