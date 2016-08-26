package main

import (
	"fmt" //
	"html/template"
	"image/png"
	"os"
	"path"

	"io"

	"strings"

	"bufio"

	"math"

	"github.com/arteev/invstick/domain"
	"github.com/arteev/invstick/flags"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type FuncMakeData func(stick domain.StickersService) error

func StringECL(level string) qr.ErrorCorrectionLevel {
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

func Encoding(value string) qr.Encoding {
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

//CreateQRCode returns create QRCode,save to png file and retruns url
func CreateQRCode(stick *domain.Sticker) {
	filename := path.Join(*flags.Dir, stick.ID+"_sticker.png")
	fout, _ := os.Create(filename)
	defer fout.Close()
	qrcode, err := qr.Encode(stick.Num, StringECL(*flags.CorrectionLevel), Encoding(*flags.Encoding))
	if err != nil {
		panic(err)
	}
	qrcode, err = barcode.Scale(qrcode, *flags.Width, *flags.Heigth)
	if err != nil {
		panic(err)
	}
	err = png.Encode(fout, qrcode)
	if err != nil {
		panic(err)
	}
	stick.QRCode = stick.ID + "_sticker.png"

}

func dataFromFile(filename string, sticks domain.StickersService) (bool, error) {
	exists := false
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		f, err := os.Open(filename)
		if err != nil {
			return false, err
		}
		bf := bufio.NewScanner(f)
		for bf.Scan() {
			exists = true
			text := bf.Text()
			stick := &domain.Sticker{
				Num: text,
			}
			sticks.Create(stick)
			CreateQRCode(stick)
		}
		f.Close()

	}
	return exists, nil
}

//read from args
func dataReadArgs(sticks domain.StickersService) error {

	for _, d := range flags.Data {
		if strings.HasPrefix(d, "@") {
			//read from file
			ok, err := dataFromFile(d[1:], sticks)
			if err != nil {
				return err
			}
			if ok {
				continue
			}
		}
		stick := &domain.Sticker{
			Num: d,
		}
		sticks.Create(stick)
		CreateQRCode(stick)
	}
	return nil
}

//generation
func dataGenerate(sticks domain.StickersService) error {

	for i := *flags.GenStart; i < (*flags.GenStart + *flags.GenCount); i++ {
		data := fmt.Sprintf(*flags.Prefix+*flags.Mask+*flags.Suffix, i)
		stick := &domain.Sticker{
			Num: data,
		}
		sticks.Create(stick)
		CreateQRCode(stick)
	}
	return nil
}

func dataFromPipe(sticks domain.StickersService) error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return nil
	}

	r := bufio.NewReader(os.Stdin)
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		s = strings.TrimRight(s, "\r\n")
		stick := &domain.Sticker{
			Num: s,
		}
		sticks.Create(stick)
		CreateQRCode(stick)
	}
	return nil
}

func mkSlice(args ...interface{}) []interface{} {
	return args
}

func mul(a, b int) int {
	return a * b
}

func add(a, b int) int {
	return a + b
}

func calcpages(onpage, count int) int {
	i := count / onpage
	if math.Mod(float64(count), float64(onpage)) > 0 {
		i++
	}
	return i
}

func mkSliceRange(from, count int) (result []int) {
	for i := from; i < from+count; i++ {
		result = append(result, i)
	}
	return
}

func DoTemplate(sticks domain.StickersService) error {
	funcMap := map[string]interface{}{
		"mkSlice":      mkSlice,
		"mkSliceRange": mkSliceRange,
		"mul":          mul,
		"add":          add,
		"calcpages":    calcpages}

	st, err := template.New(path.Base(*flags.Template)).Funcs(funcMap).ParseFiles(*flags.Template)
	if err != nil {
		return err
	}

	var out io.Writer
	if *flags.Dir == "" {
		out = os.Stdout
	} else {
		fout, err := os.Create(path.Join(*flags.Dir, "index.html"))
		if err != nil {
			return err
		}
		defer fout.Close()
		out = fout
	}
	err = st.Execute(out, sticks.Stickers())
	if err != nil {
		return err
	}
	return nil
}

func main() {

	var fmake FuncMakeData

	if !*flags.Gen && flags.Data.Count() == 0 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			flags.ExitWithError(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			flags.ExitWithError(flags.ErrNoData)
		}
		//from pipe
		fmake = dataFromPipe
	} else if *flags.Gen {
		//generate
		fmake = dataGenerate
	} else {
		//read from args
		fmake = dataReadArgs
	}

	if fmake == nil {
		flags.ExitWithError(flags.ErrNoData)
	}
	//Make stickers
	sticks := domain.StickersSlice()
	err := fmake(sticks)
	if err != nil {
		flags.ExitWithError(err)
	}
	if err := DoTemplate(sticks); err != nil {
		flags.ExitWithError(err)
	}
}
