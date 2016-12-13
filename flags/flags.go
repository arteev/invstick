package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

//Version application
const Version = "0.1"

//Errors
var (
	ErrNoData           = errors.New("Use -gen or -data for data")
	ErrEmptyData        = errors.New("Data is empty")
	ErrTemplateNotFound = errors.New("Templates not found")
	ErrGenCount         = errors.New("Incorrect gen-count")
	ErrCorrectionLevel  = errors.New("Incorrect error correction level for QR Codes")
	ErrEncoding         = errors.New("Incorrect encoding mode for QR Codes")
	ErrSizeBarCode      = errors.New("The dimensions QR Codes are incorrect")
	ErrPosition         = errors.New("The position are incorrect")
)

//Flags cli
var (
	//List of templates

	Dir      = flag.String("dir", "", "Output directory")
	Template ArrayString
	Suffix   = flag.String("suffix", "", "Optional.Use suffix for generation")
	Prefix   = flag.String("prefix", "", "Use prefix for generation")

	//generation data
	Gen      = flag.Bool("gen", false, "Use for generation data")
	GenStart = flag.Int("gen-start", 1, "Start numbers used in data generation")
	GenCount = flag.Int("gen-count", 10, "Number of data generation")
	Mask     = flag.String("mask", "", "Use mask for generation: sample: %06d")

	//User data
	Data ArrayString

	//QRCodes flags
	CorrectionLevel = flag.String("correction", "M", "Error Correction Level: L recovers 7% of data;M recovers 15% of data;Q recovers 25% of data:H recovers 30% of data")
	Encoding        = flag.String("encoding", "Auto", "Encoding mode for QR Codes. Auto,Numeric,AlphaNumeric,Unicode")
	Width           = flag.Int("width", 100, "Width of barcode")
	Height          = flag.Int("height", 100, "Height of barcode")
	Left            = flag.Int("left", 1, "Start horizontal position on sheet")
	Top             = flag.Int("top", 1, "Start vertical position on sheet")
	Barcode         = flag.Bool("barcode", true, "Generate QR Codes")

	//web
	WebAddr = flag.String("listen", "", "Use web interface on :80 or other")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Version %s\n", Version)
	fmt.Fprintf(os.Stderr, "Use flags -data or -gen or stdin pipe for recieve data for QR Codes\n")
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	//Template = flag.String("template", "", "Name of template(s). -template tmpl1 template tmpl2 ...")
	flag.Var(&Template, "template", "Name of template(s). -template tmpl1 template tmpl2 ...")
	flag.Var(&Data, "data", "List of custom values. -data one -data two ...")
	//flag.Var(&Template, "template", "Name of template(s). -template tmpl1 template tmpl2 ...")

}

func Parse() {
	flag.Parse()
	if err := check(); err != nil {
		ExitWithError(err)
	}
}

func ExitWithError(e error) {
	fmt.Fprintln(os.Stderr, e)
	flag.Usage()
	os.Exit(2)
}

func check() error {
	if *WebAddr != "" {
		return nil
	}

	if *Mask == "" {
		*Mask = "%d"
	}
	if Template.Count() == 0 {
		return ErrTemplateNotFound
	}
	if *Gen {
		if *GenCount == 0 {
			return ErrGenCount
		}
	}
	if *Width <= 0 || *Height <= 0 {
		return ErrSizeBarCode
	}
	if *Left <= 0 || *Top <= 0 {
		return ErrPosition
	}
	return nil
}
