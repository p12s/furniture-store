package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const DIR_ENV_PATH = ".env.example"

func TestNew(t *testing.T) {
	// fmt.Println("os.Getenv GOPATH:", os.Getenv("GOPATH")) //
	// goPath := os.Getenv("GOPATH")
	// if goPath == "" {
	// 	goPath = build.Default.GOPATH
	// }
	// fmt.Println("goPath:", goPath) // /home/runner/go

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("pwd:", pwd)
	parentRoot := filepath.Dir(filepath.Dir(pwd))
	fmt.Println("parentRoot:", parentRoot)

	path := fmt.Sprintf("%s/%s", parentRoot, DIR_ENV_PATH)
	fmt.Println("path-path-path:", path)

	err = godotenv.Load(os.ExpandEnv(path))
	assert.Equal(t, nil, err)

	_, err = New()
	assert.Equal(t, nil, err)
}
