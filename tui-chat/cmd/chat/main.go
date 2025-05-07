package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sakaeshinya/tui-chat/internal/app"
	"github.com/sakaeshinya/tui-chat/internal/config"
	"github.com/sakaeshinya/tui-chat/internal/domain"
	"github.com/sakaeshinya/tui-chat/internal/network"
	"github.com/sakaeshinya/tui-chat/internal/ui"
)

func main() {
	// コマンドライン引数の解析
	var (
		username   = flag.String("u", "", "ユーザー名")
		isServer   = flag.Bool("s", false, "サーバーモードで起動")
		isClient   = flag.Bool("c", false, "クライアントモードで起動")
		address    = flag.String("a", "", "接続先アドレス（クライアントモード）またはリッスンアドレス（サーバーモード）")
		showHelp   = flag.Bool("h", false, "ヘルプを表示")
		configPath = flag.String("config", "", "設定ファイルのパス")
	)
	flag.Parse()

	// ヘルプの表示
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// 設定の読み込み
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// ユーザー名の設定
	if *username == "" {
		*username = cfg.Username
	}

	// ユーザーの作成
	user, err := domain.NewUser(*username)
	if err != nil {
		log.Fatalf("ユーザーの作成に失敗しました: %v", err)
	}

	// ネットワークコンポーネントの作成
	var networkComponent network.NetworkPort
	if *isServer {
		networkComponent = network.NewTCPServer()
	} else {
		networkComponent = network.NewTCPClient()
	}

	// UIコントローラの作成（ハンドラは後で設定）
	controller := ui.NewTUIController(nil, user)

	// サービスの作成
	service := app.NewChatService(networkComponent, controller, user)

	// イベントハンドラの作成
	handler := app.NewChatEventHandler(service, controller)

	// ハンドラの設定
	controller.SetHandler(handler)

	// サーバーモードまたはクライアントモードで自動起動
	if *isServer && *address != "" {
		// サーバーモードで起動
		if *address == "" {
			*address = fmt.Sprintf(":%s", cfg.DefaultPort)
		}
		err := service.StartServer(*address)
		if err != nil {
			log.Fatalf("サーバーの起動に失敗しました: %v", err)
		}
		controller.SetMode(ui.ModeServer)
	} else if *isClient && *address != "" {
		// クライアントモードで起動
		if *address == "" {
			*address = fmt.Sprintf("%s:%s", cfg.DefaultHost, cfg.DefaultPort)
		}
		err := service.JoinChat(*address)
		if err != nil {
			log.Fatalf("サーバーへの接続に失敗しました: %v", err)
		}
		controller.SetMode(ui.ModeClient)
	}

	// UIの起動
	if err := controller.Start(); err != nil {
		log.Fatalf("UIの起動に失敗しました: %v", err)
	}
}

// loadConfig は設定を読み込むヘルパー関数
func loadConfig(path string) (config.Config, error) {
	if path == "" {
		var err error
		path, err = config.GetConfigPath()
		if err != nil {
			return config.DefaultConfig, err
		}
	}
	return config.LoadConfig(path)
}

// printHelp はヘルプを表示する関数
func printHelp() {
	fmt.Println("TUIチャットアプリケーション")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  chat [オプション]")
	fmt.Println("")
	fmt.Println("オプション:")
	fmt.Println("  -u <name>     ユーザー名を指定")
	fmt.Println("  -s            サーバーモードで起動")
	fmt.Println("  -c            クライアントモードで起動")
	fmt.Println("  -a <address>  接続先アドレス（クライアントモード）またはリッスンアドレス（サーバーモード）")
	fmt.Println("  -config <path> 設定ファイルのパス")
	fmt.Println("  -h            ヘルプを表示")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  サーバーモード: chat -s -a :8080")
	fmt.Println("  クライアントモード: chat -c -a localhost:8080")
}
