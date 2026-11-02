package organizer

import (
	"fmt"
	"os"
	"path/filepath"
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

func MoverArquivos(movimentos []Movimento) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("obtendo diretório home: %w", err)
	}

	fmt.Printf("Movendo arquivos... \n")
	for _, m := range movimentos {
		destinoDir := filepath.Join(homeDir, m.Destino)
		if err := os.MkdirAll(destinoDir, os.ModePerm); err != nil {
			return fmt.Errorf("criando diretório %s: %w", destinoDir, err)
		}

		nomeArquivo := m.Nome
		// Verifica se o arquivo de destino já existe
		if _, err := os.Lstat(filepath.Join(destinoDir, nomeArquivo)); err == nil {
			ext := filepath.Ext(m.Nome)
			nome := strings.TrimSuffix(m.Nome, ext)
			for i := 2; ; i++ {
				nomeArquivo = fmt.Sprintf("%s_%d%s", nome, i, ext)
				_, err := os.Lstat(filepath.Join(destinoDir, nomeArquivo))
				if os.IsNotExist(err) {
					break
				}
				if err != nil {
					return fmt.Errorf("verificando arquivo de destino %s: %w", filepath.Join(destinoDir, nomeArquivo), err)
				}
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("verificando arquivo de destino %s: %w", filepath.Join(destinoDir, nomeArquivo), err)
		}

		if err := os.Rename(m.Origem, filepath.Join(destinoDir, nomeArquivo)); err != nil {
			return fmt.Errorf("movendo arquivo %s para %s: %w", m.Origem, destinoDir, err)
		}
	}
	return nil
}
