package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Classificar(dir string, regras map[string]Regra) ([]Movimento, error) {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("lendo pasta %s: %w", dir, err)
	}

	var movimentos []Movimento
	for _, e := range entradas {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		regra, ok := regras[ext]
		if !ok {
			continue
		}
		fmt.Printf("Arquivo %s será movido para %s \n", e.Name(), regra.Destino)
		movimentos = append(movimentos, Movimento{
			Nome:    e.Name(),
			Origem:  filepath.Join(dir, e.Name()),
			Destino: regra.Destino,
		})
	}
	return movimentos, nil
}
