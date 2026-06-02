package factory

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/inhandnet/inconnect-cli/internal/api"
	"github.com/inhandnet/inconnect-cli/internal/config"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
)

type Factory struct {
	IO         *iostreams.IOStreams
	ConfigPath string
	configOnce sync.Once
	config     *config.Config
	configErr  error
}

func New() *Factory {
	return &Factory{
		IO:         iostreams.System(),
		ConfigPath: config.DefaultPath(),
	}
}

func (f *Factory) Config() (*config.Config, error) {
	f.configOnce.Do(func() {
		f.config, f.configErr = config.Load(f.ConfigPath)
	})
	return f.config, f.configErr
}

func (f *Factory) ReloadConfig() {
	f.configOnce = sync.Once{}
}

func (f *Factory) SaveConfig() error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}
	return config.Save(cfg, f.ConfigPath)
}

func (f *Factory) APIClient() (*api.APIClient, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}

	ctx, err := cfg.ActiveContext()
	if err != nil {
		return nil, err
	}

	transport := &api.TokenTransport{
		Token:        ctx.EffectiveToken(),
		RefreshToken: ctx.RefreshToken,
		Host:         ctx.BaseURL(),
		ClientID:     ctx.ClientID,
		ClientSecret: ctx.ClientSecret,
		OnRefresh: func(accessToken, refreshToken string, expiry time.Time) {
			ctx.Token = accessToken
			ctx.RefreshToken = refreshToken
			ctx.ExpiresAt = expiry
			_ = f.SaveConfig()
		},
	}

	verbose := 100
	if v := os.Getenv("INCONNECT_VERBOSE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			verbose = n
		}
	}

	return api.NewAPIClient(ctx.BaseURL(), transport, verbose), nil
}
