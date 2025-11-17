package main

import (
	"html/template"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Contact struct {
	Location string `yaml:"location"`
	Phone    string `yaml:"phone"`
	Email    string `yaml:"email"`
	Website  string `yaml:"website"`
	LinkedIn string `yaml:"linkedin"`
}

type Education struct {
	Institution string   `yaml:"institution"`
	Location    string   `yaml:"location"`
	Degree      string   `yaml:"degree"`
	Year        int      `yaml:"year"`
	Highlights  []string `yaml:"highlights"`
}

type Skills struct {
	Platform      []string `yaml:"platform"`
	Deployment    []string `yaml:"deployment"`
	Programming   []string `yaml:"programming"`
	Observability []string `yaml:"observability"`
	CICD          []string `yaml:"cicd"`
	Services      []string `yaml:"services"`
}

type Period struct {
	From struct {
		Month string `yaml:"month"`
		Year  string `yaml:"year"`
	} `yaml:"from"`

	To struct {
		Month string `yaml:"month"`
		Year  string `yaml:"year"`
	} `yaml:"to"`
}

type Experience struct {
	Company      string   `yaml:"company"`
	URL          string   `yaml:"url"`
	Location     string   `yaml:"location"`
	Title        string   `yaml:"title"`
	Period       Period   `yaml:"period"`
	Achievements []string `yaml:"achievements"`
}

type Interest struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Resume struct {
	Name        string       `yaml:"name"`
	Contact     Contact      `yaml:"contact"`
	Profile     string       `yaml:"profile"`
	Education   []Education  `yaml:"education"`
	Skills      Skills       `yaml:"skills"`
	Experience  []Experience `yaml:"experience"`
	Interests   []Interest   `yaml:"interests"`
	References  string       `yaml:"references"`
	LastUpdated string       `yaml:"last_updated"`

	// auto populated
	Year int
}

func main() {
	// Read YAML
	yamlBytes, err := os.ReadFile("resume.yaml")
	if err != nil {
		log.Fatalf("cannot read resume.yaml: %v", err)
	}

	var r Resume
	if err := yaml.Unmarshal(yamlBytes, &r); err != nil {
		log.Fatalf("cannot unmarshal YAML: %v", err)
	}

	// Inject current year for footer
	r.Year = time.Now().Year()

	// Load template
	tpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalf("template parse error: %v", err)
	}

	// Output file
	out, err := os.Create("index.html")
	if err != nil {
		log.Fatalf("cannot create index.html: %v", err)
	}
	defer out.Close()

	// Execute template
	if err := tpl.Execute(out, r); err != nil {
		log.Fatalf("template execution failed: %v", err)
	}

	log.Println("index.html generated successfully")
}
