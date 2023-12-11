package cmd

import (
	"github.com/golobby/container/v3"
	"github.com/latifrons/commongo/safe_viper"
	"github.com/latifrons/goerrorcollector/consumer/core"
	"github.com/latifrons/latigo"
	"github.com/latifrons/latigo/logging"
	"github.com/latifrons/latigo/program"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(run)
}

var run = &cobra.Command{
	Use:   "clearing",
	Short: "start clearing server",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		program.LoadConfigs(program.FolderConfig{
			Root: "data",
		}, "INJ")

		lvl, err := zerolog.ParseLevel(safe_viper.ViperMustGetString("debug.log_level"))
		if err != nil {
			panic(err)
		}
		logging.SetupDefaultLoggerWithColor(lvl, !viper.GetBool("debug.log_color"))

		err = core.BuildDependencies()
		if err != nil {
			panic(err)
		}

		var componentProvider *component.ClearingComponentProvider
		err = container.Resolve(&componentProvider)
		if err != nil {
			panic(err)
		}

		engine := latigo.NewDefaultEngine()
		engine.SetupComponentProvider(componentProvider)
		engine.SetupCronJob(cronJobProvider)
		engine.SetupBootJob(bootJobProvider)
		engine.Start()
	},
}
