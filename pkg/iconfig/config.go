package iconfig

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// GetConfig loads the configuration and sets up watchers for the YAML and .env files if they exist.
// It calls the callback function whenever the watched files change.
func GetConfig(cfg interface{}, yamlPath, envPath string) error {
	if err := checkData(yamlPath, envPath); err != nil {
		return err
	}
	if err := ValidateYAML(yamlPath); err != nil {
		return err
	}

	// Initial configuration load
	if err := loadConfig(cfg, yamlPath, envPath); err != nil {
		return err
	}

	return nil
}
func WatchConfig(ctx context.Context, yamlPath string, callback func(string, error)) error {
	if yamlPath == "" {
		return fmt.Errorf("yamlPath is empty")
	}
	if err := ValidateYAML(yamlPath); err != nil {
		return err
	}
	if err := checkData(yamlPath, ""); err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	absPath, err := filepath.Abs(yamlPath)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	err = watcher.Add(filepath.Dir(absPath))
	if err != nil {
		watcher.Close()
		return fmt.Errorf("error adding file to watch: %w", err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					callback("", fmt.Errorf("watcher event channel closed"))
				} else if event.Name == absPath && (event.Op&fsnotify.Write == fsnotify.Write) {
					if err := ValidateYAML(yamlPath); err != nil {
						callback("", err)
					} else {
						callback("", nil)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					callback("", fmt.Errorf("watcher error channel closed"))
				} else {
					callback("", fmt.Errorf("watcher error: %w", err))
				}
			case <-ctx.Done():
				callback("service completed via context", nil)
				return
			}
		}
	}()

	return nil
}

func ValidateYAML(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("invalid YAML: %v", err)
	}

	return nil
}

// Check if files exist when paths are not empty
func checkData(yamlPath, envPath string) error {
	if yamlPath != "" && !fileExists(yamlPath) {
		return fmt.Errorf("yaml config file does not exist: %s", yamlPath)
	}
	if envPath != "" && !fileExists(envPath) {
		return fmt.Errorf(".env file does not exist: %s", envPath)
	}
	return nil
}

// loadConfig loads the configuration from .env and YAML files
func loadConfig(cfg interface{}, configPath, envPath string) error {
	// Loading .env file if path provided
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Loading and validating configuration from YAML file if path provided
	if configPath != "" {
		// Attempting to read a file with retries in case of error
		var fileContent []byte
		var err error
		for attempts := 0; attempts < 3; attempts++ {
			fileContent, err = os.ReadFile(configPath)
			if err == nil {
				break
			}
			if os.IsNotExist(err) {
				return fmt.Errorf("config file does not exist: %s", configPath)
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			return fmt.Errorf("failed to read config file after multiple attempts: %w", err)
		}

		if len(fileContent) == 0 {
			return fmt.Errorf("config file is empty: %s", configPath)
		}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			return fmt.Errorf("error parsing config file: %w", err)
		}
	}

	// Additional environment variable loading using cleanenv for validation
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return fmt.Errorf("error reading environment variables: %w", err)
	}

	return nil
}

// fileExists checks the presence of the file at path
func fileExists(path string) bool {
	p, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
