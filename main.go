package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarlonSdS/Yuuko/logging"
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
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
}

func run() error {
	arquivo, err := logging.CriarArquivoLog()
	if err != nil {
		return fmt.Errorf("criando arquivo de logs: %w", err)
	}
	fmt.Println("Guardando logs em:", arquivo.Name())
	defer arquivo.Close()

	logger := logging.NovoLogger(io.MultiWriter(os.Stdout, arquivo))

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Erro ao localizar a Home", "erro", err)
		return fmt.Errorf("localizando a Home: %w", err)
	}

	downloads := filepath.Join(home, "Downloads")
	logger.Info("Iniciando a classificação de arquivos", "diretorio", downloads)

	c, err := contarPorExtensao(downloads)
	if err != nil {
		logger.Error("Erro ao contar arquivos", "erro", err)
		return fmt.Errorf("contando arquivos por extensão: %w", err)
	}
	logger.Info("Contagem de arquivos por extensão", "contagem", c)

	regras := organizer.PegarRegras()
	movimentos, err := organizer.Classificar(downloads, regras, logger)
	if err != nil {
		logger.Error("Erro ao classificar arquivos", "erro", err)
		return fmt.Errorf("classificando arquivos: %w", err)
	}
	logger.Info("Arquivos classificados com sucesso", "quantidade", len(movimentos))

	err = organizer.MoverArquivos(movimentos, logger)
	if err != nil {
		logger.Error("Erro ao mover arquivos", "erro", err)
		return fmt.Errorf("movendo arquivos: %w", err)
	}
	//logger.Info("Arquivos movidos com sucesso")
	fmt.Println("Arquivos movidos com sucesso!")
	return nil
}
