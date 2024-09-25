package model

// BlockState represents the state of the last processed block.
type BlockState struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	LastBlock uint64 `gorm:"not null"` // The last processed block
}

func (m *BlockState) TableName() string {
	return "block_state"
}
