package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

// Job is one batch entry after inheriting the global configuration.
type Job struct {
	Type  string
	Pages int
	Seed  uint64
	Out   string
	Title string
}

// manifest is the YAML shape of a batch file; unknown fields are ignored.
type manifest struct {
	Jobs []jobEntry `yaml:"jobs"`
}

// jobEntry keeps pointers for pages and seed so an omitted value inherits the
// global configuration while an explicit zero is still validated.
type jobEntry struct {
	Type  string  `yaml:"type"`
	Pages *int    `yaml:"pages"`
	Seed  *uint64 `yaml:"seed"`
	Out   string  `yaml:"out"`
	Title string  `yaml:"title"`
}

// LoadBatch reads a YAML manifest and resolves every job against base.
func LoadBatch(path string, base Config) ([]Job, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baca manifest %s: %w", path, err)
	}

	var m manifest
	if err := yaml.Unmarshal(content, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if len(m.Jobs) == 0 {
		return nil, fmt.Errorf("manifest %s: jobs kosong", path)
	}

	jobs := make([]Job, 0, len(m.Jobs))
	for _, entry := range m.Jobs {
		job := Job{
			Type:  entry.Type,
			Pages: 1,
			Seed:  base.Seed,
			Out:   entry.Out,
			Title: entry.Title,
		}
		if job.Type == "" {
			job.Type = base.Type
		}
		if entry.Pages != nil {
			job.Pages = *entry.Pages
		}
		if entry.Seed != nil {
			job.Seed = *entry.Seed
		}
		if job.Title == "" {
			if base.Provided("title") {
				job.Title = base.Title
			} else {
				job.Title = typeDefaults[job.Type].title
			}
		}
		job.Out = OutputPath(base.OutDir, job.Out)
		if job.Out == "" {
			name, err := DefaultOutputName(job.Type, time.Now())
			if err != nil {
				return nil, fmt.Errorf("manifest %s: %w", path, err)
			}
			job.Out = OutputPath(base.OutDir, name)
		}
		jobs = append(jobs, job)
	}

	if err := ValidateJobs(jobs); err != nil {
		return nil, fmt.Errorf("manifest %s: %w", path, err)
	}
	return jobs, nil
}

// ValidateJobs rejects an empty batch and any job that cannot be generated.
func ValidateJobs(jobs []Job) error {
	if len(jobs) == 0 {
		return fmt.Errorf("batch: minimal satu job")
	}

	messages := make([]string, 0, len(jobs))
	for i, job := range jobs {
		if err := Validate(Config{Type: job.Type, Out: job.Out}); err != nil {
			messages = append(messages, fmt.Sprintf("job %d: %v", i+1, err))
		}
		if job.Pages < 1 {
			messages = append(messages, fmt.Sprintf("job %d: Pages=%d (gte 1)", i+1, job.Pages))
		}
		if job.Pages > 1 && !SupportsMultiPage(job.Type) {
			messages = append(messages, fmt.Sprintf("job %d: Pages=%d tidak didukung untuk jenis %s (maksimal 1)", i+1, job.Pages, job.Type))
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("batch tidak valid: %s", strings.Join(messages, "; "))
}

// SupportsMultiPage reports whether a document type may exceed one page.
func SupportsMultiPage(docType string) bool {
	switch docType {
	case TypeInvoice, TypeReport, TypeDelivery:
		return true
	default:
		return false
	}
}
