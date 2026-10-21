// Package config merges configuration sources (flags, env, file, defaults)
// into a single Config and validates it.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Supported document types.
const (
	TypeInvoice     = "invoice"
	TypeLetter      = "surat-resmi"
	TypeReport      = "laporan"
	TypeReceipt     = "kuitansi"
	TypeDelivery    = "surat-jalan"
	TypeCertificate = "sertifikat"
)

// typeList is the single source of the type list; also used by the validator tag.
const typeList = TypeInvoice + "," + TypeLetter + "," + TypeReport + "," + TypeReceipt + "," + TypeDelivery + "," + TypeCertificate

const (
	envPrefix  = "GENPDF"
	configName = "generate-pdf"
)

// configKeys is read-only: viper key = flag name.
var configKeys = [...]string{"type", "out", "rows", "seed", "title", "no-input", "force", "verify", "verbose", "config"}

type defaults struct {
	rows  int
	title string
}

// typeDefaults is read-only after initialization.
var typeDefaults = map[string]defaults{
	TypeInvoice:     {rows: 8, title: "Invoice"},
	TypeLetter:      {rows: 0, title: "Surat Resmi"},
	TypeReport:      {rows: 70, title: "Laporan"},
	TypeReceipt:     {rows: 0, title: "Kuitansi"},
	TypeDelivery:    {rows: 25, title: "Surat Jalan"},
	TypeCertificate: {rows: 0, title: "Sertifikat"},
}

// Types returns the supported document types.
func Types() []string {
	return strings.Split(typeList, ",")
}

// Config is the final configuration after all sources are merged.
type Config struct {
	Type       string `validate:"required,oneof=invoice surat-resmi laporan kuitansi surat-jalan sertifikat"`
	Out        string `validate:"required"`
	Title      string
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
		Out:        v.GetString("out"),
		Title:      v.GetString("title"),
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

	defaults := typeDefaults[loose.Type]
	if !loose.Provided("rows") {
		loose.Rows = defaults.rows
	}
	if !loose.Provided("title") {
		loose.Title = defaults.title
	}
	if !loose.Provided("out") {
		loose.Out = fmt.Sprintf("%s-%d.pdf", loose.Type, loose.Seed)
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
