package dbutil

import (
	"context"
	sql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

type DBConf struct {
	DSN    string `yaml:"dsn"`
	Prefix string `yaml:"prefix"`
}

func (c DBConf) ConnectDB() *gorm.DB {
	dsnObj, err := sql.ParseDSN(c.DSN)
	if err != nil {
		log.Panic(err)
	}
	dsnObj.ParseTime = true
	dsnObj.Collation = "utf8mb4_general_ci"
	dsn := dsnObj.FormatDSN() + "&charset=utf8mb4&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{TablePrefix: c.Prefix},
	})
	if err != nil {
		panic(err)
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")
	return db
}

type dbKey struct{}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

func GetDB(ctx context.Context) *gorm.DB {
	return ctx.Value(dbKey{}).(*gorm.DB)
}
