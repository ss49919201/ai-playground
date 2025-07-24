package main

// TrieNode はトライ木の各ノードを表す構造体
// 各ノードは子ノードへのマップと、単語の終端かどうかを示すフラグを持つ
type TrieNode struct {
	// children は子ノードを格納するマップ
	// キーは文字（rune型）、値は子ノードへのポインタ
	children map[rune]*TrieNode
	// isEnd はこのノードが単語の終端であることを示すフラグ
	// 例: "cat" と "car" が挿入されている場合、't' と 'r' のノードで true になる
	isEnd bool
}

// Trie はトライ木全体を表す構造体
// ルートノードへのポインタを保持する
type Trie struct {
	// root はトライ木のルートノード
	// 通常は空のノード（どの文字も表さない）
	root *TrieNode
}

// NewTrie は新しいトライ木を作成して返す
// 初期状態では空のルートノードのみを持つ
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
			isEnd:    false,
		},
	}
}

// Insert は指定された単語をトライ木に挿入する
// 時間計算量: O(m) - mは単語の長さ
// 空間計算量: O(m) - 新しいノードが必要な場合
func (t *Trie) Insert(word string) {
	// ルートノードから開始
	current := t.root
	
	// 単語の各文字について処理
	for _, char := range word {
		// 現在のノードに該当する文字の子ノードが存在しない場合
		if _, exists := current.children[char]; !exists {
			// 新しい子ノードを作成
			current.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isEnd:    false,
			}
		}
		// 次のノードに移動
		current = current.children[char]
	}
	// 単語の終端をマーク
	current.isEnd = true
}

// Search は指定された単語がトライ木に存在するかを検索する
// 時間計算量: O(m) - mは単語の長さ
// 空間計算量: O(1) - 追加のメモリは使用しない
func (t *Trie) Search(word string) bool {
	// ルートノードから開始
	current := t.root
	
	// 単語の各文字について処理
	for _, char := range word {
		// 現在のノードに該当する文字の子ノードが存在しない場合
		if _, exists := current.children[char]; !exists {
			// 単語は存在しない
			return false
		}
		// 次のノードに移動
		current = current.children[char]
	}
	// 最後のノードが単語の終端かどうかを返す
	// 例: "car" を検索する場合、'r' のノードの isEnd が true である必要がある
	return current.isEnd
}

// StartsWith は指定されたプレフィックスで始まる単語がトライ木に存在するかを確認する
// 時間計算量: O(p) - pはプレフィックスの長さ
// 空間計算量: O(1) - 追加のメモリは使用しない
func (t *Trie) StartsWith(prefix string) bool {
	// ルートノードから開始
	current := t.root
	
	// プレフィックスの各文字について処理
	for _, char := range prefix {
		// 現在のノードに該当する文字の子ノードが存在しない場合
		if _, exists := current.children[char]; !exists {
			// プレフィックスは存在しない
			return false
		}
		// 次のノードに移動
		current = current.children[char]
	}
	// プレフィックスが見つかった場合は true を返す
	// isEnd のチェックは不要（プレフィックスの存在のみを確認）
	return true
}