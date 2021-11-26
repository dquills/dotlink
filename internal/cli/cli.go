package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dquills/dotlink/internal/linker"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Watch     bool
	ConfigDir string

	PF linker.Config `yaml:"dotlink"`
}

func Run() error {
	var c Config
	parseFlags(&c)

	if c.ConfigDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("unable to get current working directory: %s", err.Error())
		}
		c.ConfigDir = dir
	} else {
		// TODO: move this to its own utils dir
		fullDir := linker.GetFullPath(c.ConfigDir)
		os.Chdir(fullDir)
	}

	err := parsePaths(c.ConfigDir, &c.PF)
	if err != nil {
		log.Fatalf("unable to parse dotlink.yaml: %s", err.Error())
	}

	if c.Watch {
		watch(&c)
	} else {
		c.PF.LinkAll()
	}
	return nil
}

func parseFlags(c *Config) {
	flag.BoolVar(&c.Watch, "w", false,
		"Watches the config dir for changes and updates symlinks accordingly. *NOT YET IMPLEMENTED*",
	)
	flag.StringVar(&c.ConfigDir, "d", "",
		"Sets the base directory to look for .dotlink.yaml",
	)
	flag.Parse()
}

// TODO: Method instead of func...
// TODO: return errors instead of log.Fatalf'ing everywhere
func parsePaths(d string, pf *linker.Config) error {
	filePath := d + "/dotlink.yaml"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("unable to find dotlink.yaml")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open %s", filePath)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("unable to parse %s", filePath)
	}

	err = yaml.Unmarshal(b, pf)
	if err != nil {
		log.Fatalf("unable to parse %s", filePath)
	}
	return nil
}

func watch(c *Config) error {
	log.Printf("'watch' is not yet implemented")
	return nil
}
