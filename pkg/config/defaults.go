package config

var (
	registeredDefaults = make(map[string]interface{})
	// default aliases
	registeredAliases = map[string]string{
		"PluginsDirectory": "plugins-directory",
	}
)

func AddDefaultValue(key string, val interface{}) {
	registeredDefaults[key] = val
}

func AddAlias(alias, orig string) {
	registeredAliases[alias] = orig
}
