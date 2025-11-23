package portfolio

type Overview struct {
	Intro   string   `yaml:"intro"`
	Bullets []string `yaml:"bullets"`
}

type Experience struct {
	Company  string   `yaml:"company"`
	Role     string   `yaml:"role"`
	Period   string   `yaml:"period"`
	Location string   `yaml:"location"`
	Bullets  []string `yaml:"bullets"`
	Stack    string   `yaml:"stack"`
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
	Name        string       `yaml:"name"`
	Tagline     string       `yaml:"tagline"`
	Overview    Overview     `yaml:"overview"`
	Experiences []Experience `yaml:"experience"`
	Projects    []Project    `yaml:"projects"`
	Contact     Contact      `yaml:"contact"`
}
