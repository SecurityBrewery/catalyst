package database

import (
	"context"
	"database/sql"
	"errors"
)

func Paginate(ctx context.Context, f func(ctx context.Context, offset, limit int64) (nextPage bool, err error)) error {
	const pageSize int64 = 100

	for i := range int64(1000) {
		nextPage, err := f(ctx, i*pageSize, pageSize)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No more features to process, exit the loop
				return nil
			}

			return err
		}

		if !nextPage {
			return nil
		}
	}

	return errors.New("pagination limit reached, too many pages")
}

func PaginateItems[T any](ctx context.Context, f func(ctx context.Context, offset, limit int64) (items []T, err error)) ([]T, error) {
	var allItems []T

	if err := Paginate(ctx, func(ctx context.Context, offset, limit int64) (nextPage bool, err error) {
		items, err := f(ctx, offset, limit)
		if err != nil {
			return false, err
		}

		allItems = append(allItems, items...)

		return len(items) == int(limit), nil
	}); err != nil {
		return nil, err
	}

	return allItems, nil
}
