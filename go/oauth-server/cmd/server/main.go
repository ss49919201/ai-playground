package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/ss49919201/ai-kata/go/oauth-server/internal/adapters/http" // エイリアスを使用
	"github.com/ss49919201/ai-kata/go/oauth-server/internal/adapters/storage"
	"github.com/ss49919201/ai-kata/go/oauth-server/internal/app"
	"github.com/ss49919201/ai-kata/go/oauth-server/internal/config"
)

func main() {
	// --- 設定ファイルのパスをコマンドライン引数から取得 ---
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// --- 設定の読み込み ---
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// --- 依存関係の初期化 (DI: Dependency Injection) ---
	// アダプター層の初期化
	clientRepo := storage.NewInMemoryClientRepository()
	userRepo := storage.NewInMemoryUserRepository()
	codeRepo := storage.NewInMemoryAuthorizationCodeRepository()
	tokenRepo := storage.NewInMemoryTokenRepository()

	clock := storage.SystemClock{}
	hasher := storage.NewBcryptHasher(0) // bcryptのデフォルトコストを使用
	idGen := storage.UUIDGenerator{}
	codeIssuer := storage.RandomCodeIssuer{}
	tokenIssuer := storage.RandomTokenIssuer{} // または JWTTokenIssuer を実装して使用

	// アプリケーションサービス層の初期化 (アダプターを注入)
	authServiceConfig := app.AuthServiceConfig{
		AuthCodeLifetime:    cfg.Token.AuthCodeLifetime,
		AccessTokenLifetime: cfg.Token.AccessTokenLifetime, // インプリシットフロー用
	}
	authService := app.NewAuthService(
		clientRepo, userRepo, codeRepo, tokenRepo, hasher, codeIssuer, tokenIssuer, clock, authServiceConfig,
	)

	tokenServiceConfig := app.TokenServiceConfig{
		AccessTokenLifetime:  cfg.Token.AccessTokenLifetime,
		RefreshTokenLifetime: cfg.Token.RefreshTokenLifetime,
	}
	tokenService := app.NewTokenService(
		clientRepo, userRepo, codeRepo, tokenRepo, hasher, tokenIssuer, clock, tokenServiceConfig,
	)

	clientService := app.NewClientService(
		clientRepo, idGen, hasher, clock,
	)

	// HTTPアダプター (サーバー) の初期化 (サービスを注入)
	httpServer := httpadapter.NewServer(authService, tokenService, clientService)

	// --- HTTPサーバーの設定と起動 ---
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      httpServer, // httpadapter.Server が http.Handler を実装
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// サーバーをゴルーチンで起動
	go func() {
		log.Printf("OAuth 2.0 サーバーを %s で起動します...", addr)
		var serverErr error
		if cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
			log.Println("TLS を有効にして起動します。")
			serverErr = server.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
		} else {
			log.Println("TLS を無効にして起動します (開発用途)。")
			serverErr = server.ListenAndServe()
		}
		// ListenAndServe は正常終了時以外は常にエラーを返す
		if serverErr != nil && serverErr != http.ErrServerClosed {
			log.Fatalf("サーバーの起動/実行に失敗しました: %v", serverErr)
		}
	}()

	// --- Graceful Shutdown の設定 ---
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) と SIGTERM を捕捉
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シグナルを受信するまでブロック
	<-quit
	log.Println("サーバーをシャットダウンします...")

	// シャットダウン処理のためのコンテキストを作成 (タイムアウト付き)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30秒以内にシャットダウン
	defer cancel()

	// サーバーに新しいリクエストの受け付けを停止させ、現在の処理が終わるのを待つ
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーの Graceful Shutdown に失敗しました: %v", err)
	}

	log.Println("サーバーは正常にシャットダウンしました。")
}
