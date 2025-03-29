package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
)

var wordNormalizeRegex = regexp.MustCompile(`[^A-Za-zА-Яа-я]`)

type wordEntry struct {
	word  string
	count int
}

func main() {
	var resultCount int
	flag.IntVar(&resultCount, "n", 10, "Сколько наиболее используемых слов отображать")
	flag.Parse()
	fileName := flag.Arg(0)

	if fileName == "" {
		fmt.Println("Не указано имя файла для открытия")
		return
	}

	wordsUsageRaw, err := readWordsFromFile(fileName)
	if err != nil {
		slog.Error("Ошибка при открытии файла", slog.String("ошибка", err.Error()))
		return
	}

	sortedWordEntries := sortWordEntries(wordsUsageRaw)

	if resultCount > len(sortedWordEntries) {
		resultCount = len(sortedWordEntries)
	}

	for i := 0; i < resultCount; i++ {
		fmt.Printf("%v: %d\n", sortedWordEntries[i].word, sortedWordEntries[i].count)
	}
}

func readWordsFromFile(fileName string) (map[string]int, error) {
	wordsUsage := make(map[string]int)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("не получилось открыть файл: %s. Ошибка: %w", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		if word == "" {
			continue
		}

		wordsUsage[word] += 1
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return wordsUsage, nil
}

func sortWordEntries(wordsUsage map[string]int) []wordEntry {
	sliceForSort := make([]wordEntry, 0, len(wordsUsage))

	for word, count := range wordsUsage {
		sliceForSort = append(sliceForSort, wordEntry{word: word, count: count})
	}

	sort.Slice(sliceForSort, func(i, j int) bool {
		return sliceForSort[i].count > sliceForSort[j].count
	})

	return sliceForSort
}

func normalizeWord(word string) string {
	word = wordNormalizeRegex.ReplaceAllString(word, "")
	word = strings.ToLower(word)
	return word
}
