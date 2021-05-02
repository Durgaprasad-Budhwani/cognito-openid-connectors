package common

import (
	"os"
)

func GetStage() string {
	if os.Getenv(ServerlessOffline) == TrueString {
		return "/" + os.Getenv(Stage)
	}
	return ""
}
