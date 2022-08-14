package model

import "time"

// @Author: Derek
// @Description: Base model.
// @Date: 2022/5/13 14:47
// @Version 1.0

type Model struct {
	Id             uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"  json:"id"`                            // 主键Id
	CreateTime     int64     `gorm:"column:create_time;default:0;not null"  json:"createTime"`                          // 创建时间
	LastModifyTime time.Time `gorm:"column:last_modify_time;default:CURRENT_TIMESTAMP;not null"  json:"lastModifyTime"` // 最后修改时间
	Deleted        uint8     `gorm:"column:deleted;default:0;not null"  json:"deleted"`                                 // 0:正常 1:删除
}
