package packet

import (
	"github.com/CN-TU/go-flows/util"
	"github.com/google/gopacket"
)

const filterName = "filter"

// Filter represents a generic packet filter
type Filter interface {
	util.Module
	// Matches must return true, if this packet should be used
	Matches(lt gopacket.LayerType, data []byte, ci gopacket.CaptureInfo, n uint64) bool
}

// Filters holds a collection of filters that are tried one after another
type Filters []Filter

// Matches returns true if this packet matches all filters
func (f Filters) Matches(lt gopacket.LayerType, data []byte, ci gopacket.CaptureInfo, n uint64) bool {
	for _, filter := range f {
		if !filter.Matches(lt, data, ci, n) {
			return false
		}
	}
	return true
}

// RegisterFilter registers an filter (see module system in util)
func RegisterFilter(name, desc string, new util.ModuleCreator, help util.ModuleHelp) {
	util.RegisterModule(filterName, name, desc, new, help)
}

// FilterHelp displays help for a specific filter (see module system in util)
func FilterHelp(which string) error {
	return util.GetModuleHelp(filterName, which)
}

// MakeFilter creates an filter instance (see module system in util)
func MakeFilter(which string, args []string) ([]string, Filter, error) {
	args, module, err := util.CreateModule(filterName, which, args)
	if err != nil {
		return args, nil, err
	}
	return args, module.(Filter), nil
}

// ListFilters returns a list of filters (see module system in util)
func ListFilters() ([]util.ModuleDescription, error) {
	return util.GetModules(filterName)
}
