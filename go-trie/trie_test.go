package main

import "testing"

// TestNewTrie はトライ木の作成をテストする
func TestNewTrie(t *testing.T) {
	trie := NewTrie()
	if trie == nil {
		t.Error("NewTrie() should return a non-nil Trie")
	}
}

// TestInsert は単語の挿入と基本的な検索機能をテストする
func TestInsert(t *testing.T) {
	trie := NewTrie()
	
	// 単語を挿入
	trie.Insert("hello")
	
	// 挿入した単語が見つかることを確認
	if !trie.Search("hello") {
		t.Error("Expected 'hello' to be found after insertion")
	}
	
	// 挿入していない単語が見つからないことを確認
	if trie.Search("world") {
		t.Error("Expected 'world' to not be found")
	}
}

func TestStartsWith(t *testing.T) {
	trie := NewTrie()
	
	// Insert words
	trie.Insert("hello")
	trie.Insert("help")
	trie.Insert("her")
	trie.Insert("here")
	
	// Test prefix search
	if !trie.StartsWith("he") {
		t.Error("Expected prefix 'he' to be found")
	}
	
	if !trie.StartsWith("hel") {
		t.Error("Expected prefix 'hel' to be found")
	}
	
	if !trie.StartsWith("hello") {
		t.Error("Expected prefix 'hello' to be found")
	}
	
	if trie.StartsWith("world") {
		t.Error("Expected prefix 'world' to not be found")
	}
	
	if trie.StartsWith("hi") {
		t.Error("Expected prefix 'hi' to not be found")
	}
}

func TestEmptyString(t *testing.T) {
	trie := NewTrie()
	
	// Empty string should be handled correctly
	trie.Insert("")
	if !trie.Search("") {
		t.Error("Expected empty string to be found after insertion")
	}
	
	if !trie.StartsWith("") {
		t.Error("Expected empty prefix to always return true")
	}
}

func TestUnicodeSupport(t *testing.T) {
	trie := NewTrie()
	
	// Test with Unicode characters
	trie.Insert("こんにちは")
	trie.Insert("今日")
	
	if !trie.Search("こんにちは") {
		t.Error("Expected 'こんにちは' to be found")
	}
	
	if !trie.StartsWith("こん") {
		t.Error("Expected prefix 'こん' to be found")
	}
	
	if trie.Search("こんばんは") {
		t.Error("Expected 'こんばんは' to not be found")
	}
}

func TestPartialWordAsPrefix(t *testing.T) {
	trie := NewTrie()
	
	trie.Insert("hello")
	
	// "hell" is a prefix of "hello" but not a complete word
	if !trie.StartsWith("hell") {
		t.Error("Expected 'hell' to be found as prefix")
	}
	
	if trie.Search("hell") {
		t.Error("Expected 'hell' to not be found as complete word")
	}
}