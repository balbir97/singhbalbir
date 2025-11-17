package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ----------------------------
// YAML Structs
// ----------------------------

type Resume struct {
	Name        string            `yaml:"name"`
	Titles      string            `yaml:"titles"`
	Contact     Contact           `yaml:"contact"`
	Punchline   string            `yaml:"punchline"`
	Education   []EducationEntry  `yaml:"education"`
	Skills      Skills            `yaml:"skills"`
	Experience  []ExperienceEntry `yaml:"experience"`
	Interests   []InterestEntry   `yaml:"interests"`
	References  string            `yaml:"references"`
	LastUpdated string            `yaml:"last_updated"`
}

type Contact struct {
	Location string `yaml:"location"`
	Phone    string `yaml:"phone"`
	Email    string `yaml:"email"`
	Github   string `yaml:"github"`
	Twitter  string `yaml:"twitter"`
	LinkedIn string `yaml:"linkedin"`
}

type EducationEntry struct {
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

type ExperienceEntry struct {
	Company      string   `yaml:"company_short"`
	Url          string   `yaml:"url"`
	Location     string   `yaml:"location"`
	Title        string   `yaml:"title_short"`
	Period       Period   `yaml:"period"`
	Achievements []string `yaml:"achievements"`
}

type Period struct {
	From MonthYear `yaml:"from"`
	To   MonthYear `yaml:"to"`
}

type MonthYear struct {
	Month string `yaml:"month"`
	Year  string `yaml:"year"`
}

type InterestEntry struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// ----------------------------
// Template Struct (expected by HTML)
// ----------------------------

type TemplateData struct {
	Name       string
	Titles     string
	Punchline  string
	Email      string
	Experience []TemplateExperience
	Social     []TemplateSocial
	Year       int
}

type TemplateExperience struct {
	Role     string
	Company  string
	Location string
	Period   string
	Bullets  []string
	Link     string
}

type TemplateSocial struct {
	Name string
	URL  string
}

// ----------------------------
// Load YAML
// ----------------------------

func loadYAML(path string) (Resume, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Resume{}, err
	}

	var r Resume
	err = yaml.Unmarshal(data, &r)
	return r, err
}

// ----------------------------
// Main
// ----------------------------

func main() {
	// Load resume.yaml
	resume, err := loadYAML("resume.yaml")
	if err != nil {
		fmt.Println("Error reading YAML:", err)
		return
	}

	// Map YAML → TemplateData
	td := TemplateData{
		Name:      resume.Name,
		Titles:    resume.Titles,
		Punchline: resume.Punchline,
		Email:     resume.Contact.Email,
		Year:      time.Now().Year(),
		Social: []TemplateSocial{
			{Name: "Github", URL: resume.Contact.Github},
			{Name: "Twitter", URL: resume.Contact.Twitter},
			{Name: "LinkedIn", URL: resume.Contact.LinkedIn},
		},
	}

	// Convert experience
	for _, exp := range resume.Experience {
		// Build the period string with Present handling
		var period string
		if exp.Period.To.Month == "-" || exp.Period.To.Year == "-" {
			period = fmt.Sprintf("%s %s to Present",
				exp.Period.From.Month,
				exp.Period.From.Year,
			)
		} else {
			period = fmt.Sprintf("%s %s to %s %s",
				exp.Period.From.Month,
				exp.Period.From.Year,
				exp.Period.To.Month,
				exp.Period.To.Year,
			)
		}

		td.Experience = append(td.Experience, TemplateExperience{
			Role:     exp.Title,
			Company:  exp.Company,
			Location: exp.Location,
			Period:   period,
			Bullets:  exp.Achievements,
			Link:     exp.Url,
		})
	}

	// Load template file
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return
	}

	// Output file
	out, err := os.Create("index.html")
	if err != nil {
		fmt.Println("Error creating output:", err)
		return
	}
	defer out.Close()

	// Execute template
	err = tmpl.Execute(out, td)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		return
	}

	fmt.Println("Generated index.html successfully.")
}
