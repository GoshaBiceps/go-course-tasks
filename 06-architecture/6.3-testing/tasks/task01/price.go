package task01

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParsePrice парсит строку с ценой в число.
// Поддерживает форматы:
//
//	"1500"        -> 1500.00
//	"1 500"       -> 1500.00  (пробел как разделитель тысяч)
//	"1500.50"     -> 1500.50
//	"1500,50"     -> 1500.50  (запятая как десятичный разделитель)
//	"1 500,50"    -> 1500.50
//	"1 500,50 руб." -> 1500.50 (с суффиксом)
//	"₽1 500"      -> 1500.00  (с символом валюты)
func ParsePrice(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("пустая строка")
	}

	// Убираем символы валют и слова
	s = removeNonNumeric(s)
	if s == "" {
		return 0, fmt.Errorf("строка не содержит числа")
	}

	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, ".,")

	// Убираем пробелы (разделители тысяч)
	s = strings.ReplaceAll(s, " ", "")

	// Заменяем запятую на точку
	s = strings.ReplaceAll(s, ",", ".")

	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось распознать число: %w", err)
	}

	return result, nil
}

// removeNonNumeric убирает нечисловые символы кроме пробелов, точек и запятых.
func removeNonNumeric(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) || r == ' ' || r == '.' || r == ',' {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}
