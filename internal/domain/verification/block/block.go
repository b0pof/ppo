package block

import "time"

const blockPeriod = time.Hour

type Block struct{}

func New() *Block {
	return &Block{}
}

func (b *Block) GetExpiration() time.Time {
	return time.Now().Add(blockPeriod)
}
