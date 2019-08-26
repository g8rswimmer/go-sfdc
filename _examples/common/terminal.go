package common

import (
	"fmt"
	"os"
)

// TextColor is the color of the ternimal's text.
type TextColor string

const (
	// BlackText will have text color of black
	BlackText TextColor = "30"
	// RedText will have text color of red
	RedText TextColor = "31"
	// GreenText will have text color of green
	GreenText TextColor = "32"
	// YellowText will have text color of yellow
	YellowText TextColor = "33"
	// BlueText will have text color of blue
	BlueText TextColor = "34"
	// PurpleText will have text color of purple
	PurpleText TextColor = "35"
	// CyanText will have text color of cyan
	CyanText TextColor = "36"
	// WhiteText will have text color of white
	WhiteText TextColor = "37"
)

// TextStyle is the style of the terminal's text
type TextStyle string

const (
	// NoEffectStyle will have text with no effects.
	NoEffectStyle TextStyle = "0"
	// BoldStyle will have text with bold effect.
	BoldStyle TextStyle = "1"
	// UnderlineStyle will have text with underline effect.
	UnderlineStyle TextStyle = "4"
	// BlinkStyle will have text with blink effect.
	BlinkStyle TextStyle = "5"
	// InverseStyle will have text with inverse effect.
	InverseStyle TextStyle = "7"
	// HiddenStyle will have text with hidden effect.
	HiddenStyle TextStyle = "8"
)

// BackgroundColor is the color of the terminal's background
type BackgroundColor string

const (
	// BlackBackground will have a black background.
	BlackBackground BackgroundColor = "40"
	// RedBackground will have a red background.
	RedBackground BackgroundColor = "41"
	// GreenBackground will have a green background.
	GreenBackground BackgroundColor = "42"
	// YellowBackground will have a yellow background.
	YellowBackground BackgroundColor = "43"
	// BlueBackground will have a blue background.
	BlueBackground BackgroundColor = "44"
	// PurpleBackground will have a purple background.
	PurpleBackground BackgroundColor = "45"
	// CyanBackground will have a cyan background.
	CyanBackground BackgroundColor = "46"
	// WhiteBackground will have a white background.
	WhiteBackground BackgroundColor = "47"
)

// ANSIColor is the terminal test styling
type ANSIColor struct {
	Color      TextColor
	Style      TextStyle
	Background BackgroundColor
}

// Codes will return the ANSI color codes for the terminal's text.
func (a ANSIColor) Codes() string {
	return fmt.Sprintf("%s;%s;%s", a.Color, a.Style, a.Background)
}

var (
	infoColor = ANSIColor{
		Color:      YellowText,
		Style:      NoEffectStyle,
		Background: BlackBackground,
	}
	warnColor = ANSIColor{
		Color:      CyanText,
		Style:      NoEffectStyle,
		Background: BlackBackground,
	}
	errorColor = ANSIColor{
		Color:      RedText,
		Style:      BoldStyle,
		Background: BlackBackground,
	}
	titleColor = ANSIColor{
		Color:      BlueText,
		Style:      BoldStyle,
		Background: BlackBackground,
	}
)

// Println will print to the terminal with ANSI effects.
func Println(line string, ansi ANSIColor) {
	fmt.Printf("\033[%sm%s\033[0m\n", ansi.Codes(), line)
}

// Title prints the title ANSI effects.
func Title(title string) {
	Println(title, titleColor)
}

// Info prints the information ANSI effects.
func Info(info string, args ...interface{}) {
	Println(fmt.Sprintf(info, args...), infoColor)
}

// Warn prints the warn ANSI effects.
func Warn(warn string, args ...interface{}) {
	Println(fmt.Sprintf(warn, args...), warnColor)
}

// Error prints the error ANSI effects and exits.
func Error(err error) {
	if err != nil {
		Println(fmt.Sprintf("Error: %v", err), errorColor)
		os.Exit(1)
	}
}
