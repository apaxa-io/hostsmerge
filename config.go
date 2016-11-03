package main

import (
	"encoding/json"
	"errors"
	"os"
)

type config struct {
	URLs                  []string // List of URLs for fetch
	Excludes              []string // FQDNs for excluding from result file
	IP                    string   // IP address assigned to all FQDN in result file
	Resolvers             []string // List of DNS resolvers for checking FQDNs
	ResolveViaTCP         bool     // True => access to resolvers via TCP, otherwise use UDP
	ResolveTimeout        int      // Timeout (in seconds) per DNS resolver communication (for dial, write & read)
	NegativeTTL           uint32   // TTL used for caching negative DNS answer (no IP address assigned for FQDN)
	MinPositiveTTL        uint32   // TTL used for caching positive DNS answer (some IP address assigned for FQDN) if it own TTL is too small
	ServFailTTL           uint32   // TTL used for caching DNS answer with SERVFAIL error (usually mean the same as negative answer)
	SkipOnResolveServFail bool     // Do not add FQDN to result file if resolver return SERVFAIL
	SkipOnResolveError    bool     // Do not add FQDN to result file if some other error occurred in resolver
	NumParsers            int      // Number of goroutines used for fetching & parsing URLs
	NumResolvers          int      // Number if goroutines used for DNS resolution
	CacheFile             string   // Path to file with DNS cache (used to speedup)
	TargetFile            string   // Path to target file
	ChanLength            int      // Length of channels used for passing FQDNs between workers
}

func validateConfig(conf *config) error {
	if len(conf.URLs) < 1 {
		return errors.New("Should be at least 1 url")
	}
	if len(conf.IP) < 1 {
		return errors.New("Ip should be set")
	}
	if len(conf.Resolvers) < 1 {
		return errors.New("Should be at least 1 resolver")
	}
	if conf.ResolveTimeout < 1 {
		return errors.New("ResolveTimeout should be >= 1")
	}
	if conf.NegativeTTL < 0 {
		return errors.New("NegativeTtl should be >= 0")
	}
	if conf.NumParsers < 1 {
		return errors.New("NumParsers should be >= 1")
	}
	if conf.NumResolvers < 1 {
		return errors.New("NumResolvers should be >= 1")
	}
	if conf.ChanLength < 1 {
		return errors.New("ChanLength should be >= 1")
	}
	if len(conf.CacheFile) == 0 {
		conf.CacheFile = cachePath()
	}
	return nil
}

func locateConfig() (string, error) {
	switch len(os.Args) {
	case 1:
		return configPath(), nil
	case 2:
		return os.Args[1], nil
	default:
		return "", errors.New("Bad usage. Use \"" + os.Args[0] + " [<path-to-config-file>]\"")
	}
}

func readConfig() (conf config, err error) {
	fn, err := locateConfig()
	if err != nil {
		return
	}
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		return
	}
	if err = json.NewDecoder(f).Decode(&conf); err != nil {
		return
	}
	err = validateConfig(&conf)

	for i := range conf.Excludes {
		conf.Excludes[i] = normalize(conf.Excludes[i])
	}

	return
}
