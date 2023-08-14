package bucket

import (
	"context"
	"sync"
	"time"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
)

type MapValueQuantity map[string]int

type bucket struct {
	mu     *sync.Mutex
	limit  int
	values MapValueQuantity
}

type MapKeyBucket map[string]bucket

type LeakyBucket struct {
	ticker  *time.Ticker
	buckets MapKeyBucket
}

func NewLeakyBucket(cfg *config.Config) *LeakyBucket {
	return &LeakyBucket{
		ticker: time.NewTicker(60 * time.Second),
		buckets: MapKeyBucket{
			common.LoginBucketKey: bucket{
				mu:     &sync.Mutex{},
				limit:  cfg.Limitations.Login,
				values: make(MapValueQuantity),
			},
			common.PassBucketKey: bucket{
				mu:     &sync.Mutex{},
				limit:  cfg.Limitations.Password,
				values: make(MapValueQuantity),
			},
			common.IPBucketKey: bucket{
				mu:     &sync.Mutex{},
				limit:  cfg.Limitations.IP,
				values: make(MapValueQuantity),
			},
		},
	}
}

func (l *LeakyBucket) Add(req common.APIAuthRequest) bool {
	return l.addLogin(req.Login) && l.addPassword(req.Password) && l.addIP(req.IP)
}

func (l *LeakyBucket) addLogin(login string) bool {
	return l.addInBucket(common.LoginBucketKey, login)
}

func (l *LeakyBucket) addPassword(password string) bool {
	return l.addInBucket(common.PassBucketKey, password)
}

func (l *LeakyBucket) addIP(ip string) bool {
	return l.addInBucket(common.IPBucketKey, ip)
}

func (l *LeakyBucket) addInBucket(bucketKey, valueKey string) bool {
	bucket, ok := l.buckets[bucketKey]
	if !ok {
		return false
	}

	l.buckets[bucketKey].mu.Lock()
	defer l.buckets[bucketKey].mu.Unlock()

	// если лимит превышен - вход запрещен
	if l.buckets[bucketKey].values[valueKey] >= bucket.limit {
		return false
	}
	// иначе увеличиваем число запросов по этому адресу
	l.buckets[bucketKey].values[valueKey]++
	return true
}

func (l *LeakyBucket) ResetBucket() {
	for _, bucketKey := range []string{common.LoginBucketKey, common.PassBucketKey, common.IPBucketKey} {
		if bucket, ok := l.buckets[bucketKey]; ok {
			bucket.values = make(MapValueQuantity)
			l.buckets[bucketKey] = bucket
		}
	}
}

func (l *LeakyBucket) Repeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-l.ticker.C:
			l.ResetBucket()
		}
	}
}
