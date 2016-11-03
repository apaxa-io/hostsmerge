package main

func configPath() string {
	const suffix = "/etc/io.apaxa.hostsmerge.config.json"
	return suffix
}

func cachePath() string {
	const suffix = "/var/cache/io.apaxa.hostsmerge.cache.json"
	return suffix
}
