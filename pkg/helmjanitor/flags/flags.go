package flags

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func GetBoolPtr(flags *pflag.FlagSet, name string) (bool, bool, error) {
	found := false
	for _, arg := range os.Args {
		prefix := fmt.Sprintf("--%s", name)
		if strings.HasPrefix(arg, prefix) || strings.HasPrefix(arg, prefix+"=") {
			found = true
			break
		}
	}
	if !found {
		return false, false, nil
	}
	if val, err := flags.GetBool(name); err != nil {
	 	return false, false, fmt.Errorf("obtaining %s flag: %w", name, err)
	} else {
		return val, true, nil
	}
}
