package task01

import (
	"math"
	"testing"
)

// TODO: напиши табличные тесты для функции ParsePrice.
//
// Создай срез структур с полями:
//   name        string   - название теста
//   input       string   - входная строка
//   expected    float64  - ожидаемый результат
//   expectError bool     - ожидаем ли ошибку
//
// Обязательно проверь эти случаи:
//   Корректные: "1500", "1 500", "1500.50", "1500,50", "1 500,50 руб.", "₽1 500"
//   Граничные:  "0", "0,01"
//   Некорректные: "", "abc", "руб.", "---"
//
// Подсказка: запусти go test -v ./... чтобы увидеть подробный вывод

func TestParsePrice(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		// TODO: заполни таблицу тест-кейсов

		{"обычное целое число", "1500", 1500.00, false},
		{"пробел как разделитель тысяч", "1 500", 1500.00, false},
		{"точка как десятичный разделитель", "1500.50", 1500.50, false},
		{"запятая как десятичный разделитель", "1500,50", 1500.50, false},
		{"цена с суффиксом руб.", "1 500,50 руб.", 1500.50, false},
		{"цена с символом рубля", "₽1 500", 1500.00, false},

		{"ноль", "0", 0.00, false},
		{"минимальная дробная цена", "0,01", 0.01, false},

		{"пустая строка", "", 0, true},
		{"буквы", "abc", 0, true},
		{"только руб.", "руб.", 0, true},
		{"дефисы", "---", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: вызови ParsePrice(tc.input)
			got, err := ParsePrice(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("ожидали ошибку, но получили nil")
				}
				return
			}

			if err != nil {
				t.Errorf("не ожидали ошибку, но получили: %v", err)
				return
			}

			if math.Abs(got-tc.expected) > 0.001 {
				t.Errorf("ParsePrice(%q) = %f, ожидали %f",
					tc.input,
					got,
					tc.expected)
			}

		})
	}
}

// FuzzParsePrice проверяет что функция не паникует ни при каких входных данных.
// Запуск: go test -fuzz=FuzzParsePrice -fuzztime=10s
func FuzzParsePrice(f *testing.F) {
	// TODO: добавь начальные случаи через f.Add(...)
	// f.Add("1500")
	// f.Add("")

	f.Add("1500")
	f.Add("1 500,50 руб.")
	f.Add("")
	f.Add("abc")

	f.Fuzz(func(t *testing.T, input string) {
		// TODO: вызови ParsePrice(input)
		// Функция не должна паниковать - только возвращать ошибку
		// Не нужно проверять результат - главное что нет panic
		_, _ = ParsePrice(input)
	})
}
