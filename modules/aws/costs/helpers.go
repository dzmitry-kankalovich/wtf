package awscosts

import (
	"strconv"
	"strings"
)

func mapUnitToChar(unit string) string {
	if unit == "USD" {
		return "$"
	}
	return unit
}

func trunkateAmount(amount string) string {
	floatAmount, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		return amount
	}
	parts := strings.Split(strconv.FormatFloat(floatAmount, 'f', 2, 32), ".")

	// Add commas to the integer part
	integerPart := parts[0]
	var withCommas strings.Builder
	for i, digit := range []rune(integerPart) {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			withCommas.WriteString(",")
		}
		withCommas.WriteRune(digit)
	}

	// Append the decimal part, if any
	if len(parts) > 1 {
		withCommas.WriteString(".")
		withCommas.WriteString(parts[1])
	}

	return withCommas.String()
}
