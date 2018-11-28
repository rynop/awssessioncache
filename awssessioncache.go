// Package awssessioncache provides a cache of aws/session by region.
// This package is concurrency safe.
//
// Usage:
// 	import sc "github.com/rynop/awssessioncache"
//	...
//	sess, err := sc.Get(&sc.Conf{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	uploader := s3manager.NewUploader(sess)
//	...
package awssessioncache

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type atomicSessionMap struct {
	sessionsByRegion map[string]*session.Session
	rwLock           sync.RWMutex
}

func (m *atomicSessionMap) get(region string) (*session.Session, bool) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	val, found := m.sessionsByRegion[region]
	return val, found
}

func (m *atomicSessionMap) set(region string, val *session.Session) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.sessionsByRegion[region] = val
}

// Conf is the config for a session cache.
// If region is omitted, will use the `AWS_REGION` env var
type Conf struct {
	region string
}

var sessionByRegionCache = atomicSessionMap{
	sessionsByRegion: make(map[string]*session.Session),
}

// Get an aws sdk session by region.
// If Conf is empty, will use default region defined in `AWS_REGION` env var
func Get(c *Conf) (*session.Session, error) {
	if c.region == "" {
		// no region passed? use default
		region, _ := os.LookupEnv("AWS_REGION")
		c.region = region
	}

	if s, exists := sessionByRegionCache.get(c.region); exists {
		return s, nil
	}

	sess, err := session.NewSession(&aws.Config{Region: &c.region})
	if err != nil {
		return nil, err
	}
	sessionByRegionCache.set(c.region, sess)
	return sess, nil
}
