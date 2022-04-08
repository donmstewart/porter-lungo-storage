package mongdb_lungo

import (
	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/storage/plugins"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type StoreTypeEnum string

const (
	PluginKey = plugins.PluginInterface + ".porter.mongodb_lungo"

	StoreTypeEnumMemory StoreTypeEnum = "MEMORY"
	StoreTypeEnumFile   StoreTypeEnum = "FILE"
)

// PluginConfig are the configuration settings that can be defined for the
// mongodb_lungo plugin in porter.yaml
type PluginConfig struct {
	StoreType StoreTypeEnum `mapstructure:"storeType"`
	Path      string        `mapstructure:"path,omitempty"`
	Timeout   int           `mapstructure:"timeout,omitempty"`
}

func NewPlugin(cxt *portercontext.Context, pluginConfig interface{}) (store plugins.StoragePlugin, err error) {
	cfg := PluginConfig{
		StoreType: "MEMORY",
		Path:      "",
		Timeout:   10,
	}
	if err = mapstructure.Decode(pluginConfig, &cfg); err != nil {
		return nil, errors.Wrapf(err, "error decoding %s plugin config from %#v", PluginKey, pluginConfig)
	}

	switch cfg.StoreType {
	case StoreTypeEnumMemory:
		store = NewMemoryStore(cxt, cfg)
	case StoreTypeEnumFile:
		store = NewFileStore(cxt, cfg)
	default:
		return nil, errors.New("unknown StoreType specified")
	}

	return store, nil
}
