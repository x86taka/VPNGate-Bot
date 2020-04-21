package lib

type Host struct {
	HostName string `gorm:"primary_key"`
	IP       string `gorm:"size:255"`
	PORT     int
}

type DbConfig struct {
	// データベース
	Dialect string

	// ユーザー名
	DBUser string

	// パスワード
	DBPass string

	// プロトコル
	DBProtocol string

	// DB名
	DBName string
}

