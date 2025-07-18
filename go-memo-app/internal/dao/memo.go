package dao

import (
	"go-memo-app/pkg/model"
	"sync"
	"time"
)

type (
	GetAllMemosFunc    func() []model.Memo
	GetMemoByIDFunc    func(id int) (*model.Memo, error)
	CreateMemoFunc     func(title, content string) (*model.Memo, error)
	UpdateMemoFunc     func(id int, title, content string) (*model.Memo, error)
	DeleteMemoFunc     func(id int) error
)

var (
	memos  = make(map[int]*model.Memo)
	nextID = 1
	mu     sync.RWMutex
)

func NewInMemoryGetAllMemos() GetAllMemosFunc {
	return func() []model.Memo {
		mu.RLock()
		defer mu.RUnlock()

		memoList := make([]model.Memo, 0, len(memos))
		for _, memo := range memos {
			memoList = append(memoList, *memo)
		}
		return memoList
	}
}

func NewInMemoryGetMemoByID() GetMemoByIDFunc {
	return func(id int) (*model.Memo, error) {
		mu.RLock()
		defer mu.RUnlock()

		memo, exists := memos[id]
		if !exists {
			return nil, ErrMemoNotFound
		}
		return memo, nil
	}
}

func NewInMemoryCreateMemo() CreateMemoFunc {
	return func(title, content string) (*model.Memo, error) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		memo := &model.Memo{
			ID:        nextID,
			Title:     title,
			Content:   content,
			CreatedAt: now,
			UpdatedAt: now,
		}

		memos[nextID] = memo
		nextID++

		return memo, nil
	}
}

func NewInMemoryUpdateMemo() UpdateMemoFunc {
	return func(id int, title, content string) (*model.Memo, error) {
		mu.Lock()
		defer mu.Unlock()

		memo, exists := memos[id]
		if !exists {
			return nil, ErrMemoNotFound
		}

		memo.Title = title
		memo.Content = content
		memo.UpdatedAt = time.Now()

		return memo, nil
	}
}

func NewInMemoryDeleteMemo() DeleteMemoFunc {
	return func(id int) error {
		mu.Lock()
		defer mu.Unlock()

		if _, exists := memos[id]; !exists {
			return ErrMemoNotFound
		}

		delete(memos, id)
		return nil
	}
}