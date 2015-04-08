package config

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
	"os"
	"path"
	"strings"
)

const configName = ".gsckconfig"
const defaultSection = "default"
const envPrefix = "GSCK"

var configPath string
var conf *ini.File

func init() {
	homeDir, _ := homedir.Dir()
	configPath = path.Join(homeDir, configName)
	_, e := os.Lstat(configPath)
	if os.IsNotExist(e) {
		conf = ini.Empty()
		conf.SaveTo(configPath)
	} else {
		conf, _ = ini.Load(configPath)
	}
	setDefaultFromMap(map[string]string{
		"user":          os.Getenv("USER"),
		"retry":         "2",
		"formatter":     "ansi",
		"local.tmpdir":  "/tmp",
		"remote.tmpdir": "/tmp",
		"json.pretty":   "true",
	})
	conf.SaveTo(configPath)
	bindEnv()
}

func setDefault(raw, value string) {
	section, key := splitKey(raw)
	if section.Key(key).String() == "" {
		section.NewKey(key, value)
	}
}

func setDefaultFromMap(pair map[string]string) {
	for k, v := range pair {
		setDefault(k, v)
	}
}

// bindEnv sets ENV GSCK_X with value of default.X
// Actual effect, if written in bash: export GSCK_USER=default.user
func bindEnv() {
	section := conf.Section(defaultSection)
	hash := section.KeysHash()
	for k, v := range hash {
		envKey := envPrefix + "_" + strings.ToUpper(k)
		os.Setenv(envKey, v)
	}
}

func splitKey(raw string) (section *ini.Section, key string) {
	fields := strings.Split(raw, ".")
	length := len(fields)
	sectionName := defaultSection
	if length > 1 {
		sectionName = fields[0]
	}
	section = conf.Section(sectionName)
	key = fields[length-1]
	return
}

// GetString returns string value for key
// Example: default.user
func GetString(key string) (value string) {
	section, sectionKey := splitKey(key)
	value = strings.TrimSpace(section.Key(sectionKey).String())
	return
}

// GetInt returns int value for key
// Example: default.retry
func GetInt(key string) (value int) {
	section, sectionKey := splitKey(key)
	value, _ = section.Key(sectionKey).Int()
	return
}

// GetBool returns bool value for key
// Example: json.pretty
func GetBool(key string) (value bool) {
	section, sectionKey := splitKey(key)
	value, _ = section.Key(sectionKey).Bool()
	return
}

// Set writes setting to config file
func Set(key, value string) {
	section, sectionKey := splitKey(key)
	section.NewKey(sectionKey, value)
	conf.SaveTo(configPath)
}
