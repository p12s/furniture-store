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
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}

	err := godotenv.Load(os.ExpandEnv(fmt.Sprintf("%s%s", goPath, THIS_PACKAGE_ENV_PATH)))
	assert.Equal(t, nil, err)

	_, err = New()
	assert.Equal(t, nil, err)
}
