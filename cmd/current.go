package cmd

import (
	"fmt"
	"time"
	
	"github.com/EyalPazz/toggler/internal/toggl"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently running time entry",
	Long:  `Display information about the currently running time entry, if any.`,
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
		
		entry, err := client.GetCurrentTimer(workspaces[0].ID)
		if err != nil {
			fmt.Printf("Error getting current timer: %v\n", err)
			return
		}
		
		if entry == nil {
			fmt.Println("No timer is currently running")
			return
		}
		
		start, err := time.Parse(time.RFC3339, entry.Start)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return
		}
		
		duration := time.Since(start)
		
		fmt.Printf("Current timer:\n")
		fmt.Printf("• Description: %s\n", entry.Description)
		fmt.Printf("• Started: %s\n", start.Format("2006-01-02 15:04:05"))
		fmt.Printf("• Duration: %s\n", formatDuration(duration))
		fmt.Printf("• ID: %d\n", entry.ID)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}