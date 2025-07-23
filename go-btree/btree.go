// package btree はBツリーデータ構造の実装を提供します
package btree

// BTree は次数（degree）で定義されるBツリー構造体です
// Bツリーは自己平衡型の木構造で、データベースやファイルシステムで使用されます
type BTree struct {
	degree int   // ノードが持つことができる最大キー数は (2*degree-1)
	root   *Node // 根ノードへのポインタ
}

// Node はBツリーの各ノードを表す構造体です
type Node struct {
	keys     []int   // ノードに格納されているキー値の配列（ソート済み）
	children []*Node // 子ノードへのポインタ配列（葉ノードの場合はnil）
	isLeaf   bool    // 葉ノード（子を持たない）かどうかのフラグ
	parent   *Node   // 親ノードへのポインタ（根ノードの場合はnil）
}

// NewBTree は指定された次数でBツリーを作成します
// degree: ノードの最小子数（2以上である必要があります）
// 各ノードは最低degree-1個、最大2*degree-1個のキーを持ちます
func NewBTree(degree int) *BTree {
	return &BTree{
		degree: degree, // 次数を設定（通常は3以上の値を使用）
	}
}

// Insert はBツリーに新しいキーを挿入します
// key: 挿入するキー値（整数）
// Bツリーの性質を保持しながら挿入を行います
func (bt *BTree) Insert(key int) {
	// ケース1: 空のツリーの場合、新しい根ノードを作成
	if bt.root == nil {
		bt.root = &Node{
			keys:   []int{key}, // 最初のキーとして挿入
			isLeaf: true,        // 根ノードは葉ノードでもある
		}
		return
	}
	
	// ケース2: 根ノードが満杯の場合、根を分割して新しい根を作成
	if len(bt.root.keys) == 2*bt.degree-1 {
		// 新しい根ノード（内部ノード）を作成
		newRoot := &Node{
			isLeaf: false, // 内部ノードなので葉ではない
		}
		// 現在の根を新しい根の子として設定
		newRoot.children = append(newRoot.children, bt.root)
		bt.root.parent = newRoot
		// 子ノード（元の根）を分割
		newRoot.splitChild(0, bt.degree)
		// 新しい根をツリーの根として設定
		bt.root = newRoot
	}
	
	// ケース3: 根ノードが満杯でない場合、通常の挿入処理
	bt.root.insertNonFull(key, bt.degree)
}

// insertNonFull は満杯でないノードにキーを挿入します
// key: 挿入するキー値
// degree: Bツリーの次数
// 前提条件: ノードは満杯ではない（キー数 < 2*degree-1）
func (n *Node) insertNonFull(key int, degree int) {
	// 最後のキーのインデックスから開始
	i := len(n.keys) - 1
	
	if n.isLeaf {
		// 葉ノードの場合: キーを直接挿入
		// 配列を1つ拡張してスペースを作る
		n.keys = append(n.keys, 0)
		// 挿入位置を見つけるまで要素を右にシフト
		for i >= 0 && key < n.keys[i] {
			n.keys[i+1] = n.keys[i] // 要素を右にシフト
			i--
		}
		// 適切な位置にキーを挿入
		n.keys[i+1] = key
	} else {
		// 内部ノードの場合: 適切な子ノードを見つけて再帰的に挿入
		// 挿入すべき子ノードのインデックスを見つける
		for i >= 0 && key < n.keys[i] {
			i--
		}
		i++ // 子ノードのインデックスに調整
		
		// 対象の子ノードが満杯の場合、先に分割
		if len(n.children[i].keys) == 2*degree-1 {
			n.splitChild(i, degree)
			// 分割後、挿入するキーがどちらの子に行くかを判定
			if key > n.keys[i] {
				i++ // 右の子ノードを選択
			}
		}
		// 適切な子ノードに再帰的に挿入
		n.children[i].insertNonFull(key, degree)
	}
}

