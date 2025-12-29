package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_WithDefaults(t *testing.T) {
	// Clear environment variables to test defaults
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"APP_PORT", "APP_ENV", "JWT_SECRET", "JWT_EXPIRE_HOURS",
		"JWT_REFRESH_EXPIRE_HOURS", "UPLOAD_PATH", "MAX_FILE_SIZE",
	}
	
	// Store original values
	originalValues := make(map[string]string)
	for _, env := range envVars {
		originalValues[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	
	// Restore original values after test
	defer func() {
		for _, env := range envVars {
			if val, exists := originalValues[env]; exists && val != "" {
				os.Setenv(env, val)
			}
		}
	}()

	config := Load()

	// Test default values
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "3306", config.Database.Port)
	assert.Equal(t, "root", config.Database.User)
	assert.Equal(t, "", config.Database.Password)
	assert.Equal(t, "go_commerce", config.Database.Name)
	assert.Equal(t, "utf8mb4", config.Database.Charset)
	assert.Equal(t, true, config.Database.ParseTime)
	assert.Equal(t, "Local", config.Database.Loc)

	assert.Equal(t, "8080", config.App.Port)
	assert.Equal(t, "development", config.App.Env)

	assert.Equal(t, "your-secret-key", config.JWT.Secret)
	assert.Equal(t, 24, config.JWT.ExpireHours)
	assert.Equal(t, 168, config.JWT.RefreshExpireHours)

	assert.Equal(t, "./uploads", config.Upload.Path)
	assert.Equal(t, int64(5242880), config.Upload.MaxFileSize)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Set test environment variables
	testEnvVars := map[string]string{
		"DB_HOST":                   "testhost",
		"DB_PORT":                   "3307",
		"DB_USER":                   "testuser",
		"DB_PASSWORD":               "testpass",
		"DB_NAME":                   "test_db",
		"APP_PORT":                  "9000",
		"APP_ENV":                   "production",
		"JWT_SECRET":                "test-secret",
		"JWT_EXPIRE_HOURS":          "48",
		"JWT_REFRESH_EXPIRE_HOURS":  "336",
		"UPLOAD_PATH":               "./test-uploads",
		"MAX_FILE_SIZE":             "10485760",
	}

	// Store original values
	originalValues := make(map[string]string)
	for key, value := range testEnvVars {
		originalValues[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	// Restore original values after test
	defer func() {
		for key := range testEnvVars {
			if val, exists := originalValues[key]; exists && val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	config := Load()

	// Test environment variable values
	assert.Equal(t, "testhost", config.Database.Host)
	assert.Equal(t, "3307", config.Database.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "test_db", config.Database.Name)

	assert.Equal(t, "9000", config.App.Port)
	assert.Equal(t, "production", config.App.Env)

	assert.Equal(t, "test-secret", config.JWT.Secret)
	assert.Equal(t, 48, config.JWT.ExpireHours)
	assert.Equal(t, 336, config.JWT.RefreshExpireHours)

	assert.Equal(t, "./test-uploads", config.Upload.Path)
	assert.Equal(t, int64(10485760), config.Upload.MaxFileSize)
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "Empty environment variable",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original value
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					os.Setenv(tt.key, originalValue)
				} else {
					os.Unsetenv(tt.key)
				}
			}()

			// Set test environment variable
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoad_InvalidIntegerValues(t *testing.T) {
	// Test with invalid integer values
	testEnvVars := map[string]string{
		"JWT_EXPIRE_HOURS":         "invalid",
		"JWT_REFRESH_EXPIRE_HOURS": "also_invalid",
		"MAX_FILE_SIZE":            "not_a_number",
	}

	// Store original values
	originalValues := make(map[string]string)
	for key, value := range testEnvVars {
		originalValues[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	// Restore original values after test
	defer func() {
		for key := range testEnvVars {
			if val, exists := originalValues[key]; exists && val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	config := Load()

	// Should use default values when parsing fails
	assert.Equal(t, 0, config.JWT.ExpireHours)        // strconv.Atoi returns 0 on error
	assert.Equal(t, 0, config.JWT.RefreshExpireHours) // strconv.Atoi returns 0 on error
	assert.Equal(t, int64(0), config.Upload.MaxFileSize) // strconv.ParseInt returns 0 on error
}