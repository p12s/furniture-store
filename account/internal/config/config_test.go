package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/stretchr/testify/assert"
)

const DIR_ENV_PATH = ".env.example"

func TestNew(t *testing.T) {
	currentDir, err := os.Getwd()
	assert.Equal(t, nil, err)

	configPath := filepath.Dir(filepath.Dir(currentDir))
	err = godotenv.Load(os.ExpandEnv(fmt.Sprintf("%s/%s", configPath, DIR_ENV_PATH)))
	assert.Equal(t, nil, err)

	_, err = config.New()
	assert.Equal(t, nil, err)
}
