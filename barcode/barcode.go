package barcode

import (
	"image/png"
	"io"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func strToCorrectionLevel(level string) qr.ErrorCorrectionLevel {
	switch level {
	case "L":
		return qr.L
	case "M":
		return qr.M
	case "Q":
		return qr.Q
	case "H":
		return qr.H
	}
	return qr.M
}

func strToEncoding(value string) qr.Encoding {
	switch value {
	case "Auto":
		return qr.Auto
	case "Numeric":
		return qr.Numeric
	case "AlphaNumeric":
		return qr.AlphaNumeric
	case "Unicode":
		return qr.Unicode
	}
	return qr.Auto
}

// GenBarcode генерирует qr-code размером w*h в кодировке encoding и уровнем коррекции ошибок level
// С содержимым content
// Пишет результат в wr oi.Writer
// Возвращает ошибку при неудаче
func GenBarcode(wr io.Writer, content string, width, height int, encoding, level string) error {
	qrcode, err := qr.Encode(content, strToCorrectionLevel(level), strToEncoding(encoding))
	if err != nil {
		return err
	}
	qrcode, err = barcode.Scale(qrcode, width, height)
	if err != nil {
		return err
	}
	err = png.Encode(wr, qrcode)
	return err
}
