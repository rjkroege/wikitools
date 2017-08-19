package wiki

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// readConfigurationImpl is a more easily testable implementation
// core of ReadConfiguration.
func readConfigurationImpl(basedir string) (map[string]string, error) {
	configfile := filepath.Join(basedir, ".wikinewrc")
	fd, err := os.Open(configfile)
	if err != nil {
		return map[string]string{}, err
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	config := map[string]string{}
	if err := decoder.Decode(&config); err != nil {
		return map[string]string{}, err
	}
	return config, nil
}

// ReadConfiguration opens a json format configuration file
// and returns a map of template name, template file pairs.
func ReadConfiguration() map[string]string {
	homedir := os.ExpandEnv("$HOME")
	if homedir == "" {
		return map[string]string{}
	}
	config, err := readConfigurationImpl(homedir)
	if err != nil {
		log.Println("couldn't read config file because", err)
		return map[string]string{}
	}
	return config
}
