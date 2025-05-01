package domain

import (
	"errors"
	"strings"
)

// Scope はアクセス権限を表す値オブジェクトです。
// OAuth 2.0 では、スペース区切りの文字列として表現されることが多いです。
type Scope string

// ValidateScope はスペース区切りのスコープ文字列を解析し、
// Scopeのスライスとエラーを返します。
// 空白文字を含むスコープ名や重複するスコープはエラーとします。
// 空文字列は有効な入力で、空のスライスを返します。
func ValidateScope(scopeStr string) ([]Scope, error) {
	// 空文字列の場合は空のスライスを返す
	if scopeStr == "" {
		return []Scope{}, nil
	}

	// スペースで分割
	parts := strings.Split(scopeStr, " ")
	scopes := make([]Scope, 0, len(parts))
	seen := make(map[Scope]struct{}) // 重複チェック用

	for _, part := range parts {
		// 空の要素は無視 (連続するスペースなど)
		if part == "" {
			continue
		}

		scope := Scope(part)

		// スコープ名の形式チェック (RFC 6749 Section 3.3: VSCHAR (%x21-7E) のみ)
		// ここでは簡単なチェックとして、制御文字やスペースが含まれていないか確認
		if strings.ContainsAny(part, " \t\n\r\x00\x1F\x7F") {
			return nil, errors.New("スコープ名に無効な文字が含まれています: " + part)
		}

		// 重複チェック
		if _, exists := seen[scope]; exists {
			return nil, errors.New("スコープが重複しています: " + part)
		}

		scopes = append(scopes, scope)
		seen[scope] = struct{}{}
	}

	// 解析結果が空の場合 (例: " " のような入力)
	if len(scopes) == 0 && scopeStr != "" {
		// これをエラーとするか、空スライスを返すかは仕様による
		// ここでは空スライスを返す
	}

	return scopes, nil
}

// FormatScopes は Scope のスライスをスペース区切りの文字列に変換します。
// スライスが空またはnilの場合は空文字列を返します。
func FormatScopes(scopes []Scope) string {
	if len(scopes) == 0 {
		return ""
	}
	strs := make([]string, len(scopes))
	for i, s := range scopes {
		strs[i] = string(s)
	}
	return strings.Join(strs, " ")
}
