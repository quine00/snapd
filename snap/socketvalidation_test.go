// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2025 Canonical Ltd
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

package snap_test

import (
	"gopkg.in/check.v1"

	"github.com/snapcore/snapd/snap"
)

type SocketValidationSuite struct{}

var _ = check.Suite(&SocketValidationSuite{})

// Mock interface that provides global socket activation paths
type mockGlobalSocketInterface struct {
	paths []string
}

func (m *mockGlobalSocketInterface) GlobalSocketActivationPaths() []string {
	return m.paths
}

func (s *SocketValidationSuite) TestValidateSocketActivationWithInterfaces(c *check.C) {
	snapInfo := &snap.Info{
		SuggestedName: "test-snap",
		Apps: map[string]*snap.AppInfo{
			"service": {
				Name:        "service",
				DaemonScope: snap.SystemDaemon,
				Sockets: map[string]*snap.SocketInfo{
					"sock": {
						Name:         "sock",
						ListenStream: "$SNAP_REAL_XDG_RUNTIME_DIR/pulse/native",
					},
				},
			},
		},
		Slots: map[string]*snap.SlotInfo{
			"audio-playback": {
				Name:      "audio-playback",
				Interface: "audio-playback",
			},
		},
	}

	// Set up backlinks
	snapInfo.Apps["service"].Snap = snapInfo
	snapInfo.Apps["service"].Sockets["sock"].App = snapInfo.Apps["service"]
	snapInfo.Slots["audio-playback"].Snap = snapInfo

	// Create mock interface providers
	interfaceProviders := map[string]any{
		"audio-playback": &mockGlobalSocketInterface{
			paths: []string{"/run/user/[0-9]*/pulse/native"},
		},
	}

	// Test validation with allowed global path
	err := snap.ValidateSocketActivationWithInterfaces(snapInfo, interfaceProviders)
	c.Assert(err, check.IsNil)

	// Test validation with disallowed path
	snapInfo.Apps["service"].Sockets["sock"].ListenStream = "/run/forbidden/path"
	err = snap.ValidateSocketActivationWithInterfaces(snapInfo, interfaceProviders)
	c.Assert(err, check.ErrorMatches, `invalid definition of socket "sock": .*`)
}

func (s *SocketValidationSuite) TestValidateSocketWithGlobalPaths(c *check.C) {
	// This test would need to be made by making validateSocketWithGlobalPaths exported
	// or testing through the public ValidateSocketActivationWithInterfaces function
}

func (s *SocketValidationSuite) TestIsAllowedGlobalPath(c *check.C) {
	// This test would need to be made by making isAllowedGlobalPath exported
	// or testing through the public ValidateSocketActivationWithInterfaces function
}

func (s *SocketValidationSuite) TestExpandSocketPath(c *check.C) {
	// This test would need to be made by making expandSocketPath exported
	// or testing through the public ValidateSocketActivationWithInterfaces function
}
