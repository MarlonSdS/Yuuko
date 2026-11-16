package organizer

import (
	"fmt"
	"os"
	"path/filepath"
)

func MoverArquivos(movimentos []Movimento) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("obtendo diretório home: %w", err)
	}

	fmt.Println("Movendo arquivos...")
	for _, m := range movimentos {
		destinoDir := filepath.Join(homeDir, m.Destino)
		if err := os.MkdirAll(destinoDir, 0o755); err != nil {
			return fmt.Errorf("criando diretório %s: %w", destinoDir, err)
		}

		nomeArquivo, err := nomeLivre(destinoDir, m.Nome)
		if err != nil {
			return err
		}

		if err := os.Rename(m.Origem, filepath.Join(destinoDir, nomeArquivo)); err != nil {
			return fmt.Errorf("movendo arquivo %s para %s: %w", m.Origem, destinoDir, err)
		}
	}
	return nil
}
