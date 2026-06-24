package functioncaller

import (
	"runtime"
)

func PrintFuncName() string {
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {

	}

	caller := runtime.FuncForPC(fpcs[0] - 1)

	if caller == nil {
		return ""
	}

	return caller.Name()
}
