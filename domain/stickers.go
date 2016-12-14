package domain

//A Sticker storing data of sticker
type Sticker struct {
	Name   string
	ID     string
	Num    string
	QRCode string
}

//StickersService stickers collection
type StickersService interface {
	Create(s *Sticker)
	Stickers() []*Sticker
}
