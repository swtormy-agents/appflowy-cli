package appflowy

import "errors"

type Config struct {
	URL   string
	Token string
}

func LoadConfig(getenv func(string) string) (Config, error) {
	cfg := Config{
		URL:   getenv("APPFLOWY_URL"),
		Token: getenv("APPFLOWY_TOKEN"),
	}
	if cfg.URL != "" {
		n, _, err := ValidateBaseURL(cfg.URL)
		if err != nil {
			return Config{}, err
		}
		cfg.URL = n
	}
	if cfg.Token == "" && cfg.URL == "" {
		return cfg, nil
	}
	if cfg.URL == "" {
		return Config{}, errors.New("APPFLOWY_URL is not set")
	}
	return cfg, nil
}
