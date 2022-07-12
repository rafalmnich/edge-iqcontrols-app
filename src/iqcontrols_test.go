package main_test

import (
	"testing"

	"github.com/futurehomeno/cliffhanger/test/suite"

	"github.com/rafalmnich/edge-iqcontrols-app/test"
)

func TestConfiguration(t *testing.T) {
	t.SkipNow()
	s := &suite.Suite{
		Cases: []*suite.Case{
			{
				Name:     "Configuration",
				Setup:    test.ServiceSetup("not_configured"),
				TearDown: []suite.Callback{test.TearDown("not_configured")},
				Nodes: []*suite.Node{
					{
						Name:    "Configure log level",
						Command: suite.StringMessage("pt:j1/mt:cmd/rt:app/rn:iqcontrols/ad:1", "cmd.log.set_level", "iqcontrols", "warn"),
						Expectations: []*suite.Expectation{
							suite.ExpectString("pt:j1/mt:evt/rt:app/rn:iqcontrols/ad:1", "evt.log.level_report", "iqcontrols", "warning"),
						},
					},
				},
			},
		},
	}

	s.Run(t)
}
