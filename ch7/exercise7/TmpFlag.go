package exercise7

import(
	"fmt"
	"flag"
	"gopl.io/ch7/tempconv"
)

// type temp {
// 	tempconv.Celsius
// 	kelvin float64
// }

type Kelvin float64
const KelvinZero = -273.0
type tempFlag struct {
	tempconv.Celsius
}

func KToC(c Kelvin) tempconv.Celsius{
	return tempconv.Celsius(c - KelvinZero)
}

func (f *tempFlag) Set(s string) error  {
	var t float64
	var unit string
	fmt.Sscanf(s, "%f%s", &t, &unit) //error can be dealed in switch default case
	switch unit {
	case "C", "c", "\u00b0C":
		f.Celsius = tempconv.Celsius(t)
	case "F", "f":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(t))
	case "K", "k":
		f.Celsius = KToC(Kelvin(t))
	default:
		return fmt.Errorf("invalid temperature %v", t)
	}
	return nil
}

func TempFlag(option string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	var f tempFlag
	f.Celsius = value
	flag.CommandLine.Var(&f, option, usage)
	return &f.Celsius
}
