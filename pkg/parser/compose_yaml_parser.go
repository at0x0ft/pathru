package parser

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/mount"
	"gopkg.in/yaml.v3"
	"strings"
)

type ComposeYamlParser struct {
	Content string
}

type ComposeYamlTop struct {
	Services map[string]ComposeYamlService `yaml:"services"`
	Volumes  map[string]interface{}        `yaml:"volumes"`
}

type ComposeYamlService struct {
	Volumes []interface{} `yaml:"volumes"`
}

func (mp *ComposeYamlParser) Parse() (map[string]mount.BindMount, error) {
	var t ComposeYamlTop
	if err := yaml.Unmarshal([]byte(mp.Content), &t); err != nil {
		return nil, err
	}

	res := make(map[string]mount.BindMount)
	for n, s := range t.Services {
		for _, rv := range s.Volumes {
			switch v := rv.(type) {
			case string:
				m, isBind, err := mp.parseShort(v, t.Volumes)
				if err != nil {
					return nil, err
				} else if !isBind {
					continue
				}
				res[n] = m
				break
			case map[string]interface{}:
				m, isBind, err := mp.parseLong(v)
				if err != nil {
					return nil, err
				} else if !isBind {
					continue
				}
				res[n] = m
				break
			default:
				return nil, fmt.Errorf(
					"wrong compose volume format [\"%s\"]",
					rv,
				)
			}
		}
	}

	return res, nil
}

func (mp *ComposeYamlParser) parseShort(content string, volumes map[string]interface{}) (mount.BindMount, bool, error) {
	paths := strings.Split(content, ":")
	if len(paths) < 2 {
		return mount.BindMount{}, false, fmt.Errorf(
			"wrong syntax in compose volume [\"%s\"]",
			content,
		)
	}
	src, tgt := paths[0], paths[1]
	if _, exists := volumes[src]; exists {
		return mount.BindMount{}, false, nil
	}
	return mount.BindMount{src, tgt}, true, nil
}

func (mp *ComposeYamlParser) parseLong(content map[string]interface{}) (mount.BindMount, bool, error) {
	if t, ok := content["type"]; !ok {
		return mount.BindMount{}, false, fmt.Errorf(
			"wrong compose volume format (not found \"type\" key in long syntax) [\"%s\"]",
			content,
		)
	} else if t != "bind" {
		return mount.BindMount{}, false, nil
	}

	rs, ok1 := content["source"]
	rt, ok2 := content["target"]
	if !(ok1 && ok2) {
		return mount.BindMount{}, false, fmt.Errorf(
			"wrong compose bind volume format (not found \"source\" or \"target\" key in long syntax) [\"%s\"]",
			content,
		)
	}

	s, ok3 := rs.(string)
	t, ok4 := rt.(string)
	if !(ok3 && ok4) {
		return mount.BindMount{}, false, fmt.Errorf(
			"failed to cast bind volume source or target [source = \"%s\", target = \"%s\"]",
			rs,
			rt,
		)
	}
	return mount.BindMount{s, t}, true, nil
}
