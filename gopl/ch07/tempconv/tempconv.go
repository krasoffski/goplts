package tempconv

import (
	"flag"
	"fmt"
)

// Celsius represent scale degree Celsius.
type Celsius float64

// Fahrenheit represent scale degree Fahrenheit.
type Fahrenheit float64

// Kelvin represent scale degree Kelvin.
type Kelvin float64

// CToF converts Celsius to Fahrenheit.
func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32.0)
}

// FToC converts Fahrenheit to Celsius.
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32.0) * 5.0 / 9.0)
}

// KToC converts Kelvin to Celsius.
func KToC(k Kelvin) Celsius {
	return Celsius(k - 273.15)
}

type celsiusFlag struct {
	Celsius
}

func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	case "K", "°K":
		f.Celsius = KToC(Kelvin(value))
		return nil

	}
	return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag parses Celsius temperature from command line.
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
