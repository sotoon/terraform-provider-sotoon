package common

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TFProviderConfig struct {
	BepaUrls    []string
	BepaTimeout time.Duration
}

var lock = &sync.Mutex{}

var singleInstance *TFProviderConfig

func Config(ctx context.Context) *TFProviderConfig {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = loadConfig(ctx)
		}
	}

	return singleInstance
}

func loadConfig(ctx context.Context) *TFProviderConfig {
	config := TFProviderConfig{}

	config.BepaUrls = []string{
		"https://ha.bepa.sotoon.ir",
		"https://afra.bepa.sotoon.ir",
		"https://neda.bepa.sotoon.ir",
		"https://bepa.sotoon.ir",
	}

	config.BepaTimeout = 10 * time.Second

	config.logConfig(ctx)

	return &config
}

func (conf *TFProviderConfig) logConfig(ctx context.Context) {
	tflog.Info(ctx, "****** TERRAFORM PROVIDER CONFIG ******")
	tflog.Info(ctx, fmt.Sprintf("bepa urls:%s", conf.BepaUrls))
	tflog.Info(ctx, fmt.Sprintf("bepa timeout:%s", conf.BepaTimeout))
	tflog.Info(ctx, "****** TERRAFORM PROVIDER CONFIG ******")
}
