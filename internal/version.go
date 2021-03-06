package internal

import (
	"fmt"
	"runtime"
)

const notProvided = "[not provided]"
const ApplicationName = "gbuild"

// Default build-time variables
// Values are overriden via ldflags
// Example:
// export GIT_TAG=$(git describe --tags $(git rev-list --tags --max-count=1)) && \
//   go build -o gbuild \
//   -ldflags "-X 'github.com/chaordic-io/gbuild/internal.BuildDate=$(date)'\
// 			-X 'github.com/chaordic-io/gbuild/internal.Version=$GIT_TAG'"

var (
	Version   = "development"
	BuildDate = notProvided
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	GoVersion = runtime.Version()
)

func PrintVersionInfo() {
	fmt.Println("Application:    ", ApplicationName)
	fmt.Println("Version:        ", Version)
	fmt.Println("BuildDate:      ", BuildDate)
	fmt.Println("Platform:       ", Platform)
	fmt.Println("GoVersion:      ", GoVersion)
}
