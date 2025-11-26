// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func TestSortDeps(t *testing.T) {
	suite.Run(t, new(sortDepsSuite))
}

type sortDepsSuite struct {
	suite.Suite
	A, B, C, D, E *Dep[int]
}

func (s *sortDepsSuite) SetupTest() {
	f := func(key string) *Dep[int] {
		ret := NewDep(func(_ *gin.Context) int { return 0 })
		ret.setKey(key)
		return ret
	}
	s.A = f("A")
	s.B = f("B")
	s.C = f("C")
	s.D = f("D")
	s.E = f("E")
}

func (s *sortDepsSuite) dump(arr []anotherDep) string {
	ret := make([]string, len(arr))
	for idx, d := range arr {
		ret[idx] = d.getKey()
	}
	return strings.Join(ret, ", ")
}

func (s *sortDepsSuite) Test_Simple() {
	deps := []anotherDep{
		s.A, s.B, s.C, s.D, s.E,
	}
	expect := []anotherDep{
		s.E, s.D, s.C, s.B, s.A,
	}

	actual, err := sortDeps(deps)
	s.Require().NoError(err)
	s.ElementsMatch(expect, actual, "expect: %s\nactual: %s", s.dump(expect), s.dump(actual))
}

func (s *sortDepsSuite) Test_Nil() {
	actual, err := sortDeps(nil)
	s.Require().NoError(err)
	s.Empty(actual)
}

func (s *sortDepsSuite) Test_Empty() {
	actual, err := sortDeps([]anotherDep{})
	s.Require().NoError(err)
	s.Empty(actual)
}

func (s *sortDepsSuite) Test_Duplicated() {
	DependsOn(s.B, s.A)
	DependsOn(s.C, s.B)
	DependsOn(s.D, s.C)
	DependsOn(s.E, s.D)
	deps := []anotherDep{s.A, s.B, s.C, s.D, s.E}
	expect := []anotherDep{
		s.E, s.D, s.C, s.B, s.A,
	}

	actual, err := sortDeps(deps)
	s.Require().NoError(err)
	s.Equal(expect, actual, "expected: %s\nactual: %s", s.dump(expect), s.dump(actual))
}

func (s *sortDepsSuite) Test_Complex() {
	// D depends B depends A
	DependsOn(s.B, s.D)
	DependsOn(s.A, s.B)
	// C depends B depends A
	DependsOn(s.B, s.C)
	// E depends B
	DependsOn(s.B, s.E)

	deps := []anotherDep{s.D, s.B, s.A, s.C, s.E}
	expect := []string{
		"A", "B", "C", "D", "E",
	}
	actual, err := sortDeps(deps)
	s.Require().NoError(err)
	s.Require().Len(actual, len(expect))

	s.Require().ElementsMatch(deps, actual)

	idx := func(k string) int {
		for idx, v := range actual {
			if v.getKey() == k {
				return idx
			}
		}
		s.FailNow("", "%s does not exist in result %s", k, s.dump(actual))
		return -1
	}
	s.Less(idx("A"), idx("B"))
	s.Less(idx("B"), idx("C"))
	s.Less(idx("B"), idx("D"))
	s.Less(idx("B"), idx("E"))
}
