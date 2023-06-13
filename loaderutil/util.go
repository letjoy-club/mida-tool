package loaderutil

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/ttlcache"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func NewItemLoader[K comparable, V any](
	db *gorm.DB,
	loader func(ctx context.Context, keys []K) ([]V, error),
	dataMaper func(k map[K]V, v V),
	duration time.Duration,
) *dataloader.Loader[K, V] {
	c := ttlcache.New[K, V](duration)
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[V] {
		data, err := loader(ctx, keys)

		if err != nil {
			return lo.Map(data, func(m V, i int) *dataloader.Result[V] {
				return &dataloader.Result[V]{Error: err}
			})
		}

		dataMap := map[K]V{}
		for _, m := range data {
			dataMaper(dataMap, m)
		}

		return lo.Map(keys, func(key K, i int) *dataloader.Result[V] {
			ret, itemNotFound := dataMap[key]
			if itemNotFound {
				return &dataloader.Result[V]{Data: ret}
			} else {
				return &dataloader.Result[V]{Error: midacode.ErrItemNotFound}
			}
		})
	}, dataloader.WithCache[K, V](&c),
	)
}
