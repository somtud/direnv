package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// CmdAllow is `direnv allow [PATH_TO_RC]`
var CmdAllow = &Cmd{
	Name:   "allow",
	Desc:   "Grants direnv to load the given .envrc",
	Args:   []string{"[PATH_TO_RC]"},
	Action: actionWithConfig(cmdAllowAction),
}

var migrationMessage = `
Migrating the allow data to the new location

The allowed .envrc permissions used to be stored in the XDG_CONFIG_DIR. It's
better to keep that folder for user-editable configuration so the data is
being moved to XDG_DATA_HOME.
`

func cmdAllowAction(env Env, args []string, config *Config) (err error) {
	var rcPath string
	if len(args) > 1 {
		rcPath = args[1]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	if _, err = os.Stat(config.AllowDir()); os.IsNotExist(err) {
		oldAllowDir := filepath.Join(config.ConfDir, "allow")
		if _, err = os.Stat(oldAllowDir); err == nil {
			fmt.Println(migrationMessage)

			fmt.Printf("moving %s to %s\n", oldAllowDir, config.AllowDir())
			err = os.Rename(oldAllowDir, config.AllowDir())
			if err != nil {
				return
			}

			fmt.Printf("creating a symlink back from %s to %s for back-compat.\n", config.AllowDir(), oldAllowDir)
			err = os.Symlink(config.AllowDir(), oldAllowDir)
			if err != nil {
				return
			}
			fmt.Println("")
			fmt.Println("All done, have a nice day!")
		}
	}

	rc := FindRC(rcPath, config)
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Allow()
}
