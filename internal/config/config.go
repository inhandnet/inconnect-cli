package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentContext string              `yaml:"current-context"`
	Contexts       map[string]*Context `yaml:"contexts"`
}

func DefaultPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "inconnect", "config.yaml")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Contexts: make(map[string]*Context)}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]*Context)
	}
	return &cfg, nil
}

func Save(cfg *Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Config) ActiveContext() (*Context, error) {
	name := c.ActiveContextName()
	if name == "" {
		return nil, fmt.Errorf("no active context; run 'inconnect auth login' first")
	}
	ctx, ok := c.Contexts[name]
	if !ok {
		return nil, fmt.Errorf("context %q not found", name)
	}
	if h := os.Getenv("INCONNECT_HOST"); h != "" {
		ctx.Host = h
	}
	return ctx, nil
}

func (c *Config) ActiveContextName() string {
	if name := os.Getenv("INCONNECT_CONTEXT"); name != "" {
		return name
	}
	return c.CurrentContext
}

func (c *Config) SetContext(name string, ctx *Context) {
	c.Contexts[name] = ctx
}

func (c *Config) DeleteContext(name string) {
	delete(c.Contexts, name)
	if c.CurrentContext == name {
		c.CurrentContext = ""
	}
}
