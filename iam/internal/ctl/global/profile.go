package global

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
)

type Profile struct {
	enabled bool
	profile string
	outFile string
}

var DefaultProfile = &Profile{}

func (p *Profile) AddFlag(fs *pflag.FlagSet) {
	fs.BoolVar(&p.enabled, "enable-profile", false, "Enable profile.")
	fs.StringVar(
		&p.profile,
		"profile-type",
		p.profile,
		"Name of profile to capture. One of (cpu|heap|goroutine|threadcreate|block|mutex)",
	)
	fs.StringVar(&p.outFile, "profile-output", "std", "Name of the file to write the profile to, if write to io stream, put 'std'.")
}

// Run launch to run profiling service.
func (p *Profile) Run() (err error) {
	if !p.enabled {
		return nil
	}

	switch p.profile {
	case "none":
		return nil
	case "cpu":
		var w io.Writer
		if p.outFile == "std" {
			w = os.Stdout
		} else {
			w, err = os.Create(p.outFile)
			if err != nil {
				return err
			}
		}
		return pprof.StartCPUProfile(w)
	case "block":
		runtime.SetBlockProfileRate(1)
		return nil
	case "mutex":
		runtime.SetMutexProfileFraction(1)
		return nil
	default:
		if pe := pprof.Lookup(p.profile); pe == nil {
			err = fmt.Errorf("profile not found, got: %s", p.profile)
		}
	}
	return
}

// Close does some clean work when going to close.
func (p *Profile) Close() (err error) {
	if !p.enabled {
		return nil
	}
	switch p.profile {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
		return
	case "heap":
		runtime.GC()
		fallthrough
	default:
		pe := pprof.Lookup(p.profile)
		if pe == nil {
			return
		}
		var w io.Writer
		if p.outFile == "std" {
			w = os.Stdout
		} else {
			w, err = os.Create(p.outFile)
			if err != nil {
				return err
			}
		}
		return pe.WriteTo(w, 0)
	}
}
