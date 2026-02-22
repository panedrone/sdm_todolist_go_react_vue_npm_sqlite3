package dbal

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	//"gorm.io/gorm/logger"
)

var _ds = &_DS{} // private for this package

func (ds *_DS) initDb() (err error) {
	ds.rootDb, err = gorm.Open(sqlite.Open("./db/todolist.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	DS = func(ctx context.Context) DataStore {
		return ds
	}
	return
}

// WithContext can be used in a middleware to start a separate session for each incoming web-request
//
//  api := r.Group("", func(ctx *gin.Context) {
//		ctx.Set("db", dal.WithContext(ctx))
//	})

//func WithContext(ctx context.Context) *gorm.DB {
//	//	https://gorm.io/docs/method_chaining.html
//	//
//	//	getTx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
//	//	getTx := db.Where("name = ?", "jinzhu").WithContext(context.Background())
//	//	getTx := db.Where("name = ?", "jinzhu").Debug()
//	//	// `Session`, `WithContext`, `Debug` returns `*gorm.DB` marked as safe to reuse, newly initialized `*gorm.Statement` based on it keeps current conditions
//	//
//	return ds.rootDb.WithContext(ctx)
//}

func OpenDB() error {
	return _ds.Open()
}

func CloseDB() error {
	return _ds.Close()
}

func RootDb() *gorm.DB {
	return _ds.rootDb // to build subqueries subquery
}

func WithContext(ctx context.Context) *gorm.DB {
	return _ds.Session(ctx)
}

func RunTx(ctx context.Context, txFunc func(tx *gorm.DB) error) error {
	return WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return txFunc(tx)
	})
}
