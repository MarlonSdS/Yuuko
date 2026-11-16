package organizer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Reconhece o contador que o próprio Yuuko gera: "foto (2)".
var contador = regexp.MustCompile(`^(.*) \((\d+)\)$`)

// existe informa se há algo (arquivo, pasta ou link) em caminho.
func existe(caminho string) (bool, error) {
	_, err := os.Lstat(caminho)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// separarExtensao divide "foto.tar.gz" em ("foto", ".tar.gz"),
// tratando .tar.* como extensão composta.
func separarExtensao(nome string) (stem, ext string) {
	ext = filepath.Ext(nome)
	stem = strings.TrimSuffix(nome, ext)
	if len(stem) > 4 && strings.EqualFold(stem[len(stem)-4:], ".tar") {
		ext = stem[len(stem)-4:] + ext
		stem = stem[:len(stem)-4]
	}
	return stem, ext
}

// nomeLivre devolve um nome que ainda não existe em dir. Se nome já está
// livre, devolve o próprio nome; senão acrescenta (2), (3)... Se o nome já
// termina em um contador "(N)", continua a partir dele em vez de empilhar.
func nomeLivre(dir, nome string) (string, error) {
	ocupado, err := existe(filepath.Join(dir, nome))
	if err != nil {
		return "", fmt.Errorf("verificando %s: %w", filepath.Join(dir, nome), err)
	}
	if !ocupado {
		return nome, nil
	}

	stem, ext := separarExtensao(nome)
	base, n := stem, 1
	if m := contador.FindStringSubmatch(stem); m != nil {
		if num, err := strconv.Atoi(m[2]); err == nil {
			base, n = m[1], num
		}
	}

	for i := n + 1; ; i++ {
		candidato := fmt.Sprintf("%s (%d)%s", base, i, ext)
		ocupado, err := existe(filepath.Join(dir, candidato))
		if err != nil {
			return "", fmt.Errorf("verificando %s: %w", filepath.Join(dir, candidato), err)
		}
		if !ocupado {
			return candidato, nil
		}
	}
}
