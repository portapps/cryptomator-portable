//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"encoding/json"
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
	pluginDir := "../data/plugins"
	settingsPath := "../data/settings.json"
	ipcPortPath := "../data/ipcPort.bin"
	ipcSocketPath := "../data/ipc.socket"
	keychainPath := "../data/keychain.json"
	p12Path := "../data/key.p12"

	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.logDir=", "java-options=-Dcryptomator.logDir="+logDir); err != nil {
		log.Fatal().Err(err).Msg("Cannot set logDir")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.pluginDir=", "java-options=-Dcryptomator.pluginDir="+pluginDir); err != nil {
		log.Fatal().Err(err).Msg("Cannot set pluginDir")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.settingsPath=", "java-options=-Dcryptomator.settingsPath="+settingsPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set settingsPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.ipcPortPath=", "java-options=-Dcryptomator.ipcPortPath="+ipcPortPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set ipcPortPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.ipcSocketPath=", "java-options=-Dcryptomator.ipcSocketPath="+ipcSocketPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set ipcSocketPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.integrationsWin.keychainPaths=", "java-options=-Dcryptomator.integrationsWin.keychainPaths="+keychainPath); err != nil {
		log.Fatal().Err(err).Msg("Cannot set keychainPath")
	}
	if err := utl.ReplaceByPrefix(utl.PathJoin(app.AppPath, "app", "Cryptomator.cfg"), "java-options=-Dcryptomator.p12Path=", "java-options=-Dcryptomator.p12Path="+p12Path); err != nil {
		log.Fatal().Err(err).Msg("Cannot set p12Path")
	}

	// Create folders
	_ = utl.CreateFolder(utl.PathJoin(app.DataPath, "log"))
	_ = utl.CreateFolder(utl.PathJoin(app.DataPath, "plugins"))

	// Update settings
	settingsFile := utl.PathJoin(app.DataPath, "settings.json")
	if _, err := os.Stat(settingsFile); err == nil {
		log.Info().Msg("Updating settings...")
		rawSettings, err := os.ReadFile(settingsFile)
		if err == nil {
			jsonMapSettings := make(map[string]interface{})
			if err = json.Unmarshal(rawSettings, &jsonMapSettings); err != nil {
				log.Error().Err(err).Msg("Settings unmarshal")
			}
			log.Info().Interface("settings", jsonMapSettings).Msg("Current settings")

			jsonMapSettings["askedForUpdateCheck"] = false
			jsonMapSettings["checkForUpdatesEnabled"] = false
			log.Info().Interface("settings", jsonMapSettings).Msg("New settings")

			jsonSettings, err := json.MarshalIndent(jsonMapSettings, "", "  ")
			if err != nil {
				log.Error().Err(err).Msg("Settings marshal")
			}
			if err = os.WriteFile(settingsFile, jsonSettings, 0644); err != nil {
				log.Error().Err(err).Msg("Write settings")
			}
		}
	} else {
		log.Info().Msg("Initializing settings...")
		if err = os.WriteFile(settingsFile, []byte(`{
  "askedForUpdateCheck": false,
  "checkForUpdatesEnabled": false
}`), 0644); err != nil {
			log.Error().Err(err).Msg("Write settings")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
