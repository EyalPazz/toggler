package cmd

import (
	"fmt"
	
	"github.com/EyalPazz/toggler/internal/toggl"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running time entry",
	Long:  `Stop the currently running time entry if one exists.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		
		entry, err := client.StopTimer(workspaces[0].ID)
		if err != nil {
			fmt.Printf("Error stopping timer: %v\n", err)
			return
		}
		
		if entry == nil {
			fmt.Println("No timer is currently running")
			return
		}
		
		fmt.Printf("Timer stopped: \"%s\" (ID: %d)\n", entry.Description, entry.ID)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}