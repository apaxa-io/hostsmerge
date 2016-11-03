package main

import "os/user"

func configPath() string {
	const suffix = "/Library/Preferences/io.apaxa.hostsmerge.config.json"
	u, err := user.Current()
	if err != nil {
		return suffix
	}
	return u.HomeDir + suffix
}

func cachePath() string {
	const suffix = "/Library/Caches/io.apaxa.hostsmerge.cache.json"
	u, err := user.Current()
	if err != nil {
		return suffix
	}
	return u.HomeDir + suffix
}
