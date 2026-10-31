package organizer

type Regra struct {
	Nome      string   `json:"nome"`
	Extensoes []string `json:"extensoes"`
	Destino   string   `json:"destino"`
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
