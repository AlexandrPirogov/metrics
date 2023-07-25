package function

import "log"

// Helper function-values for applying
type assignFunction func()

func CompareStringsDo(check string, defaulValue string, actionOnTrue assignFunction) {
	if check != defaulValue {
		actionOnTrue()
	}
}

func CompareStringsDoOthewise(check string, defaulValue string, actionOnTrue assignFunction, actionOnFalse assignFunction) {
	if check != defaulValue {
		actionOnTrue()
		return
	}

	actionOnFalse()
}

func CompareBoolssDo(check bool, defaulValue bool, actionOnTrue assignFunction) {
	if check != defaulValue {
		actionOnTrue()
	}
}

func ErrFatalCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
