// Задача 4: Оборачивание ошибки
//
// Ожидаемый вывод:
//   load data: read file: file not found

package main

import (
"fmt"
"errors"
)

// TODO: напиши функцию readFileMock() error
//       возвращай errors.New("file not found")
func readFileMock() error {
	return errors.New("file not found")
}
// TODO: напиши функцию loadData() error
//       вызывает readFileMock() и оборачивает ошибку через fmt.Errorf с %w:
//       "read file: <err>"
//       сама тоже оборачивает: "load data: <err>"

func loadData() error {
	err := readFileMock()
	if err != nil {
		err = fmt.Errorf("read file: %w", err)
		return fmt.Errorf("load data: %w", err)
	}
	return nil
}

func main() {
	err := loadData()
	fmt.Println(err)
}
