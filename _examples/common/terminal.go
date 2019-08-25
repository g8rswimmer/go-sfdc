package common

import (
	"fmt"
	"os"
)

type TextColor string

const (
	BlackText  TextColor = "30"
	RedText    TextColor = "31"
	GreenText  TextColor = "32"
	YellowText TextColor = "33"
	BlueText   TextColor = "34"
	PurpleText TextColor = "35"
	CyanText   TextColor = "36"
	WhiteText  TextColor = "37"
)

type TextStyle string

const (
	NoEffectStyle  TextStyle = "0"
	BoldStyle      TextStyle = "1"
	UnderlineStyle TextStyle = "4"
	BlinkStyle     TextStyle = "5"
	InverseStyle   TextStyle = "7"
	HiddenStyle    TextStyle = "8"
)

type BackgroundColor string

const (
	BlackBackground  BackgroundColor = "40"
	RedBackground    BackgroundColor = "41"
	GreenBackground  BackgroundColor = "42"
	YellowBackground BackgroundColor = "43"
	BlueBackground   BackgroundColor = "44"
	PurpleBackground BackgroundColor = "45"
	CyanBackground   BackgroundColor = "46"
	WhiteBackground  BackgroundColor = "47"
)

type ANSIColor struct {
	Color      TextColor
	Style      TextStyle
	Background BackgroundColor
}

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

func Println(line string, ansi ANSIColor) {
	fmt.Printf("\033[%sm%s\033[0m\n", ansi.Codes(), line)
}

func Title(title string) {
	Println(title, titleColor)
}

func Info(info string, args ...interface{}) {
	Println(fmt.Sprintf(info, args...), infoColor)
}

func Warn(warn string, args ...interface{}) {
	Println(fmt.Sprintf(warn, args...), warnColor)
}

func Error(err error) {
	if err != nil {
		Println(fmt.Sprintf("Error: %v", err), errorColor)
		os.Exit(1)
	}
}
