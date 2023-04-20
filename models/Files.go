package models

type File struct {
	ID   uint   `gorm:"not null;primaryKey"`
	Name string `gorm:"not null;type:varchar(100)"`
	File []byte `gorm:"type:bytea"`
}
