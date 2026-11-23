package organizer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Classificar(dir string, regras map[string]Regra, logger *slog.Logger) ([]Movimento, error) {
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
			logger.Debug("Arquivo ignorado", "nome", e.Name(), "extensao", ext, "motivo", "extensão não consta nas regras")
			continue
		}
		//logger.Info("Arquivo será movido", "nome", e.Name(), "destino", regra.Destino)
		movimentos = append(movimentos, Movimento{
			Nome:    e.Name(),
			Origem:  filepath.Join(dir, e.Name()),
			Destino: regra.Destino,
		})
	}
	return movimentos, nil
}
