package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Success(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("CACHE_DIR", "tmp")
	os.Setenv("CACHE_SIZE", "1024")

	cfg, err := New()

	require.NoError(t, err)
	require.Equal(t, cfg, &Config{
		Cache: CacheConfig{
			Size: 1024,
			Dir:  "tmp",
		},
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	})
}

func TestConfig_ParsePortErr(t *testing.T) {
	os.Setenv("PORT", "abc")
	os.Setenv("HOST", "localhost")
	os.Setenv("CACHE_DIR", "tmp")
	os.Setenv("CACHE_SIZE", "1024")

	cfg, err := New()
	require.Nil(t, cfg)
	require.Error(t, err)
}

func TestConfig_ParseCacheSizeErr(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("CACHE_DIR", "tmp")
	os.Setenv("CACHE_SIZE", "abc")

	cfg, err := New()
	require.Nil(t, cfg)
	require.Error(t, err)
}

func TestConfig_NotSetEnvs(t *testing.T) {
	cfg, err := New()
	require.Nil(t, cfg)
	require.Error(t, err)
}
