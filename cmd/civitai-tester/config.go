package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	ServerPort  int    `json:"server_port"`
	TestTimeout int    `json:"test_timeout_seconds"`
	TestLimits  struct {
		ModelsLimit   int `json:"models_limit"`
		ImagesLimit   int `json:"images_limit"`
		CreatorsLimit int `json:"creators_limit"`
		TagsLimit     int `json:"tags_limit"`
	} `json:"test_limits"`
	CustomTests struct {
		SkipRateLimit bool     `json:"skip_rate_limit_test"`
		OnlyTests     []string `json:"only_tests,omitempty"`
		SkipTests     []string `json:"skip_tests,omitempty"`
	} `json:"custom_tests"`
}

func loadConfig() *Config {
	config := &Config{
		ServerPort:  9999,
		TestTimeout: 30,
	}
	
	// Set default test limits
	config.TestLimits.ModelsLimit = 5
	config.TestLimits.ImagesLimit = 5
	config.TestLimits.CreatorsLimit = 5
	config.TestLimits.TagsLimit = 10
	
	// Load from environment variables
	if apiKey := os.Getenv("CIVITAI_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}
	
	if port := os.Getenv("TESTER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.ServerPort = p
		}
	}
	
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.TestTimeout = t
		}
	}
	
	// Try to load from config file
	if data, err := os.ReadFile("config.json"); err == nil {
		var fileConfig Config
		if err := json.Unmarshal(data, &fileConfig); err == nil {
			// Merge file config with defaults
			if fileConfig.APIKey != "" {
				config.APIKey = fileConfig.APIKey
			}
			if fileConfig.ServerPort != 0 {
				config.ServerPort = fileConfig.ServerPort
			}
			if fileConfig.TestTimeout != 0 {
				config.TestTimeout = fileConfig.TestTimeout
			}
			if fileConfig.TestLimits.ModelsLimit > 0 {
				config.TestLimits.ModelsLimit = fileConfig.TestLimits.ModelsLimit
			}
			if fileConfig.TestLimits.ImagesLimit > 0 {
				config.TestLimits.ImagesLimit = fileConfig.TestLimits.ImagesLimit
			}
			if fileConfig.TestLimits.CreatorsLimit > 0 {
				config.TestLimits.CreatorsLimit = fileConfig.TestLimits.CreatorsLimit
			}
			if fileConfig.TestLimits.TagsLimit > 0 {
				config.TestLimits.TagsLimit = fileConfig.TestLimits.TagsLimit
			}
			config.CustomTests = fileConfig.CustomTests
		}
	}
	
	return config
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}

func (c *Config) Print() {
	fmt.Println("Configuration:")
	fmt.Printf("  Server Port: %d\n", c.ServerPort)
	fmt.Printf("  Test Timeout: %d seconds\n", c.TestTimeout)
	fmt.Printf("  API Key: %s\n", func() string {
		if c.APIKey == "" {
			return "Not set"
		}
		return "***" + c.APIKey[len(c.APIKey)-4:]
	}())
	fmt.Printf("  Test Limits: Models=%d, Images=%d, Creators=%d, Tags=%d\n",
		c.TestLimits.ModelsLimit,
		c.TestLimits.ImagesLimit,
		c.TestLimits.CreatorsLimit,
		c.TestLimits.TagsLimit)
	if len(c.CustomTests.OnlyTests) > 0 {
		fmt.Printf("  Only Tests: %v\n", c.CustomTests.OnlyTests)
	}
	if len(c.CustomTests.SkipTests) > 0 {
		fmt.Printf("  Skip Tests: %v\n", c.CustomTests.SkipTests)
	}
}