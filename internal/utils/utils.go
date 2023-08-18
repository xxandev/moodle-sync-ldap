package utils

import (
	"os"
	"regexp"
	"strings"
)

var tmap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "h", 'д': "d", 'е': "e", 'є': "ie", 'ж': "zh", 'з': "z",
	'и': "i", 'і': "i", 'ї': "i", 'й': "i", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "shch", 'ь': "", 'ю': "iu", 'я': "ia", 'ё': "jo", 'ы': "y", 'ъ': "",
}

func CreateFile(name string, content []byte) error {
	// if _, err := os.Stat(name); err == nil {
	// 	return fmt.Errorf("such file exists %s", name)
	// }
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return file.Close()
}

func IsStatFile(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

func Transliterate(text string) (res string) {
	text = strings.ToLower(text)
	for _, char := range text {
		if replacement, ok := tmap[char]; ok {
			res += replacement
		} else {
			res += string(char)
		}
	}
	return regexp.MustCompile("[^a-zA-Z0-9а-яА-Я]+").ReplaceAllString(res, "")
}
