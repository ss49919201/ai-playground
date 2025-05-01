package ports

import (
	"context"

	"github.com/ss49919201/ai-playground/go/oauth-server/internal/domain"
)

// ClientRepository はクライアント情報の永続化を抽象化するインターフェースです。
type ClientRepository interface {
	// Save は指定されたクライアント情報を永続化します。
	// 既に存在するクライアントIDの場合は更新します。
	Save(ctx context.Context, client domain.Client) error

	// FindByID は指定されたクライアントIDに対応するクライアント情報を取得します。
	// 見つからない場合はエラーを返します (例: ErrNotFound)。
	FindByID(ctx context.Context, id domain.ClientID) (domain.Client, error)

	// List は登録されているすべてのクライアント情報を取得します。
	// (管理機能などで使用。大量になる可能性があるため注意)
	// List(ctx context.Context) ([]domain.Client, error)

	// Delete は指定されたクライアントIDのクライアント情報を削除します。
	// Delete(ctx context.Context, id domain.ClientID) error
}

// UserRepository はユーザー情報の永続化を抽象化するインターフェースです。
type UserRepository interface {
	// Save は指定されたユーザー情報を永続化します。
	// 既に存在するユーザーIDの場合は更新します。
	Save(ctx context.Context, user domain.User) error

	// FindByID は指定されたユーザーIDに対応するユーザー情報を取得します。
	// 見つからない場合はエラーを返します (例: ErrNotFound)。
	FindByID(ctx context.Context, id domain.UserID) (domain.User, error)

	// FindByUsername は指定されたユーザー名に対応するユーザー情報を取得します。
	// 見つからない場合はエラーを返します (例: ErrNotFound)。
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

// AuthorizationCodeRepository は認可コードの永続化を抽象化するインターフェースです。
// 認可コードは一時的なものであり、一度使用されるか有効期限が切れると削除されるべきです。
type AuthorizationCodeRepository interface {
	// Save は指定された認可コード情報を永続化します。
	Save(ctx context.Context, code domain.AuthorizationCode) error

	// FindByValue は指定された認可コードの値に対応する情報を取得します。
	// 見つからない場合や、既に使用済み/有効期限切れの場合はエラーを返すことが期待されます。
	// (ただし、有効期限切れチェックはサービス層で行う方が一貫性がある場合もあります)
	FindByValue(ctx context.Context, value string) (domain.AuthorizationCode, error)

	// Delete は指定された認可コードの値を削除します。
	// 認可コードが使用された後に呼び出されます。
	Delete(ctx context.Context, value string) error
}

// TokenRepository はアクセストークンとリフレッシュトークンの永続化を抽象化するインターフェースです。
type TokenRepository interface {
	// Save は指定されたトークン情報 (アクセスまたはリフレッシュ) を永続化します。
	Save(ctx context.Context, token domain.Token) error

	// FindByValue は指定されたトークンの値に対応する情報を取得します。
	// アクセストークンとリフレッシュトークンの両方を検索対象とします。
	// 見つからない場合はエラーを返します (例: ErrNotFound)。
	FindByValue(ctx context.Context, value string) (domain.Token, error)

	// Delete は指定されたトークンの値を削除します。
	// トークンが失効された場合や、リフレッシュトークンがローテーションされる場合に呼び出されます。
	Delete(ctx context.Context, value string) error

	// FindByUserAndClient は特定のユーザーとクライアントに発行されたトークンを取得します。
	// (例: 同一ユーザー/クライアントへの同時セッション数を制限する場合などに使用)
	// FindByUserAndClient(ctx context.Context, userID domain.UserID, clientID domain.ClientID) ([]domain.Token, error)
}

// TODO: 標準的なエラー型 (例: ErrNotFound) を定義する
// var ErrNotFound = errors.New("resource not found")
