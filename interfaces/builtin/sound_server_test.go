// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2026 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin_test

import (
	"gopkg.in/check.v1"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/builtin"
)

type SoundServerSuite struct {
	iface interfaces.Interface
}

var _ = check.Suite(&SoundServerSuite{})

func (s *SoundServerSuite) SetUpSuite(c *check.C) {
	s.iface = builtin.MustInterface("sound-server")
}

func (s *SoundServerSuite) TestGlobalSocketActivationPaths(c *check.C) {
	provider, ok := s.iface.(interfaces.GlobalSocketActivationPathsProvider)
	c.Assert(ok, check.Equals, true, check.Commentf("sound-server interface should implement GlobalSocketActivationPathsProvider"))

	paths := provider.GlobalSocketActivationPaths()
	c.Assert(len(paths), check.Not(check.Equals), 0, check.Commentf("should provide at least one global path"))

	expectedPaths := map[string]bool{
		"/run/user/[0-9]*/pulse/native":   false,
		"/run/user/[0-9]*/pipewire-[0-9]": false,
	}

	for _, path := range paths {
		if _, exists := expectedPaths[path]; exists {
			expectedPaths[path] = true
		}
	}

	for path, found := range expectedPaths {
		c.Assert(found, check.Equals, true, check.Commentf("expected path %s not found", path))
	}
}
