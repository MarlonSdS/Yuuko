package organizer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func MoverArquivos(movimentos []Movimento, logger *slog.Logger) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("obtendo diretório home: %w", err)
	}

	logger.Info("Movendo arquivos...")

	movidos, renomeados := 0, 0
	for _, m := range movimentos {
		destinoDir := filepath.Join(homeDir, m.Destino)
		if err := os.MkdirAll(destinoDir, 0o755); err != nil {
			logger.Error("Erro ao criar diretório de destino", "diretorio", destinoDir, "erro", err)
			return fmt.Errorf("criando diretório %s: %w", destinoDir, err)
		}

		nomeArquivo, err := nomeLivre(destinoDir, m.Nome)
		if err != nil {
			logger.Error("Erro ao gerar nome livre para arquivo", "nome_original", m.Nome,
				"diretorio", destinoDir, "erro", err)
			return err
		}
		renomeado := nomeArquivo != m.Nome

		if err := os.Rename(m.Origem, filepath.Join(destinoDir, nomeArquivo)); err != nil {
			logger.Error(
				"Erro ao mover arquivo", "origem", m.Origem, "destino", filepath.Join(destinoDir, nomeArquivo), "erro", err)
			return fmt.Errorf("movendo arquivo %s para %s: %w", m.Origem, destinoDir, err)
		} else {

			if renomeado {
				renomeados++
			}
		}
		logger.Info(
			"Arquivo movido com sucesso", "origem", m.Origem, "destino", filepath.Join(destinoDir, nomeArquivo),
			"renomeado", renomeado, "nome_original", m.Nome, "novo_nome", nomeArquivo)
		movidos++
	}
	logger.Info("Resumo da operação", "movidos", movidos, "renomeados", renomeados)
	return nil
}
