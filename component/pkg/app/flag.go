package app

import (
	"fmt"
	"github.com/spf13/pflag"
)

// FlagOptions abstracts an interface giving to give App's components to put their options in.
type FlagOptions interface {
	// Flags gives app's NamedFlagSets to options to put their flags in it.
	Flags(*NamedFlagSets)

	// Validate run after load options from user input.
	Validate() []error
}

// addFlags add options to command flags.
func (a *App) addFlags() {

	manager := NewNamedFlagSets()

	if a.options != nil {
		a.options.Flags(manager)
	}

	if !a.noVersion {
		// TODO: see it.
		manager.GlobalFlagSet().AddFlag(pflag.Lookup("version"))
	}

	manager.GlobalFlagSet().BoolP("help", "h", false,
		fmt.Sprintf("help for %s", a.cmd.Name()))

	rootFlagSet := a.cmd.Flags()
	rootFlagSet.SortFlags = true
	for _, f := range manager.FlagSetsMap {
		rootFlagSet.AddFlagSet(f)
	}
}

// NamedFlagSets provides sorted flagSets
// It's recommended to use NewNamedFlagSets to create an object.
type NamedFlagSets struct {
	OrderName   []string
	FlagSetsMap map[string]*pflag.FlagSet
}

func NewNamedFlagSets() *NamedFlagSets {
	return &NamedFlagSets{
		OrderName:   []string{},
		FlagSetsMap: map[string]*pflag.FlagSet{},
	}
}

// AddFlagSet 's default behavior is to skip when has the same name.
func (s *NamedFlagSets) AddFlagSet(name string) *pflag.FlagSet {
	if _, ok := s.FlagSetsMap[name]; !ok {
		s.OrderName = append(s.OrderName, name)
		s.FlagSetsMap[name] = pflag.NewFlagSet(name, pflag.ExitOnError) // pflag has a global variable, review it.
	}
	return s.FlagSetsMap[name]
}

// GlobalFlagSet provides global level flag set.
func (s *NamedFlagSets) GlobalFlagSet() *pflag.FlagSet {
	g := "global"
	if _, ok := s.FlagSetsMap[g]; !ok {
		s.OrderName = append(s.OrderName, g)
		s.FlagSetsMap[g] = pflag.NewFlagSet(g, pflag.ExitOnError)
	}
	return s.FlagSetsMap[g]
}
