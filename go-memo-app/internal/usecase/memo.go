package usecase

import (
	"go-memo-app/internal/dao"
	"go-memo-app/pkg/model"
)

type (
	GetAllMemosFunc    = dao.GetAllMemosFunc
	GetMemoByIDFunc    = dao.GetMemoByIDFunc
	CreateMemoFunc     func(title, content string) (*model.Memo, error)
	UpdateMemoFunc     func(id int, title, content string) (*model.Memo, error)
	DeleteMemoFunc     = dao.DeleteMemoFunc
)

func NewGetAllMemos(getAllMemos dao.GetAllMemosFunc) GetAllMemosFunc {
	return getAllMemos
}

func NewGetMemoByID(getMemoByID dao.GetMemoByIDFunc) GetMemoByIDFunc {
	return getMemoByID
}

func NewCreateMemo(createMemo dao.CreateMemoFunc) CreateMemoFunc {
	return func(title, content string) (*model.Memo, error) {
		if title == "" {
			return nil, ErrInvalidInput
		}
		return createMemo(title, content)
	}
}

func NewUpdateMemo(updateMemo dao.UpdateMemoFunc) UpdateMemoFunc {
	return func(id int, title, content string) (*model.Memo, error) {
		if title == "" {
			return nil, ErrInvalidInput
		}
		return updateMemo(id, title, content)
	}
}

func NewDeleteMemo(deleteMemo dao.DeleteMemoFunc) DeleteMemoFunc {
	return deleteMemo
}