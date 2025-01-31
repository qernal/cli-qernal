package charm

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
)

func RenderProjectTable(projects []openapi_chaos_client.ProjectResponse) string {
	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	// Define table columns with increased width for IDs
	columns := []table.Column{
		{Title: "ID", Width: 36},
		{Title: "Org ID", Width: 36},
		{Title: "Name", Width: 20},
		{Title: "Date Created", Width: 16},
	}

	// Prepare rows
	var rows []table.Row
	for _, proj := range projects {
		// Parse the date string
		date, err := time.Parse(time.RFC3339Nano, proj.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		// Format the date to be more human-friendly
		formattedDate := date.Format("2006-01-02 15:04")

		row := table.Row{
			proj.Id,
			proj.OrgId,
			proj.Name,
			formattedDate,
		}
		rows = append(rows, row)
	}

	// Create and style the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(projects)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	// Render the table
	return t.View()
}

func RenderOrgTable(orgs []openapi_chaos_client.OrganisationResponse) string {
	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	// Define table columns with increased width for IDs
	columns := []table.Column{
		{Title: "Org ID", Width: 36},
		{Title: "Name", Width: 20},
		{Title: "Date Created", Width: 16},
		{Title: "User ID", Width: 20},
	}

	// Prepare rows
	var rows []table.Row
	for _, org := range orgs {
		// Parse the date string
		date, err := time.Parse(time.RFC3339Nano, org.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		// Format the date to be more human-friendly
		formattedDate := date.Format("2006-01-02 15:04")

		row := table.Row{
			org.Id,
			org.Name,
			formattedDate,
			org.UserId,
		}
		rows = append(rows, row)
	}

	// Create and style the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(orgs)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	// Render the table
	return t.View()
}

func RenderSecretsTable(secrets []openapi_chaos_client.SecretMetaResponse) string {

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	// Define table columns with increased width for IDs
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 14},
		{Title: "Revison", Width: 10},
		{Title: "Date Created", Width: 16},
	}

	// Prepare rows
	var rows []table.Row
	for _, secret := range secrets {
		// Parse the date string
		date, err := time.Parse(time.RFC3339Nano, secret.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		// Format the date to be more human-friendly
		formattedDate := date.Format("2006-01-02 15:04")

		row := table.Row{
			secret.Name,
			string(secret.Type),
			string(secret.Revision),
			formattedDate,
		}
		rows = append(rows, row)
	}

	// Create and style the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(secrets)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	// Render the table
	return t.View()

}

func RenderProviderTable(providers []openapi_chaos_client.Provider) string {
	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	providerNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C8C8C8")).
		Padding(0, 1, 1)

	// Prepare rows
	var rows []table.Row
	for _, provider := range providers {
		// Join array values with commas for display
		countries := strings.Join(provider.Locations.Countries, ", ")
		cities := strings.Join(provider.Locations.Cities, ", ")
		continents := strings.Join(provider.Locations.Continents, ", ")

		row := table.Row{
			providerNameStyle.Render(provider.Name),
			countries,
			cities,
			continents,
		}
		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "Name", Width: 35},
		{Title: "Countries", Width: 35},
		{Title: "City", Width: 35},
		{Title: "Continents", Width: 15},
	}

	// Create and style the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(providers)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle

	t.SetStyles(s)

	return t.View()
}
