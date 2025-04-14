package ports

import (
	"time"
	// "github.com/ss49919201/ai-kata/go/oauth-server/internal/domain" // 必要に応じて
)

// --- 副作用を抽象化するインターフェース ---

// Clock は現在時刻を取得する機能を提供します。
// これにより、テスト時に時刻を固定または操作することが可能になります。
type Clock interface {
	Now() time.Time
}

// PasswordHasher はパスワードのハッシュ化と比較を行う機能を提供します。
// bcrypt や argon2 など、具体的なアルゴリズムの実装を隠蔽します。
type PasswordHasher interface {
	// Hash は平文のパスワードを受け取り、ハッシュ化された文字列とエラーを返します。
	Hash(password string) (string, error)
	// Compare はハッシュ化されたパスワードと平文のパスワードを受け取り、
	// 一致するかどうかを示すブール値とエラーを返します。
	// 不一致の場合はエラーではなく false を返します。
	Compare(hashedPassword, password string) (bool, error)
}

// IDGenerator は一意な識別子 (UUIDなど) やランダムな文字列 (シークレットなど) を
// 生成する機能を提供します。
type IDGenerator interface {
	// Generate は一意なID (例: ClientID, UserID) を生成します。
	Generate() (string, error)
	// GenerateSecret はクライアントシークレットなどのランダムな秘密文字列を生成します。
	GenerateSecret() (string, error)
}

// CodeIssuer は認可コードフローで使用される一時的な認可コードの値を生成します。
// 推測困難なランダム文字列を生成する必要があります。
type CodeIssuer interface {
	IssueCode() (string, error)
}

// TokenIssuer はアクセストークンやリフレッシュトークンの値を生成します。
// 実装によっては、JWTの生成と署名、または単純なランダム文字列の生成を行います。
// JWTの場合は、検証機能も提供することがあります。
type TokenIssuer interface {
	// IssueToken はトークンとして使用する文字列 (JWTまたはランダム文字列) を生成します。
	IssueToken() (string, error)

	// --- JWT を使用する場合のオプションメソッド ---
	// IssueJWT は指定された情報を含むJWTを生成し、署名して返します。
	// IssueJWT(claims JWTPayload) (string, error)
	// Verify はJWT文字列を検証し、ペイロードとエラーを返します。
	// Verify(tokenValue string) (JWTPayload, error)
}

/*
// JWTPayload は TokenIssuer が JWT を扱う場合に、
// トークンに含めるクレームを表す構造体です (例)。
type JWTPayload struct {
	Issuer    string           `json:"iss"` // 発行者 (サーバー自身)
	Subject   string           `json:"sub"` // 主体 (ユーザーIDなど)
	Audience  []string         `json:"aud"` // 対象者 (クライアントIDなど)
	ExpiresAt int64            `json:"exp"` // 有効期限 (Unixタイムスタンプ)
	IssuedAt  int64            `json:"iat"` // 発行日時 (Unixタイムスタンプ)
	NotBefore int64            `json:"nbf"` // 有効開始日時 (オプション)
	JwtID     string           `json:"jti"` // JWT ID (オプション)
	ClientID  domain.ClientID  `json:"cid"` // カスタムクレーム: クライアントID
	Scope     string           `json:"scp"` // カスタムクレーム: スコープ (スペース区切り文字列)
	// 他のカスタムクレーム...
}
*/
