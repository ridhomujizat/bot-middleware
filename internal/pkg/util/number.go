package util

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateId() (string, error) {
	const urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"

	id, err := gonanoid.Generate(urlAlphabet, 10)
	if err != nil {
		return "", err
	}
	return id, nil
}
