package command

import (
	"os"
	"testing"

	"github.com/hashicorp/hcptf-cli/internal/config"
)

func TestMain(m *testing.M) {
	_ = os.Setenv(config.DisableEnvFileVariable, "1")
	os.Exit(m.Run())
}
