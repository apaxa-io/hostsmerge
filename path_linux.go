package main

func configPath() string {
	const suffix = "/etc/hostsmerge.config.json"
	return suffix
}

func cachePath() string {
	const suffix = "/var/cache/hostsmerge.cache.json"
	return suffix
}
