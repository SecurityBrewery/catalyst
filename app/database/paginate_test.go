package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginate_AllPages(t *testing.T) {
	calls := 0
	err := Paginate(t.Context(), func(ctx context.Context, offset, limit int64) (bool, error) {
		calls++
		if calls < 3 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 3, calls, "expected 3 calls")
}

func TestPaginate_EarlyStop(t *testing.T) {
	calls := 0
	err := Paginate(t.Context(), func(ctx context.Context, offset, limit int64) (bool, error) {
		calls++
		return false, nil
	})
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 1, calls, "expected 1 call")
}

func TestPaginate_Error(t *testing.T) {
	errTest := errors.New("fail")
	err := Paginate(t.Context(), func(ctx context.Context, offset, limit int64) (bool, error) {
		return false, errTest
	})
	assert.ErrorIs(t, err, errTest, "expected error")
}

func TestPaginate_NoRows(t *testing.T) {
	err := Paginate(t.Context(), func(ctx context.Context, offset, limit int64) (bool, error) {
		return false, sql.ErrNoRows
	})
	assert.NoError(t, err, "expected no error")
}

func TestPaginateItems(t *testing.T) {
	calls := 0
	f := func(ctx context.Context, offset, limit int64) ([]int, error) {
		calls++
		if offset >= 100 {
			return nil, sql.ErrNoRows
		}
		return []int{1}, nil
	}
	items, err := PaginateItems(t.Context(), f)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, []int{1}, items, "expected items to match")
	assert.Equal(t, 2, calls, "expected 2 calls")
}

func TestPaginateItemsLarge(t *testing.T) {
	calls := 0
	f := func(ctx context.Context, offset, limit int64) ([]int, error) {
		calls++
		if offset >= 200 {
			return nil, sql.ErrNoRows
		}
		return []int{1}, nil
	}
	items, err := PaginateItems(t.Context(), f)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, []int{1, 1}, items, "expected items to match")
	assert.Equal(t, 3, calls, "expected 3 calls")
}
