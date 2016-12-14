package domain

import "testing"

func TestStickers(t *testing.T) {
	stickers := StickersSlice()
	stickersSlice, ok := stickers.(*sliceStickers)
	if !ok {
		t.Fatal("Excepted sliceStickers ")
	}

	s := &Sticker{
		Num:    "1",
		QRCode: "code",
	}
	stickers.Create(s)
	if s.ID != "1" {
		t.Errorf("Excepted sticker.ID=1, got %s\n", s.ID)
	}

	if stickersSlice.genid != 1 {
		t.Errorf("Excepted stikersSlice.genid=1, got %d\n", stickersSlice.genid)
	}

	sarr := stickers.Stickers()
	if len(sarr) != 1 || sarr[0] != s {
		t.Errorf("Excepted Stickers() %#v,got %#v", s, sarr[0])
	}

}
