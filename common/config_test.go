package common

import (
	"fmt"
	"testing"
)

func TestGenServerConfig(t *testing.T) {
	SetServerConfigPath("./settings.json")
	cfg := GenServerConfig()
	fmt.Println(cfg)
}
