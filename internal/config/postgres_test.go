package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPgConfig_Dsn(t *testing.T) {
	pgConf := PgConfig{
		PgHost:           "PgHost",
		PgUser:           "PgUser",
		PgPassword:       "PgPassword",
		PgDB:             "PgDB",
		PgPort:           1111,
		PgMigrationsPath: ".Pg/Migrations/Path",
	}

	expectedDsn := "host=PgHost port=1111 user=PgUser password=PgPassword dbname=PgDB sslmode=disable"

	require.Equal(t, expectedDsn, pgConf.Dsn())
}
