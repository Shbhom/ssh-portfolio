package portfolio

type Overview struct {
	Intro   string   `yaml:"intro"`
	Bullets []string `yaml:"bullets"`
}

type ExperienceItem struct {
	Company string   `yaml:"company"`
	Role    string   `yaml:"role"`
	Period  string   `yaml:"period"`
	Bullets []string `yaml:"bullets"`
}

type Project struct {
	Name    string   `yaml:"name"`
	Tags    []string `yaml:"tags"`
	Summary string   `yaml:"summary"`
}

type Contact struct {
	Email    string `yaml:"email"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
}

type Portfolio struct {
	Name       string           `yaml:"name"`
	Tagline    string           `yaml:"tagline"`
	Overview   Overview         `yaml:"overview"`
	Experience []ExperienceItem `yaml:"experience"`
	Projects   []Project        `yaml:"projects"`
	Contact    Contact          `yaml:"contact"`
}
