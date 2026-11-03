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
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
	fmt.Println(c)
	regras := organizer.PegarRegras()
	movimentos, err := organizer.Classificar(downloads, regras)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao classificar arquivos:", err)
		return
	}

	err = organizer.MoverArquivos(movimentos)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao mover arquivos:", err)
		return
	}
	fmt.Println("Arquivos movidos com sucesso!")
}
