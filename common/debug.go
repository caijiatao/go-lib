package common

import "os"

const (
	debugModeEnvKey = "DEBUG_MODE"
)

func IsDebugMode() bool {
	_, truth := os.LookupEnv(debugModeEnvKey)
	return truth
}

func SetDebugMode() {
	_ = os.Setenv(debugModeEnvKey, "true")
}
