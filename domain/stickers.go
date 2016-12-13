package domain

type Sticker struct {
	Name   string
	ID     string
	Num    string
	QRCode string
}

type StickersService interface {
	Create(s *Sticker)
	Stickers() []*Sticker
}
