package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	BaseModel interface {
		GetOne(dest interface{}, options ...Option) error
		GetByCond(dest interface{}, options ...Option) error
		Create(data interface{}) (rowsAffected int64, err error)
		Update(update interface{}, options ...Option) (rowsAffected int64, err error)
		Count(options ...Option) (int64, error)
		Clauses(cond ...clause.Expression) (tx *gorm.DB)
		Upsert(data interface{}, columns []clause.Column, doUpdates clause.Set) (rowsAffected int64, err error)
	}

	defaultBaseModel struct {
		db *gorm.DB
	}
)

func (m *defaultBaseModel) Upsert(data interface{}, columns []clause.Column, doUpdates clause.Set) (rowsAffected int64, err error) {
	if len(columns) == 0 {
		return 0, errors.New("clause.Column must > 1")
	}
	if len(doUpdates) == 0 {
		return 0, errors.New("clause.Set must > 1")
	}

	db := m.Clauses(clause.OnConflict{
		Columns:   columns,
		DoUpdates: doUpdates,
	}).Create(data)

	return db.RowsAffected, db.Error
}

func NewBaseModel(db *gorm.DB) BaseModel {
	return &defaultBaseModel{db}
}

func (m *defaultBaseModel) buildOption(opts ...Option) *gorm.DB {
	for _, op := range opts {
		m.db = op(m.db)
	}
	return m.db
}

func (m *defaultBaseModel) GetOne(dest interface{}, options ...Option) (err error) {
	return m.buildOption(options...).Take(&dest).Error
}

func (m *defaultBaseModel) GetByCond(dest interface{}, options ...Option) (err error) {
	return m.buildOption(options...).Find(dest).Error
}

func (m *defaultBaseModel) Create(data interface{}) (rowsAffected int64, err error) {
	db := m.db.Create(data)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) Update(data interface{}, options ...Option) (rowsAffected int64, err error) {
	db := m.buildOption(options...).Updates(data)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) Count(options ...Option) (count int64, err error) {
	err = m.buildOption(options...).Count(&count).Error
	return
}

func (m *defaultBaseModel) Clauses(cond ...clause.Expression) (tx *gorm.DB) {
	return m.db.Clauses(cond...)
}
