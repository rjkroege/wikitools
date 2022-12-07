package wiki

// TODO(rjk): Should be in the "config" directory.

import (
	"encoding/json"
	"fmt"
	"os"
)

// Toplevel settings.
type Settings struct {
	Wikidir       string                    `json:"wikidir"`
	TemplateForTag    map[string]string `json:"templatefortag"`
}

// Read opens a json format configuration file.
func Read(path string) (*Settings, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no config file %q: %v", path, err)
	}

	settings := &Settings{}
	decoder := json.NewDecoder(fd)
	if err := decoder.Decode(settings); err != nil {
		return nil, fmt.Errorf("error parsing config %q: %v", path, err)
	}

	// TODO(rjk): Validate the configurable settings.
	return settings, nil
}
