package model

func MapWordsToStringWords(m_words *map[string]string) *string {
	words := new(string)

	for k, v := range *m_words {
		*words += k + " -> " + v + "\n"
	}
	if *words == "" {
		*words += "empty :("
	}
	return words
}
