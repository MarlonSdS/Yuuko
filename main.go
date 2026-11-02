package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarlonSdS/Yuuko/organizer"
)

func contarPorExtensao(dir string) (map[string]int, error) {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("lendo pasta %s: %w", dir, err)
	}

	contagem := make(map[string]int)
	for _, e := range entradas {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		contagem[ext]++
	}
	return contagem, nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Erro ao localizar a Home: %v\n", err)
		return
	}

	downloads := filepath.Join(home, "Downloads")
	c, err := contarPorExtensao(downloads)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(c)
	regras := organizer.PegarRegras()
	//fmt.Println(regras)
	organizer.Classificar(downloads, regras)
}
