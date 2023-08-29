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

// NewAggregatorLoader 用于聚合数据的 loader
// 例如获取用户的所有 Tag，那么 K 为用户 user id 类型，T 为单个 UserTag 类型，V 为聚合后的 UserTags 类型
func NewAggregatorLoader[K comparable, T any, V any](
	db *gorm.DB,
	loader func(ctx context.Context, keys []K) ([]T, error),
	groupBy func(m map[K]V, item T),
	duration time.Duration,
	options ...func(*Option[K, V]),
) *dataloader.Loader[K, V] {
	option := Option[K, V]{duration: duration}
	for _, optModifier := range options {
		optModifier(&option)
	}
	c := ttlcache.New[K, V](duration)
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[V] {
		items, err := loader(ctx, keys)
		if err != nil {
			return lo.Map(keys, func(m K, i int) *dataloader.Result[V] {
				return &dataloader.Result[V]{Error: err}
			})
		}

		dataMap := map[K]V{}
		for _, item := range items {
			groupBy(dataMap, item)
		}

		notFoundKeys := []K{}
		notFoundIndexes := []int{}

		result := lo.Map(keys, func(key K, i int) *dataloader.Result[V] {
			ret, itemFound := dataMap[key]
			if itemFound {
				return &dataloader.Result[V]{Data: ret}
			}
			if option.createIfNotFound != nil {
				notFoundKeys = append(notFoundKeys, key)
				notFoundIndexes = append(notFoundIndexes, i)
				// 如果找不到数据，且有创建函数，则记录 key 和 index，等待后续创建
				return &dataloader.Result[V]{}
			}
			if option.placeholderFunc != nil {
				// 如果找不到数据，且有占位函数，则调用占位函数
				item, err := option.placeholderFunc(ctx, key)
				return &dataloader.Result[V]{Data: item, Error: err}
			}
			return &dataloader.Result[V]{Error: midacode.ErrItemNotFound.(midacode.Error2).WithExtra(key)}
		})

		if len(notFoundKeys) > 0 {
			// 如果有找不到的数据，则调用创建函数
			createdItems, errs := option.createIfNotFound(ctx, notFoundKeys)
			for i := range createdItems {
				itemIndex := notFoundIndexes[i]
				var err error
				if len(errs) > i {
					err = errs[i]
				}
				result[itemIndex] = &dataloader.Result[V]{Data: createdItems[i], Error: err}
			}
		}
		return result
	}, dataloader.WithCache[K, V](&c))
}

// NewItemLoader 用于单个数据的 loader，实际上是 NewAggregatorLoader 的特例
func NewItemLoader[K comparable, V any](
	db *gorm.DB,
	loader func(ctx context.Context, keys []K) ([]V, error),
	dataMaper func(k map[K]V, v V),
	duration time.Duration,
	options ...func(*Option[K, V]),
) *dataloader.Loader[K, V] {
	return NewAggregatorLoader(db, loader, dataMaper, duration, options...)
}

type Option[K comparable, V any] struct {
	duration         time.Duration
	createIfNotFound func(ctx context.Context, ids []K) ([]V, []error)
	placeholderFunc  func(ctx context.Context, id K) (V, error)
}

// CreateIfNotFound 如果找不到数据，会调用这个函数
func CreateIfNotFound[K comparable, V any](createIfNotFound func(ctx context.Context, ids []K) ([]V, []error)) func(*Option[K, V]) {
	return func(o *Option[K, V]) {
		o.createIfNotFound = createIfNotFound
	}
}

// Placeholder 如果找不到数据，会调用这个函数，返回一个占位数据
func Placeholder[K comparable, V any](placeholderFunc func(ctx context.Context, id K) (V, error)) func(*Option[K, V]) {
	return func(o *Option[K, V]) {
		o.placeholderFunc = placeholderFunc
	}
}
