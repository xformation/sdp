package commands

import (
	"github.com/fatih/color"
	"github.com/xformation/sdp/pkg/cmd/sdp-cli/logger"
	s "github.com/xformation/sdp/pkg/cmd/sdp-cli/services"
)

func upgradeCommand(c CommandLine) error {
	pluginsDir := c.PluginDirectory()
	pluginName := c.Args().First()

	localPlugin, err := s.ReadPlugin(pluginsDir, pluginName)

	if err != nil {
		return err
	}

	v, err2 := s.GetPlugin(localPlugin.Id, c.RepoDirectory())

	if err2 != nil {
		return err2
	}

	if ShouldUpgrade(localPlugin.Info.Version, v) {
		s.RemoveInstalledPlugin(pluginsDir, pluginName)
		return InstallPlugin(localPlugin.Id, "", c)
	}

	logger.Infof("%s %s is up to date \n", color.GreenString("✔"), localPlugin.Id)
	return nil
}
