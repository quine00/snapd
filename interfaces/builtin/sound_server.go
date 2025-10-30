package builtin

import (
	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/snap"
)

const soundServerInterfaceName = "sound-server"

// soundServerInterface implements the sound-server interface.
type soundServerInterface struct{}

func (i *soundServerInterface) Name() string {
	return soundServerInterfaceName
}

func (i *soundServerInterface) AutoConnect(plug *snap.PlugInfo, slot *snap.SlotInfo) bool {
    return false
}

// StaticInfo returns the static information for the interface.
func (i *soundServerInterface) StaticInfo() interfaces.StaticInfo {
	return interfaces.StaticInfo{
		Summary: `allows providing a sound server to the system`,
	}
}

// GlobalSocketActivationPaths returns host paths allowed for socket activation
// for snaps that provide a sound server.
func (i *soundServerInterface) GlobalSocketActivationPaths() []string {
	return []string{
		"/run/user/[0-9]*/pipewire-[0-9]",
		"/run/user/[0-9]*/pulse/native",
	}
}

func init() {
	registerIface(&soundServerInterface{})
}
