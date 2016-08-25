package domain

import (
	"strconv"
)

type sliceStickers struct {
	genid    int
	stickers []*Sticker
}

func (b *sliceStickers) Create(s *Sticker) {
	b.genid++
	s.ID = strconv.Itoa(b.genid)
	b.stickers = append(b.stickers, s)
}

func (b *sliceStickers) Stickers() (result []*Sticker) {
	return b.stickers
}

func StickersSlice() StickersService {
	return &sliceStickers{}
}
