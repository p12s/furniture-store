package config

import (
	"fmt"
	"go/build"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const THIS_PACKAGE_ENV_PATH = "/src/github.com/p12s/furniture-store/account/.env.example"

func TestNew(t *testing.T) {
	fmt.Println("os.Getenv GOPATH:", os.Getenv("GOPATH"))
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	fmt.Println("goPath:", goPath)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("pwd:", pwd)

	err = godotenv.Load(os.ExpandEnv(fmt.Sprintf("%s%s", goPath, THIS_PACKAGE_ENV_PATH)))
	assert.Equal(t, nil, err)

	_, err = New()
	assert.Equal(t, nil, err)
}
