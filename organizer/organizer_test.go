package organizer

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// escrever cria um arquivo (e as pastas necessárias) com o conteúdo dado.
func escrever(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatalf("criando pasta de %s: %v", caminho, err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatalf("escrevendo %s: %v", caminho, err)
	}
}

func TestClassificar(t *testing.T) {
	regras := map[string]Regra{
		".png": {Nome: "Imagens", Destino: "Imagens"},
		".nds": {Nome: "NDS", Destino: "Games/Roms/NDS"},
		".gz":  {Nome: "Comprimidos", Destino: "Comprimidos"},
	}

	tests := []struct {
		name     string
		arquivos []string    // arquivos criados na pasta de origem
		pastas   []string    // subpastas criadas na pasta de origem
		esperado []Movimento // Origem é preenchida pelo teste
	}{
		{
			name:     "extensão conhecida",
			arquivos: []string{"foto.png"},
			esperado: []Movimento{{Nome: "foto.png", Destino: "Imagens"}},
		},
		{
			name:     "extensão em maiúsculas",
			arquivos: []string{"FOTO.PNG"},
			esperado: []Movimento{{Nome: "FOTO.PNG", Destino: "Imagens"}},
		},
		{
			name:     "ROM vai para a subpasta do emulador",
			arquivos: []string{"jogo.nds"},
			esperado: []Movimento{{Nome: "jogo.nds", Destino: "Games/Roms/NDS"}},
		},
		{
			name:     "só a última extensão conta",
			arquivos: []string{"backup.tar.gz"},
			esperado: []Movimento{{Nome: "backup.tar.gz", Destino: "Comprimidos"}},
		},
		{
			name:     "sem extensão é ignorado",
			arquivos: []string{"LEIAME"},
			esperado: nil,
		},
		{
			name:     "extensão desconhecida é ignorada",
			arquivos: []string{"musica.xyz"},
			esperado: nil,
		},
		{
			name:     "pasta é ignorada mesmo com extensão conhecida",
			pastas:   []string{"album.png"},
			esperado: nil,
		},
		{
			name:     "vários arquivos, só os com regra entram",
			arquivos: []string{"a.png", "b.nds", "c.xyz"},
			esperado: []Movimento{
				{Nome: "a.png", Destino: "Imagens"},
				{Nome: "b.nds", Destino: "Games/Roms/NDS"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, nome := range tt.arquivos {
				escrever(t, filepath.Join(dir, nome), "x")
			}
			for _, nome := range tt.pastas {
				if err := os.Mkdir(filepath.Join(dir, nome), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			esperado := make([]Movimento, len(tt.esperado))
			for i, m := range tt.esperado {
				m.Origem = filepath.Join(dir, m.Nome)
				esperado[i] = m
			}

			got, err := Classificar(dir, regras)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !slices.Equal(got, esperado) {
				t.Errorf("obtido %+v, esperado %+v", got, esperado)
			}
		})
	}
}

func TestClassificarPastaInexistente(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nao-existe")

	_, err := Classificar(dir, PegarRegras())
	if err == nil {
		t.Fatal("esperava erro, obtive nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("esperava erro que envolva os.ErrNotExist, obtive: %v", err)
	}
}

// Guarda contra typos na tabela de regras: extensões sem ponto, com
// maiúsculas, ou repetidas em duas regras.
func TestRegrasValidas(t *testing.T) {
	vistas := make(map[string]string)
	for _, r := range regras {
		for _, ext := range r.Extensoes {
			if len(ext) < 2 || ext[0] != '.' {
				t.Errorf("regra %s: extensão %q deveria começar com ponto", r.Nome, ext)
			}
			if ext != strings.ToLower(ext) {
				t.Errorf("regra %s: extensão %q deveria estar em minúsculas", r.Nome, ext)
			}
			if outra, ok := vistas[ext]; ok {
				t.Errorf("extensão %q aparece em %s e em %s", ext, outra, r.Nome)
			}
			vistas[ext] = r.Nome
		}
	}
}

func TestMoverArquivos(t *testing.T) {
	tests := []struct {
		name       string
		existentes map[string]string // já na home: caminho relativo -> conteúdo
		origens    map[string]string // na pasta de origem: nome -> conteúdo
		movimentos []Movimento       // Origem é preenchida pelo teste
		esperado   map[string]string // na home após mover: caminho relativo -> conteúdo
		wantErr    bool
	}{
		{
			name:       "cria a pasta de destino e move",
			origens:    map[string]string{"foto.png": "novo"},
			movimentos: []Movimento{{Nome: "foto.png", Destino: "Imagens"}},
			esperado:   map[string]string{"Imagens/foto.png": "novo"},
		},
		{
			name:       "cria subpastas aninhadas",
			origens:    map[string]string{"jogo.nds": "novo"},
			movimentos: []Movimento{{Nome: "jogo.nds", Destino: "Games/Roms/NDS"}},
			esperado:   map[string]string{"Games/Roms/NDS/jogo.nds": "novo"},
		},
		{
			name:       "conflito: não sobrescreve o arquivo existente",
			existentes: map[string]string{"Imagens/foto.png": "antigo"},
			origens:    map[string]string{"foto.png": "novo"},
			movimentos: []Movimento{{Nome: "foto.png", Destino: "Imagens"}},
			esperado: map[string]string{
				"Imagens/foto.png":   "antigo",
				"Imagens/foto_2.png": "novo",
			},
		},
		{
			name: "conflito: usa o próximo número livre",
			existentes: map[string]string{
				"Imagens/foto.png":   "antigo1",
				"Imagens/foto_2.png": "antigo2",
			},
			origens:    map[string]string{"foto.png": "novo"},
			movimentos: []Movimento{{Nome: "foto.png", Destino: "Imagens"}},
			esperado: map[string]string{
				"Imagens/foto.png":   "antigo1",
				"Imagens/foto_2.png": "antigo2",
				"Imagens/foto_3.png": "novo",
			},
		},
		{
			name:    "vários arquivos em destinos diferentes",
			origens: map[string]string{"a.png": "A", "b.nds": "B"},
			movimentos: []Movimento{
				{Nome: "a.png", Destino: "Imagens"},
				{Nome: "b.nds", Destino: "Games/Roms/NDS"},
			},
			esperado: map[string]string{
				"Imagens/a.png":        "A",
				"Games/Roms/NDS/b.nds": "B",
			},
		},
		{
			name:       "lista vazia não faz nada",
			movimentos: nil,
			esperado:   nil,
		},
		{
			name:       "origem inexistente retorna erro",
			movimentos: []Movimento{{Nome: "fantasma.png", Destino: "Imagens"}},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origem := t.TempDir()
			home := t.TempDir()

			// Faz o os.UserHomeDir() apontar para a pasta temporária.
			t.Setenv("HOME", home)
			if got, err := os.UserHomeDir(); err != nil || got != home {
				t.Fatalf("HOME não foi redirecionado (got %q, err %v); abortando para não tocar na home real", got, err)
			}

			for nome, conteudo := range tt.origens {
				escrever(t, filepath.Join(origem, nome), conteudo)
			}
			for caminho, conteudo := range tt.existentes {
				escrever(t, filepath.Join(home, caminho), conteudo)
			}

			movimentos := make([]Movimento, len(tt.movimentos))
			for i, m := range tt.movimentos {
				m.Origem = filepath.Join(origem, m.Nome)
				movimentos[i] = m
			}

			err := MoverArquivos(movimentos)
			if (err != nil) != tt.wantErr {
				t.Fatalf("erro = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			for caminho, esperado := range tt.esperado {
				b, err := os.ReadFile(filepath.Join(home, caminho))
				if err != nil {
					t.Errorf("lendo %s: %v", caminho, err)
					continue
				}
				if string(b) != esperado {
					t.Errorf("%s: conteúdo %q, esperado %q", caminho, b, esperado)
				}
			}
			for _, m := range movimentos {
				if _, err := os.Stat(m.Origem); !errors.Is(err, os.ErrNotExist) {
					t.Errorf("a origem %s deveria ter sido movida", m.Origem)
				}
			}
		})
	}
}
