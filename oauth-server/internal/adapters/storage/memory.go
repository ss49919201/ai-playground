package storage

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ss49919201/ai-playground/go/oauth-server/internal/domain"
)

// --- エラー定義 ---
var (
	ErrClientNotFound   = errors.New("クライアントが見つかりません")
	ErrUserNotFound     = errors.New("ユーザーが見つかりません")
	ErrCodeNotFound     = errors.New("認可コードが見つかりません")
	ErrTokenNotFound    = errors.New("トークンが見つかりません")
	ErrUsernameTaken    = errors.New("ユーザー名が既に使用されています")
	ErrDataInconsistent = errors.New("内部データ不整合")
)

// --- InMemoryClientRepository ---

// InMemoryClientRepository は ports.ClientRepository のインメモリ実装です。
// sync.RWMutex を使用してスレッドセーフ性を確保します。
type InMemoryClientRepository struct {
	mu      sync.RWMutex
	clients map[domain.ClientID]domain.Client
}

// NewInMemoryClientRepository は InMemoryClientRepository の新しいインスタンスを生成します。
func NewInMemoryClientRepository() *InMemoryClientRepository {
	return &InMemoryClientRepository{
		clients: make(map[domain.ClientID]domain.Client),
	}
}

// Save はクライアント情報をメモリに保存または更新します。
func (r *InMemoryClientRepository) Save(ctx context.Context, client domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// ドメインオブジェクトをそのまま保存 (コピーはドメイン層のファクトリで行われている想定)
	r.clients[client.ID] = client
	return nil
}

// FindByID は指定されたIDのクライアント情報をメモリから取得します。
func (r *InMemoryClientRepository) FindByID(ctx context.Context, id domain.ClientID) (domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	client, ok := r.clients[id]
	if !ok {
		return domain.Client{}, ErrClientNotFound
	}
	// 念のためコピーを返す (必須ではないが、より安全)
	// clientCopy := client // 構造体は値渡しなのでコピーされる
	return client, nil
}

// --- InMemoryUserRepository ---

// InMemoryUserRepository は ports.UserRepository のインメモリ実装です。
type InMemoryUserRepository struct {
	mu       sync.RWMutex
	users    map[domain.UserID]domain.User // UserID -> User
	username map[string]domain.UserID      // Username -> UserID (検索用インデックス)
}

// NewInMemoryUserRepository は InMemoryUserRepository の新しいインスタンスを生成します。
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:    make(map[domain.UserID]domain.User),
		username: make(map[string]domain.UserID),
	}
}

// Save はユーザー情報をメモリに保存または更新します。ユーザー名の重複チェックも行います。
func (r *InMemoryUserRepository) Save(ctx context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ユーザー名重複チェック (更新時も考慮)
	existingID, usernameExists := r.username[user.Username]
	if usernameExists && existingID != user.ID {
		return ErrUsernameTaken
	}

	// 古いユーザー名インデックスを削除 (更新の場合)
	if oldUser, ok := r.users[user.ID]; ok {
		if oldUser.Username != user.Username {
			delete(r.username, oldUser.Username)
		}
	}

	r.users[user.ID] = user
	r.username[user.Username] = user.ID
	return nil
}

// FindByID は指定されたIDのユーザー情報をメモリから取得します。
func (r *InMemoryUserRepository) FindByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}

// FindByUsername は指定されたユーザー名のユーザー情報をメモリから取得します。
func (r *InMemoryUserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.username[username]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	user, ok := r.users[id]
	if !ok {
		// データ不整合の可能性
		return domain.User{}, fmt.Errorf("%w: username %s points to non-existent user ID %s", ErrDataInconsistent, username, id)
	}
	return user, nil
}

// --- InMemoryAuthorizationCodeRepository ---

// InMemoryAuthorizationCodeRepository は ports.AuthorizationCodeRepository のインメモリ実装です。
type InMemoryAuthorizationCodeRepository struct {
	mu    sync.RWMutex
	codes map[string]domain.AuthorizationCode // Code Value -> AuthorizationCode
}

// NewInMemoryAuthorizationCodeRepository は InMemoryAuthorizationCodeRepository の新しいインスタンスを生成します。
func NewInMemoryAuthorizationCodeRepository() *InMemoryAuthorizationCodeRepository {
	return &InMemoryAuthorizationCodeRepository{
		codes: make(map[string]domain.AuthorizationCode),
	}
}

