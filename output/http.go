package output

type Http struct {
	Url  string   `yaml:"url"`
	Method string `yaml:"method"`
	Headers []string `yaml:"headers"`
	Format  string   `yaml:"format"`
	Body string `yaml:"body"`
}