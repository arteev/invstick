package flags

import "fmt"

//ArrayString user set of strings
type ArrayString []string

//String  -  interface Value String() string
func (d *ArrayString) String() string {
	return fmt.Sprintf("%s", *d)
}

//Set interface Value Set(s string) error
func (d *ArrayString) Set(value string) error {
	*d = append(*d, value)
	return nil
}

//Count returns count of strings
func (d *ArrayString) Count() int {
	return len(*d)
}

func (d *ArrayString) Strings() []string {
	return []string(*d)
}