// splitChild は満杯の子ノードを2つに分割します
// index: 分割する子ノードのインデックス
// degree: Bツリーの次数
// 中央のキーを親ノードに上げ、左右のキーを別ノードに分割
func (n *Node) splitChild(index int, degree int) {
	// 分割対象の満杯ノード
	fullChild := n.children[index]
	// 新しく作成する右側のノード
	newChild := &Node{
		isLeaf: fullChild.isLeaf, // 元ノードと同じタイプ（葉/内部）
		parent: n,                // 親ノードを設定
	}
	
	// 中央インデックス（degree-1番目のキーが中央）
	midIndex := degree - 1
	
	// 右側ノードに中央より右のキーをコピー
	newChild.keys = make([]int, len(fullChild.keys[midIndex+1:]))
	copy(newChild.keys, fullChild.keys[midIndex+1:])
	
	if !fullChild.isLeaf {
		// 内部ノードの場合、子ノードも分割
		// 右側ノードに中央より右の子ノードをコピー
		newChild.children = make([]*Node, len(fullChild.children[midIndex+1:]))
		copy(newChild.children, fullChild.children[midIndex+1:])
		// コピーした子ノードの親ポインタを更新
		for _, child := range newChild.children {
			if child != nil {
				child.parent = newChild
			}
		}
		// 元ノードの子ノード配列を中央でカット
		fullChild.children = fullChild.children[:midIndex+1]
	}
	
	// 中央のキーを取り出し（これが親に上がる）
	midKey := fullChild.keys[midIndex]
	// 元ノードのキー配列を中央でカット
	fullChild.keys = fullChild.keys[:midIndex]
	
	// 親ノードの子ノード配列に新ノードを挿入
	n.children = append(n.children, nil)
	copy(n.children[index+2:], n.children[index+1:]) // 要素を右にシフト
	n.children[index+1] = newChild                   // 新ノードを挿入
	
	// 親ノードのキー配列に中央キーを挿入
	n.keys = append(n.keys, 0)
	copy(n.keys[index+1:], n.keys[index:]) // 要素を右にシフト
	n.keys[index] = midKey                 // 中央キーを挿入
}

// Search はBツリー内で指定されたキーを検索します
// key: 検索するキー値
// 戻り値: キーが見つかった場合true、そうでなけれfalse
// 時間計算量: O(log n)
func (bt *BTree) Search(key int) bool {
	// 空のツリーの場合は見つからない
	if bt.root == nil {
		return false
	}
	// 根ノードから検索を開始
	return bt.root.search(key)
}

// search はノードから再帰的にキーを検索します
// key: 検索するキー値
// 戻り値: キーが見つかった場合true、そうでなけれfalse
// Bツリーの検索アルゴリズムを実装
func (n *Node) search(key int) bool {
	// ステップ1: 現在ノード内でキーの位置を特定
	i := 0
	// キーより小さい値の範囲をスキップ
	for i < len(n.keys) && key > n.keys[i] {
		i++
	}
	
	// ステップ2: キーが現在ノードで見つかったかチェック
	if i < len(n.keys) && key == n.keys[i] {
		return true // キーを発見！
	}
	
	// ステップ3: 葉ノードの場合、これ以上探す場所がない
	if n.isLeaf {
		return false // キーは存在しない
	}
	
	// ステップ4: 内部ノードの場合、適切な子ノードで再帰検索
	return n.children[i].search(key)
}

// Delete はBツリーから指定されたキーを削除します
// key: 削除するキー値
// 注意: この実装は簡単化された版です（完全なBツリー削除は非常に複雑）
func (bt *BTree) Delete(key int) {
	// 空のツリーの場合は何もしない
	if bt.root == nil {
		return
	}
	// 根ノードから削除を開始
	bt.root.delete(key)
	
	// 根ノードが空になった場合の処理
	if len(bt.root.keys) == 0 && !bt.root.isLeaf {
		// 子ノードを新しい根として昇格
		bt.root = bt.root.children[0]
		if bt.root != nil {
			bt.root.parent = nil // 親ポインタをクリア
		}
	}
}

// delete はノードから再帰的にキーを削除します
// key: 削除するキー値
// 注意: この実装は簡単化された版で、葉ノードからのみ削除します
func (n *Node) delete(key int) {
	// ステップ1: 現在ノード内でキーの位置を特定
	i := 0
	// キーより小さい値の範囲をスキップ
	for i < len(n.keys) && key > n.keys[i] {
		i++
	}
	
	// ステップ2: キーが現在ノードで見つかった場合
	if i < len(n.keys) && key == n.keys[i] {
		if n.isLeaf {
			// 葉ノードの場合: キーを直接削除
			// 削除位置以降の要素を左にシフト
			copy(n.keys[i:], n.keys[i+1:])
			// 配列のサイズを1つ減らす
			n.keys = n.keys[:len(n.keys)-1]
		}
		// 内部ノードからの削除は未実装（非常に複雑なため）
		return
	}
	
	// ステップ3: 内部ノードの場合、適切な子ノードで再帰削除
	if !n.isLeaf && i < len(n.children) {
		n.children[i].delete(key)
	}
}