package main

func configPath() string {
	const suffix = "\\AppData\\Roaming\\apaxa.io\\hostsmerge\\config.json"
	u, err := user.Current()
	if err != nil {
		return suffix
	}
	return u.HomeDir + suffix
}

func cachePath() string {
	const suffix = "\\AppData\\Local\\Publishers\\apaxa.io\\LocalCache\\hostsmerge.cache.json"
	u, err := user.Current()
	if err != nil {
		return suffix
	}
	return u.HomeDir + suffix
}
