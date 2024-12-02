package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aeilang/urlshortener/config"
	_ "github.com/lib/pq"
)

func InitDB(cfg config.DataBaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed open db: %w", err)
	}
	log.Println(cfg.DSN())

	// 配置数据库连接池
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}
