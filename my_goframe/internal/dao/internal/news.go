// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NewsDao is the data access object for the table news.
type NewsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  NewsColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// NewsColumns defines and stores column names for the table news.
type NewsColumns struct {
	Id         string // 新闻id
	UserId     string // 发布者id
	Title      string // 新闻标题
	Article    string // 正文
	CreateTime string // 发布时间
	UpdateTime string // 最后修改时间
	DeleteTime string // 删除时间
}

// newsColumns holds the columns for the table news.
var newsColumns = NewsColumns{
	Id:         "id",
	UserId:     "user_id",
	Title:      "title",
	Article:    "article",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	DeleteTime: "delete_time",
}

// NewNewsDao creates and returns a new DAO object for table data access.
func NewNewsDao(handlers ...gdb.ModelHandler) *NewsDao {
	return &NewsDao{
		group:    "default",
		table:    "news",
		columns:  newsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NewsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NewsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NewsDao) Columns() NewsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NewsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NewsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *NewsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
