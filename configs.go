package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

type Configs struct {
	Directory string
	Patterns  []string
	Command   []string
}

var configfile = flag.String("config", ".wu.json", "Config file")
var directory = flag.String("dir", "", "Directory to watch")
var pattern = flag.String("pattern", "", "Patterns to filter filenames")
var saveconf = flag.Bool("save", false, "Save options to conf")

func getConfigs() Configs {
	flag.Parse()

	conf := readConfigFile()

	if dir := parseDirectory(); dir != "" {
		conf.Directory = dir
	}

	if patterns := parsePatterns(); patterns != nil {
		conf.Patterns = patterns
	}

	if command := parseCommand(); command != nil {
		conf.Command = command
	}

	if *saveconf {
		saveConfigFile(conf)
	}

	return conf
}

func readConfigFile() Configs {
	file, err := os.Open(*configfile)
	defer file.Close()

	if err == nil {
		log.Println("Reading options from", *configfile)
		var conf Configs
		if err := json.NewDecoder(file).Decode(&conf); err != nil {
			log.Fatalln("Failed to parse config file:", err)
		}
		return conf
	}
	return Configs{".", []string{"*"}, []string{}}
}

func saveConfigFile(conf Configs) {
	log.Println("Saving options to", *configfile)
	file, err := os.Create(*configfile)
	defer file.Close()

	if err != nil {
		log.Fatalln("Failed to open config file:", err)
	}
	if bytes, err := json.MarshalIndent(conf, "", "  "); err == nil {
		if _, err := file.Write(bytes); err != nil {
			log.Fatalln("Failed to write config file:", err)
		}
	} else {
		log.Fatalln("Failed to encode options:", err)
	}
}

func parseDirectory() string {
	dir := *directory
	if info, err := os.Stat(dir); err == nil {
		if !info.IsDir() {
			log.Fatal(dir, "is not a directory")
		}
	}
	return dir
}

func parsePatterns() []string {
	pat := strings.Trim(*pattern, " ")
	if pat == "" {
		return nil
	}

	patternSep, _ := regexp.Compile("[,\\s]+")

	patternMap := make(map[string]bool)
	ret := []string{}

	for _, part := range patternSep.Split(pat, -1) {
		patternMap[part] = true
	}
	for part := range patternMap {
		ret = append(ret, part)
	}

	return ret
}

func parseCommand() []string {
	if flag.NArg() == 0 {
		return nil
	}
	return flag.Args()
}
