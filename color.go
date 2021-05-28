package main

import "fmt"
const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}


func HBlack(str string) string {
	return textColorhighlight(textBlack, str)
}

func HRed(str string) string {
	return textColorhighlight(textRed, str)
}
func HYellow(str string) string {
	return textColorhighlight(textYellow, str)
}
func HGreen(str string) string {
	return textColorhighlight(textGreen, str)
}
func HCyan(str string) string {
	return textColorhighlight(textCyan, str)
}
func HBlue(str string) string {
	return textColorhighlight(textBlue, str)
}
func HPurple(str string) string {
	return textColorhighlight(textPurple, str)
}
func HWhite(str string) string {
	return textColorhighlight(textWhite, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}

func textColorhighlight(color int, str string) string {
	return fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", color, str)
}
