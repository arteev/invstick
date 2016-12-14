package config

//A Mode obtaining data for stickers
type Mode int

//Mode receiving data
const (
	ModeGenerate = Mode(0)
	ModeUserData = Mode(1)
	ModeFile     = Mode(2)
)

//Config settings for rendering stickers
type Config struct {
	Template        string
	Name            string
	Prefix          string
	Suffix          string
	ModeData        Mode
	Genstart        int
	Gencount        int
	Mask            string
	Data            []string // for load from file, from text(user input)
	Barcode         bool
	Correctionlevel string
	Encoding        string
	Width           int
	Height          int
	Left            int
	Top             int
}
