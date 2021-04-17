package internal

import (
	"fmt"
	"runtime"
)

const notProvided = "[not provided]"
const ApplicationName = "gbuild"

// Default build-time variables
// Values are overriden via ldflags
var (
	version   = "development"
	buildDate = notProvided
	platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	goVersion = runtime.Version()
)

func PrintVersion() {
	fmt.Println("Application:    ", ApplicationName)
	fmt.Println("Version:        ", version)
	fmt.Println("BuildDate:      ", buildDate)
	fmt.Println("Platform:       ", platform)
	fmt.Println("GoVersion:      ", goVersion)
}
