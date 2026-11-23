package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func definirPasta() (string, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" || !filepath.IsAbs(stateHome) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("localizando o diretório do usuário: %w", err)
		}
		stateHome = filepath.Join(home, ".local/state/yuuko/logs")
	} else {
		stateHome = filepath.Join(stateHome, "yuuko/logs")
	}
	return stateHome, nil
}

func montarNomeArquivoLog(t time.Time) string {
	dataFormatada := t.Format("2006-01-02_15-04-05")
	nomeArquivo := "yuuko-" + dataFormatada + ".log"
	return nomeArquivo
}

func CriarArquivoLog() (*os.File, error) {
	pasta, err := definirPasta()
	if err != nil {
		return nil, fmt.Errorf("definindo a pasta de logs: %w", err)
	}

	nomeArquivo := montarNomeArquivoLog(time.Now())

	err = os.MkdirAll(pasta, 0o755)
	if err != nil {
		return nil, fmt.Errorf("criando o diretório de logs %s: %w", pasta, err)
	}

	caminhoCompleto := filepath.Join(pasta, nomeArquivo)
	arquivo, err := os.OpenFile(caminhoCompleto, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("abrindo o arquivo de log %s: %w", caminhoCompleto, err)
	}
	return arquivo, nil
}

func NovoLogger(w io.Writer) *slog.Logger {
	opcoes := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(w, opcoes)
	return slog.New(handler)
}
