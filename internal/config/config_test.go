package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	os.Setenv("APP_PORT", "1111")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "pass")
	os.Setenv("POSTGRES_DB", "postgres")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_MIGRATIONS_PATH", ".Pg/Migrations/Path")

	cfg, err := FromEnv()
	require.NoError(t, err)

	assert.Equal(t, 1111, cfg.AppPort)
	assert.Equal(t, "localhost", cfg.PgConfig.PgHost)
	assert.Equal(t, "user", cfg.PgConfig.PgUser)
	assert.Equal(t, "pass", cfg.PgConfig.PgPassword)
	assert.Equal(t, "postgres", cfg.PgConfig.PgDB)
	assert.Equal(t, 5432, cfg.PgConfig.PgPort)
	assert.Equal(t, ".Pg/Migrations/Path", cfg.PgConfig.PgMigrationsPath)
}