// Save は認可コード情報をメモリに保存します。
func (r *InMemoryAuthorizationCodeRepository) Save(ctx context.Context, code domain.AuthorizationCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.codes[code.Value] = code
	return nil
}

// FindByValue は指定された値の認可コード情報をメモリから取得します。
func (r *InMemoryAuthorizationCodeRepository) FindByValue(ctx context.Context, value string) (domain.AuthorizationCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	code, ok := r.codes[value]
	if !ok {
		return domain.AuthorizationCode{}, ErrCodeNotFound
	}
	// 有効期限チェックはリポジトリ層ではなく、サービス層で行う方が一貫性があることが多い
	// if code.IsExpired(time.Now()) {
	// 	return domain.AuthorizationCode{}, ErrCodeNotFound // 期限切れも Not Found として扱う
	// }
	return code, nil
}

// Delete は指定された値の認可コード情報をメモリから削除します。
func (r *InMemoryAuthorizationCodeRepository) Delete(ctx context.Context, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 存在しなくてもエラーにはしない (冪等性)
	delete(r.codes, value)
	return nil
}

// --- InMemoryTokenRepository ---

// InMemoryTokenRepository は ports.TokenRepository のインメモリ実装です。
type InMemoryTokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]domain.Token // Token Value -> Token
}

// NewInMemoryTokenRepository は InMemoryTokenRepository の新しいインスタンスを生成します。
func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]domain.Token),
	}
}

// Save はトークン情報 (アクセスまたはリフレッシュ) をメモリに保存します。
func (r *InMemoryTokenRepository) Save(ctx context.Context, token domain.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.Value] = token
	return nil
}

// FindByValue は指定された値のトークン情報をメモリから取得します。
func (r *InMemoryTokenRepository) FindByValue(ctx context.Context, value string) (domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	token, ok := r.tokens[value]
	if !ok {
		return domain.Token{}, ErrTokenNotFound
	}
	// 有効期限チェックはサービス層で行う
	return token, nil
}

// Delete は指定された値のトークン情報をメモリから削除します。
func (r *InMemoryTokenRepository) Delete(ctx context.Context, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 存在しなくてもエラーにはしない
	delete(r.tokens, value)
	return nil
}

// --- 副作用インターフェースのインメモリ実装 ---

// SystemClock は ports.Clock を実装します。
type SystemClock struct{}

// Now は現在の時刻を返します。
func (c SystemClock) Now() time.Time {
	return time.Now()
}

// BcryptHasher は ports.PasswordHasher を bcrypt を使用して実装します。
type BcryptHasher struct {
	Cost int // ハッシュ化のコスト (0の場合はデフォルト)
}

// NewBcryptHasher は BcryptHasher の新しいインスタンスを生成します。
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{Cost: cost}
}

// Hash はパスワードを bcrypt でハッシュ化します。
func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	return string(bytes), err
}

// Compare は bcrypt ハッシュと平文パスワードを比較します。
func (h *BcryptHasher) Compare(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true, nil // 一致
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil // 不一致 (エラーではない)
	}
	return false, err // その他のエラー (ハッシュ形式不正など)
}

// UUIDGenerator は ports.IDGenerator を UUID v4 を使用して実装します。
type UUIDGenerator struct{}

// Generate は新しい UUID v4 文字列を生成します。
func (g UUIDGenerator) Generate() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

// GenerateSecret は暗号学的に安全なランダム文字列を Base64 URL エンコードして生成します。
func (g UUIDGenerator) GenerateSecret() (string, error) {
	b := make([]byte, 32) // 256ビットのエントロピー
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("安全なランダム文字列の生成に失敗しました: %w", err)
	}
	// パディングなしの Base64 URL エンコーディングを使用
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// RandomCodeIssuer は ports.CodeIssuer を実装します。
type RandomCodeIssuer struct{}

// IssueCode は暗号学的に安全なランダム文字列を認可コードとして生成します。
func (i RandomCodeIssuer) IssueCode() (string, error) {
	b := make([]byte, 24) // 192ビットのエントロピー
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("認可コードの生成に失敗しました: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// RandomTokenIssuer は ports.TokenIssuer を実装します (JWTではない単純なランダム文字列)。
type RandomTokenIssuer struct{}

// IssueToken は暗号学的に安全なランダム文字列をトークン値として生成します。
func (i RandomTokenIssuer) IssueToken() (string, error) {
	b := make([]byte, 32) // 256ビットのエントロピー
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("トークン値の生成に失敗しました: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
