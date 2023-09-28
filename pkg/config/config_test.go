package config

import "testing"

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./testdata/config.yml")
	if err != nil {
		t.Fatalf("Failed to load config, err: %v\n", err)
	}
	t.Logf("Loaded config: %v\n", config)
}
