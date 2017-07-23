package conf_test

import (
	"testing"

	"github.com/juusechec/goresource-proxy/conf"
)

func TestConf(t *testing.T) {
	p := conf.Parameters
	if p.HOSTNAME != "localhost" {
		t.Error("unspectec HOSTNAME")
	}
	if p.CONTEXT != "/proxy" {
		t.Error("unspectec CONTEXT")
	}
	if p.PORT != "12345" {
		t.Error("unspectec PORT")
	}
}
