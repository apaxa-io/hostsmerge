package main

import (
	"github.com/miekg/dns"
	"log"
	"math"
	"math/rand"
	"time"
)

var cache dnsCache

func initCache(c config) (err error) {
	err = cache.ReadFromFile(c.CacheFile)
	cache.CleanOutdated()
	return
}

func closeCache(c config) error {
	return cache.WriteToFile(c.CacheFile)
}

func resolveDNS(fqdn string, conf config) (exists bool, ttl time.Duration, err error) {
	//log.Println("LookUp")
	m := dns.Msg{}
	m.SetQuestion(fqdn+".", dns.TypeA)
	proto := "udp"
	if conf.ResolveViaTCP {
		proto = "tcp"
	}
	c := dns.Client{Net: proto, Timeout: time.Duration(conf.ResolveTimeout) * time.Second}

	for _, i := range rand.Perm(len(conf.Resolvers)) {
		var r *dns.Msg
		r, _, err = c.Exchange(&m, conf.Resolvers[i]+":53")

		switch {
		case err != nil:
			continue
		case len(r.Answer) == 0 || r.Rcode == dns.RcodeNameError:
			return false, time.Duration(conf.NegativeTTL) * time.Second, nil
		case len(r.Answer) > 0 && r.Rcode == dns.RcodeSuccess:

			// Find min TTL in answer
			t := uint32(math.MaxUint32)
			for _, ans := range r.Answer {
				if ans.Header().Ttl < t {
					t = ans.Header().Ttl
				}
			}

			// Fix ttl
			if t < conf.MinPositiveTTL {
				t = conf.MinPositiveTTL
			}

			return true, time.Duration(t) * time.Second, nil
		case r.Rcode == dns.RcodeServerFailure:
			return !conf.SkipOnResolveServFail, time.Duration(conf.ServFailTTL) * time.Second, nil
		default:
			continue
		}
	}

	return !conf.SkipOnResolveError, 0, err
}

func resolve(fqdn string, conf config) (exists bool, err error) {
	inCache, exists := cache.Get(fqdn)
	if !inCache {
		var ttl time.Duration
		exists, ttl, err = resolveDNS(fqdn, conf)
		if err == nil {
			cache.Add(fqdn, exists, time.Now().Add(ttl))
		}
	}
	return
}

func resolver(requestChan <-chan string, responseChan chan<- string, conf config) {
	for fqdn := range requestChan {
		exists, err := resolve(fqdn, conf)
		if err != nil {
			log.Printf("Resolver: error with %v : %v\n", fqdn, err)
		}
		if exists {
			//log.Println("Resolver: add " + fqdn)
			responseChan <- fqdn
		} else {
			//log.Println("Resolver: skip " + fqdn)
		}
	}
}
