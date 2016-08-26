package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const Version = "0.1"

//Errors
var (
	ErrNoData           = errors.New("Use -gen or -data for data")
	ErrEmptyData        = errors.New("Data is empty")
	ErrTemplateNotFound = errors.New("Templates not found")
	ErrGenCount         = errors.New("Incorrect gen-count")
	ErrCorrectionLevel  = errors.New("Incorrect error correction level for QR Codes")
	ErrEncoding         = errors.New("Incorrect encoding mode for QR Codes")
	ErrSizeBarCode      = errors.New("the dimensions QR Codes are incorrect")
)

//Flags cli
var (
	//List of templates
	Template = flag.String("template", "", "Name of template(s). -template tmpl1 template tmpl2 ...")
	Dir      = flag.String("dir", "", "Output directory")

	Suffix = flag.String("suffix", "", "Optional.Use suffix for generation")
	Prefix = flag.String("prefix", "", "Use prefix for generation")

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
	Width           = flag.Int("width", 100, "Width of each barcode")
	Heigth          = flag.Int("heigth", 100, "Heigth of each barcode")
)

//UserData user set of strings
type ArrayString []string

func (d *ArrayString) String() string {
	return fmt.Sprintf("%s", *d)
}

func (d *ArrayString) Set(value string) error {
	*d = append(*d, value)
	return nil
}

func (d *ArrayString) Count() int {
	return len(*d)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Version %s\n", Version)
	fmt.Fprintf(os.Stderr, "Use flags -data or -gen or stdin pipe for recieve data for QR Codes\n")
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Var(&Data, "data", "List of custom values. -data one -data two ...")
	//flag.Var(&Template, "template", "Name of template(s). -template tmpl1 template tmpl2 ...")
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
	if *Mask == "" {
		*Mask = "%d"
	}
	if *Template == "" {
		return ErrTemplateNotFound
	}
	if *Gen {
		if *GenCount == 0 {
			return ErrGenCount
		}
	}
	if *Width <= 0 || *Heigth <= 0 {
		return ErrSizeBarCode
	}
	return nil
}
