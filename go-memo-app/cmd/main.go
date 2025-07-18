package main

import (
	"fmt"
	"go-memo-app/internal/dao"
	"go-memo-app/internal/handler"
	"go-memo-app/internal/usecase"
	"log"
	"net/http"
)

func main() {
	getAllMemosDAO := dao.NewInMemoryGetAllMemos()
	getMemoByIDDAO := dao.NewInMemoryGetMemoByID()
	createMemoDAO := dao.NewInMemoryCreateMemo()
	updateMemoDAO := dao.NewInMemoryUpdateMemo()
	deleteMemoDAO := dao.NewInMemoryDeleteMemo()

	getAllMemosUC := usecase.NewGetAllMemos(getAllMemosDAO)
	getMemoByIDUC := usecase.NewGetMemoByID(getMemoByIDDAO)
	createMemoUC := usecase.NewCreateMemo(createMemoDAO)
	updateMemoUC := usecase.NewUpdateMemo(updateMemoDAO)
	deleteMemoUC := usecase.NewDeleteMemo(deleteMemoDAO)

	memoHandler := handler.NewMemoHandler(
		getAllMemosUC,
		getMemoByIDUC,
		createMemoUC,
		updateMemoUC,
		deleteMemoUC,
	)

	http.HandleFunc("/api/memos", memoHandler)
	http.HandleFunc("/api/memos/", memoHandler)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}