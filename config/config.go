package config

type Mode int

const (
	ModeGenerate = Mode(0)
	ModeUserData = Mode(1)
	ModeFile     = Mode(2)
)

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
