// Package config loads yaml configs
package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seanenck/blap/internal/cli"
	"github.com/seanenck/blap/internal/config/types"
	"gopkg.in/yaml.v3"
)

func doDecode[T any](in string, o T) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(o); err != nil {
		return fmt.Errorf("file: %s -> %v", in, err)
	}
	return nil
}

// Load will initialize the configuration from a file
func Load(input string, context cli.Settings) (Configuration, error) {
	c := Configuration{}
	c.handler = &processHandler{}
	c.context = context
	c.Applications = make(map[string]types.Application)
	if err := doDecode(input, &c); err != nil {
		return c, err
	}
	if len(c.Include) > 0 {
		var including []string
		for _, i := range c.Include {
			r := context.Resolve(i)
			res := []string{r}
			c.context.LogDebug("including: %s\n", i)
			if strings.Contains(r, "*") {
				globbed, err := filepath.Glob(r)
				if err != nil {
					return c, err
				}
				res = globbed
			}
			including = append(including, res...)
		}
		for _, include := range including {
			c.context.LogDebug("loading included: %s\n", include)
			type included struct {
				Applications map[string]types.Application `yaml:"applications"`
				Disable      bool                         `yaml:"disable"`
			}
			var apps included
			if err := doDecode(include, &apps); err != nil {
				return c, err
			}
			if apps.Disable {
				continue
			}
			for k, v := range apps.Applications {
				if _, ok := c.Applications[k]; ok {
					return c, fmt.Errorf("%s is overwritten by config: %s", k, include)
				}
				c.Applications[k] = v
			}
		}
	}
	canFilter := context.FilterApplications()
	sub := make(map[string]types.Application)
	for n, a := range c.Applications {
		if a.Disable {
			continue
		}
		allowed := true
		if canFilter {
			allowed = context.AllowApplication(n)
		}
		if allowed {
			sub[n] = a
		}
	}
	c.Applications = sub
	return c, nil
}
