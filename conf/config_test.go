package conf_test

import (
	"testing"

	"github.com/acd19ml/TalentRank/conf"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromToml(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../etc/demo.toml")
	if should.NoError(err) {
		should.Equal("TalentRank", conf.C().App.Name)
	}
}
