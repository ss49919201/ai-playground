// package btree_test はBツリーのテストケースを提供します
// TDD手法にRed-Green-Refactorサイクルで開発されたテスト群
package btree

import "testing"

// TestNewBTree はBTreeのコンストラクタの基本動作をテストします
// テスト内容: NewBTreeがnon-nilのBTreeインスタンスを返すことを確認
func TestNewBTree(t *testing.T) {
	// 次数03のBTreeを作成（各ノードは最大5個までのキーを保持可能）
	bt := NewBTree(3)
	if bt == nil {
		t.Error("NewBTree should return a non-nil BTree")
	}
}

// TestEmptyBTreeSearch は空のBTreeでの検索動作をテストします
// テスト内容: 空のツリーではどのキーも見つからないことを確認
func TestEmptyBTreeSearch(t *testing.T) {
	// 空のBTreeを作成
	bt := NewBTree(3)
	// 任意のキーで検索しても見つからないことを確認
	if bt.Search(42) {
		t.Error("Empty BTree should not find any key")
	}
}

// TestInsertSingleElement は単一要素の挿入と検索をテストします
// テスト内容: 1. 挿入した要素が検索できること 2. 挿入していない要素は見つからないこと
func TestInsertSingleElement(t *testing.T) {
	// 次数03のBTreeを作成
	bt := NewBTree(3)
	// キー42を挿入
	bt.Insert(42)
	
	// 挿入したキー42が検索できることを確認
	if !bt.Search(42) {
		t.Error("Should find inserted element")
	}
	
	// 挿入していないキー999は見つからないことを確認
	if bt.Search(999) {
		t.Error("Should not find non-inserted element")
	}
}

// TestInsertMultipleElements は複数要素の挿入と検索をテストします
// テスト内容: 複数のキーをランダムな順序で挿入し、すべて検索できることを確認
func TestInsertMultipleElements(t *testing.T) {
	// 次数03のBTreeを作成
	bt := NewBTree(3)
	// テスト用のキー配列（意図的にランダムな順序）
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
	
	// すべてのキーを挿入
	for _, key := range keys {
		bt.Insert(key)
	}
	
	// 挿入したすべてのキーが検索できることを確認
	for _, key := range keys {
		if !bt.Search(key) {
			t.Errorf("Should find inserted element %d", key)
		}
	}
	
	// 挿入していないキーは見つからないことを確認
	if bt.Search(999) {
		t.Error("Should not find non-inserted element")
	}
}

// TestNodeSplitting はノード分割機能をテストします
// テスト内容: ノードが満杯になるまで要素を挿入し、分割後もすべての要素が検索できることを確認
// 次数03の場合、各ノードは最大5個のキーを保持でき、それを超えると分割が発生
func TestNodeSplitting(t *testing.T) {
	// 次数03のBTreeを作成
	bt := NewBTree(3)
	
	// 1から10までの連続した数値を挿入（ノード分割を発生させる）
	for i := 1; i <= 10; i++ {
		bt.Insert(i)
	}
	
	// 分割後もすべての要素が正しく検索できることを確認
	for i := 1; i <= 10; i++ {
		if !bt.Search(i) {
			t.Errorf("Should find element %d after node splitting", i)
		}
	}
}

// TestDelete は要素の削除機能をテストします
// テスト内容: 1. 削除した要素が検索できなくなること 2. 他の要素は影響を受けないこと
// 注意: 現在の実装は簡単化された版で、葉ノードからのみ削除可能
func TestDelete(t *testing.T) {
	// 次数03のBTreeを作成
	bt := NewBTree(3)
	// テスト用のキー配列
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
	
	// すべてのキーを挿入
	for _, key := range keys {
		bt.Insert(key)
	}
	
	// キー6を削除
	bt.Delete(6)
	// 削除したキー6が見つからないことを確認
	if bt.Search(6) {
		t.Error("Should not find deleted element 6")
	}
	
	// 削除していない他のキーはまだ検索できることを確認
	for _, key := range []int{10, 20, 5, 12, 30, 7, 17} {
		if !bt.Search(key) {
			t.Errorf("Should still find element %d after deletion", key)
		}
	}
}