package arke

import (
	"strings"

	. "gopkg.in/check.v1"
)

type ClassSuite struct{}

var _ = Suite(&ClassSuite{})

func (s *ClassSuite) TestNameMapping(c *C) {
	testdata := []struct {
		Name  string
		Class NodeClass
	}{
		{Name: "Zeus", Class: ZeusClass},
		{Name: "zeus", Class: ZeusClass},
		{Name: "zEUs", Class: ZeusClass},
		{Name: "Celaeno", Class: CelaenoClass},
		{Name: "celaeno", Class: CelaenoClass},
		{Name: "Helios", Class: HeliosClass},
		{Name: "helios", Class: HeliosClass},
		{Name: "Broadcast", Class: 0},
	}

	for _, d := range testdata {
		res, err := Class(d.Name)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(res, Equals, d.Class)

		c.Check(ClassName(res), Equals, strings.Title(strings.ToLower(d.Name)))
	}

	_, err := Class("hades")
	c.Check(err, ErrorMatches, `Unknown node class 'hades'`)

	c.Check(ClassName(NodeClass(1)), Equals, "<unknown>")
}
