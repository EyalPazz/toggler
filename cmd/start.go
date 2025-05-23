package cmd

import (
	"fmt"
	
	"github.com/EyalPazz/toggler/internal/toggl"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new time entry",
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")
		token := viper.GetString("api_token")
		err := toggl.StartTimer(token, description)
		if err != nil {
			fmt.Println("Error starting timer:", err)
		} else {
			fmt.Println("Timer started.")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("description", "d", "Working", "Description for the time entry")
}
