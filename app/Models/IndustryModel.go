package Models

var IndustryModel *industrymodel

type industrymodel struct {
	BaseModel
}

type Wx_industry struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"column:user_id;not null"`
	Name      string    `gorm:"column:name;not null"`
	CreatedAt LocalTime `gorm:"column:created_at"`
	UpdatedAt LocalTime `gorm:"column:updated_at"`
	DeletedAt LocalTime `gorm:"column:deleted_at;index"`
}
