package dao

import (
	"context"
	"github.com/jinzhu/gorm"
)

type BaseDao struct {
	ctx    context.Context
	master *gorm.DB
}

func (d *BaseDao) Active(ctx context.Context) {
	d.ctx = ctx
	d.master = ctx.Value("master").(*gorm.DB)
}
