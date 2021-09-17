package config

import (
	"crypto/md5"
	"errors"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"io/ioutil"
)

// Namespace storing at most one configuration section by name and sub-section.
type Namespace struct {
	name   string `config:"name,required"`
	config *ucfg.Config
}

var configOpts = []ucfg.Option{
	ucfg.PathSep("."),
	ucfg.ResolveEnv,
	ucfg.VarExp,
	ucfg.PrependValues,
}

// NewConfigFrom creates a new Config object from the given input.
// From can be any kind of structured data (struct, map, array, slice).
//
// If from is a string, the contents is treated like raw YAML input. The string
// will be parsed and a structure config object is build from the parsed
// result.
func NewConfigFrom(from interface{}) (*ucfg.Config, error) {
	if str, ok := from.(string); ok {
		c, err := yaml.NewConfig([]byte(str), configOpts...)
		return c, err
	}

	c, err := ucfg.NewFrom(from, configOpts...)
	return c, err
}

func NewConfigWithYAML(in []byte, source string) (*ucfg.Config, error) {
	opts := append(
		[]ucfg.Option{
			ucfg.MetaData(ucfg.Meta{Source: source}),
		},
		configOpts...,
	)
	c, err := yaml.NewConfig(in, opts...)
	return c, err
}

// OverwriteConfigOpts allow to change the globally set config option
func OverwriteConfigOpts(options []ucfg.Option) {
	configOpts = options
}

func LoadFile(path string) (*ucfg.Config, [md5.Size]byte, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, [md5.Size]byte{}, err
	}
	hash := md5.Sum(bs)
	c, err := yaml.NewConfig(bs, configOpts...)
	if err != nil {
		return nil, hash, err
	}
	return c, hash, err
}

// Unpack unpacks a configuration with at most one sub object. An sub object is
// ignored if it is disabled by setting `enabled: false`. If the configuration
// passed contains multiple active sub objects, Unpack will return an error.
func (ns *Namespace) Unpack(cfg *ucfg.Config) error {
	fields := cfg.GetFields()
	if len(fields) == 0 {
		return nil
	}

	var (
		err   error
		found bool
	)

	for _, name := range fields {
		var sub *ucfg.Config

		sub, err = cfg.Child(name, -1)
		if err != nil {
			// element is no configuration object -> continue so a namespace
			// Config unpacked as a namespace can have other configuration
			// values as well
			continue
		}

		if ns.name != "" {
			return errors.New("more than one namespace configured")
		}

		ns.name = name
		ns.config = sub
		found = true
	}

	if !found {
		return err
	}
	return nil
}

// Name returns the configuration sections it's name if a section has been set.
func (ns *Namespace) Name() string {
	return ns.name
}

// Config return the sub-configuration section if a section has been set.
func (ns *Namespace) Config() *ucfg.Config {
	return ns.config
}

// IsSet returns true if a sub-configuration section has been set.
func (ns *Namespace) IsSet() bool {
	return ns.config != nil
}
