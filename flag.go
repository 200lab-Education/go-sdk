package goservice

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/facebookgo/flagenv"
)

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false":
		return true
	case "":
		return true
	case "0":
		return true
	}
	return false
}

func getEnvName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	if flagenv.Prefix != "" {
		name = flagenv.Prefix + name
	}
	return strings.ToUpper(name)
}

type AppFlagSet struct {
	*flag.FlagSet
}

func newFlagSet(name string, fs *flag.FlagSet) *AppFlagSet {
	fSet := &AppFlagSet{fs}
	fSet.Usage = flagCustomUsage(name, fSet)
	return fSet
}

func (f *AppFlagSet) GetSampleEnvs() {
	f.VisitAll(func(f *flag.Flag) {
		if f.Name == "outenv" {
			return
		}

		s := fmt.Sprintf("## %s (-%s)\n", f.Usage, f.Name)
		s += fmt.Sprintf("#%s=", getEnvName(f.Name))

		if !isZeroValue(f, f.DefValue) {
			t := fmt.Sprintf("%T", f.Value)
			if t == "*flag.stringValue" {
				// put quotes on the value
				s += fmt.Sprintf("%q", f.DefValue)
			} else {
				s += fmt.Sprintf("%v", f.DefValue)
			}
		}
		fmt.Print(s, "\n\n")
	})
}

func (f *AppFlagSet) Parse(args []string) {
	flagenv.Parse()
	_ = f.FlagSet.Parse(args)
}

// inspect from PrintDefaults
func flagCustomUsage(appname string, fSet *AppFlagSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", appname)

		fSet.VisitAll(func(f *flag.Flag) {
			s := fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}
			// Boolean flags of one ASCII letter are so common we
			// treat them specially, putting their usage on the same line.
			if len(s) <= 4 { // space, space, '-', 'x'.
				s += "\t"
			} else {
				// Four spaces before the tab triggers good alignment
				// for both 4- and 8-space tab stops.
				s += "\n    \t"
			}
			s += usage
			if !isZeroValue(f, f.DefValue) {
				t := fmt.Sprintf("%T", f.Value)
				if t == "*flag.stringValue" {
					// put quotes on the value
					s += fmt.Sprintf(" (default %q)", f.DefValue)
				} else {
					s += fmt.Sprintf(" (default %v)", f.DefValue)
				}
			}
			s += fmt.Sprintf(" [$%s]", getEnvName(f.Name))
			_, _ = fmt.Fprint(os.Stderr, s, "\n")
		})
	}
}
