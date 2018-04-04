package config

var GlobCfg = Config{}

type Config struct {
	PORT             int64    `toml:"port"`
	ALLOW_ORIGIN     []string `toml:"allow_origin"`
	DB_FILE          string   `toml:"db_file"`
	JWT_KEY          string   `toml:"jwt_key"`
	JWT_EXPIRETIME   int64    `toml:"jwt_expiretime"`
	USE_STRICT_SLASH bool     `toml:"use_strict_slash"`
}
