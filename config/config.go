package config

var GlobCfg = Config{}

type Config struct {
	PORT             int64    `toml:"port"`
	ALLOW_ORIGIN     []string `toml:"allow_origin"`
	DB_FILE          string   `toml:"db_file"`
	PLUGIN_DIR       string   `toml:"plugin_dir"`
	USE_STRICT_SLASH bool     `toml:"use_strict_slash"`
}
