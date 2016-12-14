package main

import (
	"fmt" //
	"html/template"
	"mime/multipart"
	"os"
	"path"
	"regexp"

	"io"

	"strings"

	"bufio"

	"log"
	"net/http"

	"strconv"

	"path/filepath"

	"net/url"

	"github.com/arteev/invstick/barcode"
	"github.com/arteev/invstick/config"
	"github.com/arteev/invstick/domain"
	"github.com/arteev/invstick/flags"
	"github.com/arteev/invstick/helpers"
	"github.com/arteev/invstick/web/model"

	hw "github.com/arteev/invstick/web/view/helpers"
)

type funcMakeData func(stick domain.StickersService) error

//CreateQRCode create QRCode,save to png file. retruns error when failed
func CreateQRCode(stick *domain.Sticker) error {
	filename := path.Join(*flags.Dir, stick.ID+"_sticker.png")
	fout, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fout.Close()
	err = barcode.GenBarcode(fout, stick.Num, *flags.Width, *flags.Height, *flags.CorrectionLevel, *flags.Encoding)
	if err != nil {
		return err
	}
	stick.QRCode = stick.ID + "_sticker.png"
	return nil
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
			if *flags.Barcode {
				if err := CreateQRCode(stick); err != nil {
					return false, err
				}
			}
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
		if *flags.Barcode {
			if err := CreateQRCode(stick); err != nil {
				return err
			}
		}
	}
	return nil
}

//generation
func dataGenerator(c config.Config) []string {
	var r []string
	if c.Mask == "" {
		c.Mask = "%d"
	}
	for i := c.Genstart; i < (c.Genstart + c.Gencount); i++ {
		data := fmt.Sprintf(c.Prefix+c.Mask+c.Suffix, i)
		r = append(r, data)
	}
	return r
}

//TODO:refactor
//generation
func dataGenerate(sticks domain.StickersService) error {

	for i := *flags.GenStart; i < (*flags.GenStart + *flags.GenCount); i++ {
		data := fmt.Sprintf(*flags.Prefix+*flags.Mask+*flags.Suffix, i)
		stick := &domain.Sticker{
			Num: data,
		}
		sticks.Create(stick)
		if *flags.Barcode {
			if err := CreateQRCode(stick); err != nil {
				return err
			}
		}
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
		if err := CreateQRCode(stick); err != nil {
			return err
		}
	}
	return nil
}

func execTemplate(wr io.Writer, sticks domain.StickersService, tmpls ...string) error {
	funcMap := map[string]interface{}{
		"mkSlice":      helpers.MkSlice,
		"mkSliceRange": helpers.MkSliceRange,
		"mul":          helpers.Mul,
		"add":          helpers.Add,
		"calcpages":    helpers.Calcpages,
		"left":         func() int { return *flags.Left },
		"top":          func() int { return *flags.Top }}

	st, err := template.New(path.Base(tmpls[0])).Funcs(funcMap).ParseFiles(tmpls[:]...)
	if err != nil {
		return err
	}

	err = st.Execute(wr, sticks.Stickers())
	return err
}
func doTemplate(sticks domain.StickersService) error {

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
	err := execTemplate(out, sticks, flags.Template.Strings()...)
	return err
}

var tpl *template.Template
var validPath = regexp.MustCompile("^/(index|do|barcode)")

//var validPath = regexp.MustCompile("^/(index|go)/([a-zA-Z0-9]+)$")

