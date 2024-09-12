package iconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDir   string
	testMutex sync.Mutex
)

func TestMain(m *testing.M) {
	// Setup performed once before all tests
	testDir, err := os.MkdirTemp("", "config-tests-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup after all tests
	err = os.RemoveAll(testDir)
	if err != nil {
		fmt.Printf("Failed to remove temp directory: %v\n", err)
		os.Exit(1)
	}

	// Exit with the code returned from m.Run()
	os.Exit(exitCode)
}

func setupTestValid(t *testing.T) (string, string, func()) {
	testMutex.Lock()
	t.Helper()

	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	yamlPath := filepath.Join(testDir, fmt.Sprintf("config.test.%s.yaml", salt))
	envPath := filepath.Join(testDir, fmt.Sprintf(".env.test.%s", salt))

	envContent := []byte("SERVER_PORT=9090\n")
	err := os.WriteFile(envPath, envContent, 0o644)
	require.NoError(t, err, "Failed to create env file")

	yamlContent := []byte("server:\n  port: \"8081\"\n")
	err = os.WriteFile(yamlPath, yamlContent, 0o644)
	require.NoError(t, err, "Failed to create yaml file")

	return yamlPath, envPath, func() {
		os.Clearenv()
		err := os.Unsetenv("SERVER_PORT")
		require.NoError(t, err, "Failed to unset env")

		removeFile := func(path string) error {
			var err error
			for i := 0; i < 5; i++ {
				time.Sleep(10 * time.Millisecond)

				err = os.Remove(path)
				if err == nil {
					return nil
				}
				if !os.IsExist(err) {
					return err
				}
			}
			return err
		}

		err = removeFile(yamlPath)
		require.NoError(t, err, "Failed to empty env file")

		err = removeFile(envPath)
		require.NoError(t, err, "Failed to empty yaml file")
		testMutex.Unlock()
	}
}

type AppConfig struct {
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	} `yaml:"server"`
}

func TestGetConfig(t *testing.T) {
	t.Run("File change detection without env", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		callbackCalled := make(chan error, 1)

		err := os.Unsetenv("SERVER_PORT")
		require.NoError(t, err, "Failed to unset env")
		err = os.WriteFile(envPath, []byte{}, 0o644)
		assert.NoError(t, err)
		err = GetConfig(&cfg, yamlPath, envPath)
		assert.NoError(t, err)

		err = WatchConfig(ctx, yamlPath, func(msg string, err error) {
			assert.NoError(t, err)
			var newConfig AppConfig
			err = GetConfig(&newConfig, yamlPath, envPath)
			assert.NoError(t, err)
			cfg = newConfig
			callbackCalled <- err
		})
		assert.NoError(t, err)

		assert.Equal(t, "8081", cfg.Server.Port)

		// Modify the YAML file
		newYamlContent := []byte("server:\n  port: \"9000\"\n")
		err = os.WriteFile(yamlPath, newYamlContent, 0o644)
		assert.NoError(t, err)

		// Wait for the callback to be called
		select {
		case err := <-callbackCalled:
			cancel()
			assert.NoError(t, err)
			assert.Equal(t, "9000", cfg.Server.Port, "Updated YAML value should be used")
		case <-time.After(10 * time.Second):
			t.Fatal("Callback was not called within 10 seconds after YAML change")
		}
	})

	t.Run("File change detection", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		callbackCalled := make(chan error, 1)

		err := GetConfig(&cfg, yamlPath, envPath)
		assert.NoError(t, err)

		err = WatchConfig(ctx, yamlPath, func(msg string, err error) {
			assert.NoError(t, err)
			var newConfig AppConfig
			err = GetConfig(&newConfig, yamlPath, envPath)
			assert.NoError(t, err)
			cfg = newConfig
			callbackCalled <- err
		})
		assert.NoError(t, err)

		assert.Equal(t, "9090", cfg.Server.Port)

		// Modify the YAML file
		newYamlContent := []byte("server:\n  port: \"9000\"\n")
		err = os.WriteFile(yamlPath, newYamlContent, 0o644)
		assert.NoError(t, err)

		// Wait for the callback to be called
		select {
		case err := <-callbackCalled:
			cancel()
			assert.NoError(t, err)
			assert.Equal(t, "9090", cfg.Server.Port, "Updated YAML value should be used")
		case <-time.After(10 * time.Second):
			t.Fatal("Callback was not called within 10 seconds after YAML change")
		}
	})

	t.Run("Load from .env and YAML", func(t *testing.T) {
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		err := GetConfig(&cfg, yamlPath, envPath)
		assert.NoError(t, err)
		assert.Equal(t, "9090", cfg.Server.Port) // .env takes precedence over YAML
	})

	t.Run("Missing .env file", func(t *testing.T) {
		yamlPath, _, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		err := GetConfig(&cfg, yamlPath, "nonexistent.env")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("Missing YAML file", func(t *testing.T) {
		_, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		err := GetConfig(&cfg, "nonexistent.yaml", envPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		yamlContent := []byte("invalid: yaml: content")
		err := os.WriteFile(yamlPath, yamlContent, 0o644)
		assert.NoError(t, err)
		var cfg AppConfig
		err = GetConfig(&cfg, yamlPath, envPath)
		assert.Error(t, err)
	})

	t.Run("Invalid ENV", func(t *testing.T) {
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		err := os.WriteFile(envPath, []byte("SERVER_PORT-===9090\n"), 0o644)
		assert.NoError(t, err)
		var cfg AppConfig
		err = GetConfig(&cfg, yamlPath, envPath)
		assert.Equal(t, "", cfg.Server.Port)
	})

	t.Run("Validation of env variables", func(t *testing.T) {
		_, _, cleanup := setupTestValid(t)
		defer cleanup()

		_ = os.Setenv("SERVER_PORT", "7070")
		var cfg AppConfig
		err := GetConfig(&cfg, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "7070", cfg.Server.Port)
	})

	t.Run("Default values", func(t *testing.T) {
		_, _, cleanup := setupTestValid(t)
		defer cleanup()

		err := os.Unsetenv("SERVER_PORT")
		require.NoError(t, err, "Failed to unset env")
		var cfg AppConfig
		err = GetConfig(&cfg, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "8080", cfg.Server.Port) // Check default value
	})

	t.Run("Concurrent access", func(t *testing.T) {
		yamlPath, envPath, cleanup := setupTestValid(t)
		defer cleanup()

		var cfg AppConfig
		err := GetConfig(&cfg, yamlPath, envPath)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Just read cfg to check for absence of data races
				_ = cfg.Server.Port
			}()
		}
		wg.Wait()
	})
}
