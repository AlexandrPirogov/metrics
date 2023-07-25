package function

import "log"

// Helper function-values for applying
type onSuccessFunction func()

func CompareStringsDo(check string, defaulValue string, actionOnTrue onSuccessFunction) {
	if check != defaulValue {
		actionOnTrue()
	}
}

func CompareStringsDoOthewise(check string, defaulValue string, actionOnTrue onSuccessFunction, actionOnFalse onSuccessFunction) {
	if check != defaulValue {
		actionOnTrue()
		return
	}

	actionOnFalse()
}

func CompareBoolssDo(check bool, defaulValue bool, actionOnTrue onSuccessFunction) {
	if check != defaulValue {
		actionOnTrue()
	}
}

func CompareIntsDo(check int, with int, actionOnTrue onSuccessFunction) {
	if check != with {
		actionOnTrue()
	}
}

func ErrFatalCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
