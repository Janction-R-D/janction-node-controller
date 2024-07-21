package utils

import "runtime"

func MaxParallelism() int {
	maxProc := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	if maxProc == 0 {
		// in some case it could be 0
		// see https://stackoverflow.com/questions/13234749/golang-how-to-verify-number-of-processors-on-which-a-go-program-is-running#comment32938656_13245047
		return numCPU
	}

	if maxProc < numCPU {
		return maxProc
	}
	return numCPU
}
