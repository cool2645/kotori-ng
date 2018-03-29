package config

var GlobCfg = Config{}

type Config struct {
	PORT         int64    `toml:"port"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
	DB_FILE      string   `toml:"db_file"`
}
