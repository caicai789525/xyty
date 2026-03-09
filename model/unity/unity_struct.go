package unity

type UnityStruct struct {
	ID     uint   `gorm:"primaryKey"`
	Model  string `gorm:"column:model"`
	Struct string `gorm:"column:struct"`
}
