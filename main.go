//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("cryptomator-portable", "Cryptomator"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(utl.PathJoin(app.DataPath, "log"))
	app.Process = utl.PathJoin(app.AppPath, "Cryptomator.exe")

	log.Info().Msg("Updating configuration...")
	logDir := "../data/log"
	settingsPath := "../data/settings.json"
	ipcPortPath := "../data/ipcPort.bin"
	keychainPath := "../data/keychain.json"

	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.logDir=", "java-options=-Dcryptomator.logDir="+logDir); err != nil {
		log.Fatal().Err(err).Msg("Cannot set logDir")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.settingsPath=", "java-options=-Dcryptomator.settingsPath="+settingsPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set settingsPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.ipcPortPath=", "java-options=-Dcryptomator.ipcPortPath="+ipcPortPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set ipcPortPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.keychainPath=", "java-options=-Dcryptomator.keychainPath="+keychainPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set keychainPath")
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
