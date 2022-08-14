package model

import (
	"fmt"
	"gorm.io/gorm"
)

// @Author: Derek
// @Description: DB Option.
// @Date: 2022/8/14 09:03
// @Version 1.0

type Option func(*gorm.DB) *gorm.DB

func Select(fields string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields, args...)
	}
}

func DistinctSelect(fields string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fmt.Sprintf("distinct %s", fields), args...)
	}
}

func Where(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func Limit(limit int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Preload(query string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func Order(order string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func WithId(id uint64) Option {
	return Where("id = ?", id)
}

func WithIds(ids []uint64) Option {
	return Where("id in (?)", ids)
}

func Delete() Option {
	return Where("deleted = ?", DeletedYes)
}

func UnDelete() Option {
	return Where("deleted = ?", DeletedNo)
}
