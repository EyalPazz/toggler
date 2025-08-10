package cmd

import (
	"fmt"
	
	"github.com/EyalPazz/toggler/internal/toggl"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start [description]",
	Short: "Start a new time entry",
	Long:  `Start a new time entry with an optional description. If no description is provided via flag or argument, "Working" will be used as default.`,
	Run: func(cmd *cobra.Command, args []string) {
		description, _ := cmd.Flags().GetString("description")
		
		if description == "" && len(args) > 0 {
			description = args[0]
		}
		
		if description == "" {
			description = "Working"
		}
		
		token := viper.GetString("api_token")
		if token == "" {
			fmt.Println("Error: API token not configured. Set TOGGLER_API_TOKEN environment variable or add api_token to config file.")
			return
		}
		
		client := toggl.NewClient(token)
		
		workspaces, err := client.GetWorkspaces()
		if err != nil {
			fmt.Printf("Error getting workspaces: %v\n", err)
			return
		}
		
		if len(workspaces) == 0 {
			fmt.Println("Error: No workspaces found")
			return
		}
		
		entry, err := client.StartTimer(workspaces[0].ID, description)
		if err != nil {
			fmt.Printf("Error starting timer: %v\n", err)
			return
		}
		
		fmt.Printf("Timer started: \"%s\" (ID: %d)\n", entry.Description, entry.ID)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("description", "d", "Working", "Description for the time entry")
}
