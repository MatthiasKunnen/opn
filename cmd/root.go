package cmd

import (
	"github.com/MatthiasKunnen/opn/cmd/query"
	"github.com/spf13/cobra"
	"log"
)

var (
	cfgFile     string
	userLicense string
)

var rootCmd = &cobra.Command{
	Use:   "opn",
	Short: "opn is a terminal file opener",
	Long:  `opn `,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.AddCommand(updateDbCmd)
	rootCmd.AddCommand(openFileCmd)
	rootCmd.AddCommand(query.QueryCmd)
}
