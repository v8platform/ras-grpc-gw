//go:build ignore

package main

import (
	"github.com/v8platform/ras-grpc-gw/swagger/patcher"
	"os"
)

var File string

func main() {

	_, err := os.Stat(File)
	if err != nil {
		panic(err)
	}

	json, err := os.ReadFile(File)
	if err != nil {
		panic(err)
	}

	swagger, err := patcher.PatchSwagger(json, "applications", "users", "access")
	if err != nil {
		return
	}

	err = os.WriteFile(File, swagger, 0777)
	if err != nil {
		panic(err)
	}

}
