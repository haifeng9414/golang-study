package main

type Error string

func (e Error) Error() string {
	return string(e)
}

type Demo struct {
}

func (d *Demo) do() *Demo {
	panic(Error("test"))
}

func (d *Demo) error(err string) {
	panic(Error(err))
}
func Compile(str string) (regexp *Demo, err error) {
	regexp = new(Demo)
	defer func() {
		if e := recover(); e != nil {
			regexp = nil
			err = e.(Error)
		}
	}()
	return regexp.do(), nil
}

func main() {
	Compile("d")
}