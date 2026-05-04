// ============================================================
// Задача: Concurrent Web Crawler  🔴 Senior
// ============================================================
//
// Классика (LeetCode 1242). Дан граф ссылок: URL → список URL-соседей.
// Нужно посетить все достижимые URL-ы (BFS/DFS) используя не более N
// горутин одновременно.
//
// Интерфейс:
//
//   type Fetcher interface {
//       Fetch(url string) (links []string, err error)
//   }
//
//   func Crawl(start string, f Fetcher, parallelism int) []string
//
// Требования:
//   - Каждый URL обрабатывается РОВНО ОДИН РАЗ (дедупликация)
//   - Параллельно обрабатываются максимум parallelism URL-ов
//   - Возврат: отсортированный список всех посещённых URL-ов
//   - Нет утечек горутин, нет дедлока при любом графе (включая циклы!)
//
// Подвох №1: очевидное "обходим рекурсивно с sync.WaitGroup" даёт дедлок
// если ограничение parallelism меньше глубины графа и мы блокируем
// ожидание внутри горутины которая сама держит семафор.
//
// Подвох №2: как понять что "всё обработано" при BFS через канал-очередь?
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sort"
	"time"
)

type Fetcher interface {
	Fetch(url string) ([]string, error)
}

// TODO: реализуй Crawl
// Подсказка 1: используй sync.Map или map+mutex для "уже посещён"
// Подсказка 2: для ограничения parallelism — семафор-канал с буфером N.
// Но НЕ захватывай семафор внутри горутины, которая ждёт других горутин!
// Разделяй "обработать одну страницу" и "управлять пулом".
// Подсказка 3: завершение определяй по sync.WaitGroup + горутина-закрывашка.
func Crawl(start string, f Fetcher, parallelism int) []string {
	// TODO
	return nil
}

// === Mock fetcher для main ===

type mockFetcher map[string][]string

func (m mockFetcher) Fetch(url string) ([]string, error) {
	time.Sleep(10 * time.Millisecond)
	if links, ok := m[url]; ok {
		return links, nil
	}
	return nil, fmt.Errorf("not found: %s", url)
}

func main() {
	// Граф с циклом: / → /a, /b; /a → /a1, /; /b → /; /a1 → /
	g := mockFetcher{
		"/":   {"/a", "/b"},
		"/a":  {"/a1", "/"},
		"/b":  {"/"},
		"/a1": {"/"},
	}

	visited := Crawl("/", g, 3)
	sort.Strings(visited)
	fmt.Println(visited) // [/ /a /a1 /b]
}
