package config

import (
	"fmt"
	. "main/src/utils"
	"os"

	"github.com/tkanos/gonfig"
)

func GetConfig(configPath string) ProgramSettings {
	var ps ProgramSettings
	_, err := os.Stat(configPath)
	if err != nil {
		fmt.Println("Failed to open configuration file")
		fmt.Errorf("error: %v", err)
		return ps
	}
	gonfig.GetConf(configPath, &ps)
	return ps
}
