package resiver

type ResiverConfig struct {
	Name     string   `yaml:"name"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Position Position `yaml:"position"`
	Type     string   `yaml:"type"`
}

type Position struct {
	Lat float64 `yaml:"lat"`
	Lon float64 `yaml:"lon"`
}
