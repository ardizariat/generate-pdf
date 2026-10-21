package cli

import (
	"fmt"
	"os"

	"charm.land/huh/v2/spinner"

	"github.com/ardizariat/generate-pdf/internal/config"
)

// progress shows that work is running. It uses an animated spinner only on an
// interactive terminal; otherwise it prints one plain line to stderr so stdout
// stays parseable. Verbose mode relies on zap logs instead.
func progress(cfg config.Config, label string, action func()) {
	if cfg.Verbose {
		action()
		return
	}
	if isCharDevice(os.Stdin) && isCharDevice(os.Stderr) {
		_ = spinner.New().Title(label).WithOutput(os.Stderr).Action(action).Run()
		return
	}
	fmt.Fprintln(os.Stderr, label)
	action()
}

func isCharDevice(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
