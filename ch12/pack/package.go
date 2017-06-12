package pack

import "fmt"

type Test struct {
	Es
	us
	Big   int
	small int
}

type Es struct {
	F float32
	f float32
}

type us struct {
	B bool
	b bool
}

func (Es) Efunc() {
	fmt.Println("Efunc")
}
func (Es) efunc() {
	fmt.Println("efunc")
}
func (*us) Ufunc() {
	fmt.Println("Ufunc")
}
func (*us) ufunc() {
	fmt.Println("ufunc")
}

func (t Test) Tfunc() {
	fmt.Println("Tfunc")
	t.tfunc()
	t.Ufunc()
	t.ufunc()
	t.Efunc()
	t.efunc()
}

func (t Test) tfunc() {
	fmt.Println("tfunc")

}