func index(w http.ResponseWriter, r *http.Request) {

	if os.Getenv("mode") == "dev" {
		funcMap := map[string]interface{}{
			"translate":  hw.Translate,
			"localename": hw.NameLocale,
		}
		tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/*"))
	}

	data := model.Data
	lang := r.FormValue("lang")
	if lang != "" {
		data.Locale = lang
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func strtointdef(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func rendersticks(w http.ResponseWriter, r *http.Request, c config.Config) error {
	var e error
	switch c.ModeData {
	case config.ModeUserData:
		c.Data = strings.FieldsFunc(r.FormValue("userdata"), func(r rune) bool {
			return r == ' ' || r == ';' || r == ',' || r == '\n'
		})
	case config.ModeGenerate:
		c.Data = dataGenerator(c)
	case config.ModeFile:
		var f multipart.File
		f, _, e = r.FormFile("userfile")
		if e == nil {
			bf := bufio.NewScanner(f)
			for bf.Scan() {
				c.Data = append(c.Data, bf.Text())
			}
			f.Close()
		}
	}

	if len(c.Data) == 0 || e != nil {
		if e != nil {
			return fmt.Errorf("No data found: %q", e)
		}
		return fmt.Errorf("No data found")

	}

	sticks := domain.StickersSlice()

	for _, s := range c.Data {
		stick := &domain.Sticker{
			Num:  s,
			Name: c.Name,
		}
		sticks.Create(stick)
		if c.Barcode {

			//stick.QRCode = "/barcode?val="+stick.Num+
			u := &url.URL{}
			u.Scheme = r.URL.Scheme
			u.Host = r.URL.Host
			u.Path = "/barcode"
			q := u.Query()
			q.Set("val", stick.Num)
			q.Set("w", strconv.Itoa(c.Width))
			q.Set("h", strconv.Itoa(c.Height))
			q.Set("c", c.Correctionlevel)
			q.Set("e", c.Encoding)
			u.RawQuery = q.Encode()
			stick.QRCode = u.String()
		}
	}
	err := execTemplate(w, sticks, c.Template)
	return err
}

func getbarcode(w http.ResponseWriter, r *http.Request) {
	//  /barcode?c=L&e=&h=45&val=1&w=45
	width := strtointdef(r.FormValue("w"), 100)
	height := strtointdef(r.FormValue("h"), 100)
	err := barcode.GenBarcode(w, r.FormValue("val"), width, height, r.FormValue("e"), r.FormValue("c"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func do(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c := &config.Config{}
		//TODO: folder stickers from config
		c.Name = r.FormValue("name")
		c.Template = filepath.Join("stickers", r.FormValue("template"))
		c.Prefix = r.FormValue("prefix")
		c.Suffix = r.FormValue("suffix")
		//BarCode
		c.Barcode = r.FormValue("barcode") == "on"
		c.Correctionlevel = r.FormValue("level")
		c.Encoding = r.FormValue("encoding")
		c.Width = strtointdef(r.FormValue("width"), 100)
		c.Height = strtointdef(r.FormValue("height"), 100)
		c.Top = strtointdef(r.FormValue("top"), 1)
		c.Left = strtointdef(r.FormValue("left"), 1)
		switch r.FormValue("datain") {
		case "gen":
			c.ModeData = config.ModeGenerate
		case "file":
			c.ModeData = config.ModeFile
		case "data":
			c.ModeData = config.ModeUserData
		default:
			c.ModeData = config.ModeGenerate
		}
		c.Genstart = strtointdef(r.FormValue("gen-start"), 1)
		c.Gencount = strtointdef(r.FormValue("gen-count"), 10)
		c.Mask = r.FormValue("gen-mask")
		//TODO: CheckValues
		err := rendersticks(w, r, *c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func makeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if r.URL.Path != "/" && m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}
func starthttp() {
	//TODO: name folder with stickers from config
	hw.LoadLocales("config/locales")
	templates, err := helpers.GetStickerTemplates("stickers", ".gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	if len(templates) == 0 {
		templates = append(templates, "(not found)")
	}
	model.Data.Templates = templates
	model.Data.Locales = hw.Locales()
	model.Data.Locale = "en"
	funcMap := map[string]interface{}{
		"translate":  hw.Translate,
		"localename": hw.NameLocale,
	}
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/*"))
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", makeHandler(index))
	http.HandleFunc("/do", makeHandler(do))
	http.HandleFunc("/barcode", makeHandler(getbarcode))
	log.Fatalln(http.ListenAndServe(*flags.WebAddr, nil))
}

func main() {

	var fmake funcMakeData

	flags.Parse()

	if *flags.WebAddr != "" {
		starthttp()
		return
	}

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
	if err := doTemplate(sticks); err != nil {
		flags.ExitWithError(err)
	}
}
