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

package snap

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GlobalSocketActivationPathsProvider defines an interface that provides
// allowed paths for socket activation outside the standard snap directories.
// This interface is expected to be implemented by interfaces that need
// privileged socket access.
type GlobalSocketActivationPathsProvider interface {
	GlobalSocketActivationPaths() []string
}

// ValidateSocketActivationWithInterfaces validates socket activation paths
// against both standard snap directories and interface-specific global paths.
// This function should be called when interface information is available,
// typically during snap installation.
func ValidateSocketActivationWithInterfaces(snapInfo *Info, interfaceProviders map[string]any) error {
	for _, app := range snapInfo.Apps {
		if len(app.Sockets) == 0 {
			continue
		}

		// Get all allowed global paths from connected slots
		var globalPaths []string
		for _, slot := range app.Slots {
			if provider, ok := interfaceProviders[slot.Interface]; ok {
				if socketProvider, ok := provider.(GlobalSocketActivationPathsProvider); ok {
					globalPaths = append(globalPaths, socketProvider.GlobalSocketActivationPaths()...)
				}
			}
		}

		// Validate each socket
		for _, socket := range app.Sockets {
			if err := validateSocketWithGlobalPaths(socket, globalPaths); err != nil {
				return fmt.Errorf("invalid definition of socket %q: %v", socket.Name, err)
			}
		}
	}
	return nil
}

// validateSocketWithGlobalPaths validates a socket path against both standard
// snap directories and a list of allowed global paths.
func validateSocketWithGlobalPaths(socket *SocketInfo, globalPaths []string) error {
	if socket.ListenStream == "" {
		return fmt.Errorf("\"listen-stream\" is not defined")
	}

	address := socket.ListenStream

	// Skip validation for non-path addresses (network addresses, abstract sockets)
	switch address[0] {
	case '/', '$':
		// Path-based socket, validate it
	case '@':
		// Abstract socket, use existing validation
		return validateSocketAddrAbstract(socket, "listen-stream", address)
	default:
		// Network address, use existing validation
		return validateSocketAddrNet(socket, "listen-stream", address)
	}

	// Validate path format
	if clean := filepath.Clean(address); clean != address {
		return fmt.Errorf("invalid \"listen-stream\": %q should be written as %q", address, clean)
	}

	// Check against standard snap directories first
	if isStandardSnapPath(socket, address) {
		return nil
	}

	// Check against global paths provided by interfaces
	if isAllowedGlobalPath(address, globalPaths) {
		return nil
	}

	// Path not allowed
	switch socket.App.DaemonScope {
	case SystemDaemon:
		return fmt.Errorf(
			"invalid \"listen-stream\": system daemon sockets must have a prefix of $SNAP_DATA, $SNAP_COMMON, $XDG_RUNTIME_DIR, or be allowed by a connected privileged interface")
	case UserDaemon:
		return fmt.Errorf(
			"invalid \"listen-stream\": user daemon sockets must have a prefix of $SNAP_USER_DATA, $SNAP_USER_COMMON, $XDG_RUNTIME_DIR, or be allowed by a connected privileged interface")
	default:
		return fmt.Errorf("invalid \"listen-stream\": cannot validate sockets for daemon-scope %q", socket.App.DaemonScope)
	}
}

// isStandardSnapPath checks if the path follows standard snap directory patterns
func isStandardSnapPath(socket *SocketInfo, path string) bool {
	switch socket.App.DaemonScope {
	case SystemDaemon:
		return strings.HasPrefix(path, "$SNAP_DATA/") ||
			strings.HasPrefix(path, "$SNAP_COMMON/") ||
			strings.HasPrefix(path, "$XDG_RUNTIME_DIR/") ||
			strings.HasPrefix(path, "$SNAP_REAL_XDG_RUNTIME_DIR/")
	case UserDaemon:
		return strings.HasPrefix(path, "$SNAP_USER_DATA/") ||
			strings.HasPrefix(path, "$SNAP_USER_COMMON/") ||
			strings.HasPrefix(path, "$XDG_RUNTIME_DIR/") ||
			strings.HasPrefix(path, "$SNAP_REAL_XDG_RUNTIME_DIR/")
	default:
		return false
	}
}

// isAllowedGlobalPath checks if the path matches any of the allowed global paths
func isAllowedGlobalPath(path string, globalPaths []string) bool {
	// Expand environment variables for comparison
	expandedPath := expandSocketPath(path)

	for _, pattern := range globalPaths {
		if matched, err := filepath.Match(pattern, expandedPath); err == nil && matched {
			return true
		}
		// Also check direct string match for exact paths
		if pattern == expandedPath {
			return true
		}
	}
	return false
}

// expandSocketPath expands common environment variables in socket paths
// for pattern matching purposes
func expandSocketPath(path string) string {
	expanded := path
	expanded = strings.Replace(expanded, "$SNAP_REAL_XDG_RUNTIME_DIR", "/run/user/*", 1)
	return expanded
}
