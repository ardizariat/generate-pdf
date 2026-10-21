// Package config merges configuration sources (flags, env, file, defaults)
// into a single Config and validates it.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ardizariat/generate-pdf/internal/locale"
)

// Supported document types. Slugs stay English in every language.
const (
	TypeInvoice     = "invoice"
	TypeLetter      = "letter"
	TypeReport      = "report"
	TypeReceipt     = "receipt"
	TypeDelivery    = "delivery-note"
	TypeCertificate = "certificate"
)

// typeList is the single source of the type list; also used by the validator tag.
const typeList = TypeInvoice + "," + TypeLetter + "," + TypeReport + "," + TypeReceipt + "," + TypeDelivery + "," + TypeCertificate

const (
	envPrefix  = "GENPDF"
	configName = "generate-pdf"
)

// DefaultOutDir is where generated PDFs go when no path is given.
const DefaultOutDir = "output"

// configKeys is read-only: viper key = flag name.
var configKeys = [...]string{"type", "lang", "out", "out-dir", "rows", "seed", "title", "batch", "no-input", "force", "verify", "verbose", "config"}

type defaults struct {
	rows int
}

// typeDefaults is read-only after initialization.
var typeDefaults = map[string]defaults{
	TypeInvoice:     {rows: 8},
	TypeLetter:      {rows: 0},
	TypeReport:      {rows: 70},
	TypeReceipt:     {rows: 0},
	TypeDelivery:    {rows: 25},
	TypeCertificate: {rows: 0},
}

// Types returns the supported document types.
func Types() []string {
	return strings.Split(typeList, ",")
}

// DefaultRows returns the default table row count for docType.
func DefaultRows(docType string) int {
	return typeDefaults[docType].rows
}

// OutputPath places a relative path under baseDir. Absolute paths and an empty
// baseDir are returned unchanged.
func OutputPath(baseDir, path string) string {
	if path == "" || baseDir == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

// DefaultOutputName builds a unique default name:
// <type>_<UTC timestamp>_<uuid v7>.pdf.
func DefaultOutputName(docType string, at time.Time) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("create output uuid: %w", err)
	}
	return fmt.Sprintf("%s_%s_%s.pdf", docType, at.UTC().Format("20060102T150405Z"), id.String()), nil
}

// Config is the final configuration after all sources are merged.
type Config struct {
	Type       string `validate:"required,oneof=invoice letter report receipt delivery-note certificate"`
	Language   string `validate:"omitempty,oneof=en id"`
	Out        string `validate:"required"`
	OutDir     string
	Title      string
	Batch      string
	Rows       int `validate:"gte=0"`
	Seed       uint64
	NoInput    bool
	Verify     bool
	Verbose    bool
	Force      bool
	ConfigPath string

	provided map[string]bool
}

// Provided reports whether key came from a changed flag, a GENPDF_* env var, or a config file.
func (c Config) Provided(key string) bool {
	return c.provided[key]
}

// Load merges flags, GENPDF_* env vars, config file, and per-type defaults.
func Load(cmd *cobra.Command) (Config, error) {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	if cmd != nil {
		for _, key := range configKeys {
			flag := cmd.Flags().Lookup(key)
			if flag == nil {
				continue
			}
			if err := v.BindPFlag(key, flag); err != nil {
				return Config{}, fmt.Errorf("bind flag %s: %w", key, err)
			}
		}
	}

	if err := readConfigFile(v, cmd); err != nil {
		return Config{}, err
	}

	loose := Config{
		Type:       v.GetString("type"),
		Language:   v.GetString("lang"),
		Out:        v.GetString("out"),
		OutDir:     v.GetString("out-dir"),
		Title:      v.GetString("title"),
		Batch:      v.GetString("batch"),
		Rows:       v.GetInt("rows"),
		Seed:       v.GetUint64("seed"),
		NoInput:    v.GetBool("no-input"),
		Verify:     v.GetBool("verify"),
		Verbose:    v.GetBool("verbose"),
		Force:      v.GetBool("force"),
		ConfigPath: v.ConfigFileUsed(),
		provided:   providedKeys(v, cmd),
	}

	if loose.Type == "" {
		loose.Type = TypeInvoice
	}
	if loose.Language == "" {
		loose.Language = string(locale.English)
	}
	if loose.OutDir == "" {
		loose.OutDir = DefaultOutDir
	}

	if !loose.Provided("rows") {
		loose.Rows = DefaultRows(loose.Type)
	}
	if !loose.Provided("title") {
		loose.Title = locale.Get(locale.OrDefault(loose.Language)).TypeTitles[loose.Type]
	}
	if !loose.Provided("out") {
		name, err := DefaultOutputName(loose.Type, time.Now())
		if err != nil {
			return Config{}, err
		}
		loose.Out = OutputPath(loose.OutDir, name)
	}

	return loose, nil
}

// Validate checks Config; it is called after Load and again after prompting.
func Validate(c Config) error {
	err := validator.New().Struct(c)
	if err == nil {
		return nil
	}

	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		return fmt.Errorf("validate config: %w", err)
	}

	messages := make([]string, 0, len(invalid))
	for _, fieldErr := range invalid {
		messages = append(messages, validationMessage(fieldErr))
	}
	return fmt.Errorf("invalid config: %s", strings.Join(messages, "; "))
}

func validationMessage(fieldErr validator.FieldError) string {
	if param := fieldErr.Param(); param != "" {
		return fmt.Sprintf("%s=%v (%s %s)", fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), param)
	}
	return fmt.Sprintf("%s=%v (%s)", fieldErr.Field(), fieldErr.Value(), fieldErr.Tag())
}

func readConfigFile(v *viper.Viper, cmd *cobra.Command) error {
	if flagChanged(cmd, "config") {
		v.SetConfigFile(cmd.Flags().Lookup("config").Value.String())
	} else if path := os.Getenv(envPrefix + "_CONFIG"); path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName(configName)
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	err := v.ReadInConfig()
	var notFound viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func providedKeys(v *viper.Viper, cmd *cobra.Command) map[string]bool {
	provided := make(map[string]bool, len(configKeys))
	for _, key := range configKeys {
		provided[key] = flagChanged(cmd, key) || envSet(key) || v.InConfig(key)
	}
	return provided
}

func flagChanged(cmd *cobra.Command, key string) bool {
	if cmd == nil {
		return false
	}
	flag := cmd.Flags().Lookup(key)
	return flag != nil && flag.Changed
}

func envSet(key string) bool {
	name := envPrefix + "_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	_, ok := os.LookupEnv(name)
	return ok
}
