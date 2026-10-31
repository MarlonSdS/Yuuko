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

func classificar(dir string, regras map[string]organizer.Regra) error {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("lendo pasta %s: %w", dir, err)
	}

	for _, e := range entradas {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		regra, ok := regras[ext]
		if !ok {
			continue
		}
		destino := filepath.Join(dir, regra.Destino)
		err = os.MkdirAll(destino, 0755)
		if err != nil {
			return fmt.Errorf("criando pasta %s: %w", destino, err)
		}
		err = os.Rename(filepath.Join(dir, e.Name()), filepath.Join(destino, e.Name()))
		if err != nil {
			return fmt.Errorf("movendo arquivo %s para %s: %w", e.Name(), destino, err)
		}
	}
	return nil
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
	classificar(downloads, regras)
}
