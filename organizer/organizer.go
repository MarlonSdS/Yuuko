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

type Regra struct {
	Nome      string   `json:"nome"`
	Extensoes []string `json:"extensoes"`
	Destino   string   `json:"destino"`
}

type Movimento struct {
	Nome    string // Nome do arquivo
	Origem  string // Caminho completo do arquivo de origem
	Destino string // Diretório de destino relativo à home
}

var regras = []Regra{
	{Nome: "Imagens", Extensoes: []string{".png", ".jpg", ".jpeg",
		".gif", ".webp"}, Destino: "Imagens"},
	{Nome: "Videos", Extensoes: []string{".mp4", ".avi", ".mkv"},
		Destino: "Vídeos"},
	{Nome: "Documentos", Extensoes: []string{".pdf", ".docx", ".txt",
		".xlsx", ".csv"}, Destino: "Documentos"},
	{Nome: "Musicas", Extensoes: []string{".mp3", ".wav", ".flac"},
		Destino: "Músicas"},
	{Nome: "Comprimidos", Extensoes: []string{".zip", ".rar", ".7z",
		".tar", ".gz", ".xz"}, Destino: "Comprimidos"},
	{Nome: "Programas", Extensoes: []string{".appimage", ".deb", ".rpm",
		".flatpak", ".exe"}, Destino: "Programas"},
	{Nome: "NDS", Extensoes: []string{".nds"}, Destino: "Games/Roms/NDS"},
	{Nome: "GBA", Extensoes: []string{".gba"}, Destino: "Games/Roms/GBA"},
	{Nome: "SNES", Extensoes: []string{".smc"}, Destino: "Games/Roms/SNES"},
	{Nome: "NES", Extensoes: []string{".nes"}, Destino: "Games/Roms/NES"},
	{Nome: "PSX", Extensoes: []string{".iso", ".bin"}, Destino: "Games/Roms/PSX"},
	{Nome: "3DS", Extensoes: []string{".3ds"}, Destino: "Games/Roms/3DS"},
	{Nome: "Switch", Extensoes: []string{".nsp", ".xci"}, Destino: "Games/Roms/Switch"},
}

func PegarRegras() map[string]Regra {
	regrasMap := make(map[string]Regra)
	for _, regra := range regras {
		for _, extensao := range regra.Extensoes {
			regrasMap[extensao] = regra
		}
	}

	return regrasMap
}

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
