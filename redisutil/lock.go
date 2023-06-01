package redisutil

import (
	"context"
	"sort"
	"time"

	"github.com/letjoy-club/mida-tool/midacode"
)

func LockAll(ctx context.Context, keys ...string) (release func(context.Context), err error) {
	sort.Strings(keys)
	locker := GetLocker(ctx)

	releasers := []func(context.Context) error{}

	for _, key := range keys {
		if key == "" {
			continue
		}
		lock, err := locker.Obtain(ctx, key, time.Second*10, nil)
		if err != nil {
			if len(releasers) > 0 {
				for _, releaser := range releasers {
					releaser(ctx)
				}
			}
			return nil, midacode.ErrResourceBusy
		}
		releasers = append(releasers, lock.Release)
	}

	return func(ctx context.Context) {
		for _, releaser := range releasers {
			releaser(ctx)
		}
	}, nil
}
