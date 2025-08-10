package cmd

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/EyalPazz/toggler/internal/toggl"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent time entries",
	Long:  `List time entries for today or a specified date range.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("api_token")
		if token == "" {
			fmt.Println("Error: API token not configured. Set TOGGLER_API_TOKEN environment variable or add api_token to config file.")
			return
		}
		
		days, _ := cmd.Flags().GetInt("days")
		startDate := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")
		endDate := time.Now().Format("2006-01-02")
		
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
		
		entries, err := client.GetTimeEntries(workspaces[0].ID, startDate, endDate)
		if err != nil {
			fmt.Printf("Error getting time entries: %v\n", err)
			return
		}
		
		if len(entries) == 0 {
			fmt.Printf("No time entries found for the last %d day(s)\n", days)
			return
		}
		
		fmt.Printf("Time entries for the last %d day(s):\n\n", days)
		
		var totalDuration int
		for _, entry := range entries {
			start, _ := time.Parse(time.RFC3339, entry.Start)
			var duration time.Duration
			var status string
			
			if entry.Stop != "" {
				stop, _ := time.Parse(time.RFC3339, entry.Stop)
				duration = stop.Sub(start)
				status = "Completed"
			} else {
				duration = time.Since(start)
				status = "Running"
			}
			
			totalDuration += int(duration.Seconds())
			
			fmt.Printf("• %s - %s (%s)\n", 
				formatDuration(duration), 
				entry.Description, 
				status)
			fmt.Printf("  Started: %s\n", start.Format("2006-01-02 15:04:05"))
			if entry.Stop != "" {
				stop, _ := time.Parse(time.RFC3339, entry.Stop)
				fmt.Printf("  Stopped: %s\n", stop.Format("2006-01-02 15:04:05"))
			}
			fmt.Println()
		}
		
		fmt.Printf("Total time: %s\n", formatDuration(time.Duration(totalDuration)*time.Second))
	},
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	parts := []string{}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}
	
	return strings.Join(parts, " ")
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("days", "d", 1, "Number of days to look back (default: 1 - today only)")
}