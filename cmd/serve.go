package cmd

import (
	"github.com/polyxia-org/morty-gateway/config"
	"github.com/polyxia-org/morty-gateway/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the gateway server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.NewConfig(cmd.Flags())
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		newServer, err := server.NewServer(config)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		newServer.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Uint16P("port", "p", 8080, "Port to listen on")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	serveCmd.Flags().StringP("controller", "c", "http://localhost:5000", "Address of the RIK controller")
	//serveCmd.MarkFlagRequired("controller")
	viper.BindPFlag("controller", serveCmd.Flags().Lookup("controller"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
