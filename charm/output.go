package charm

import (
	"fmt"
	"strconv"
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
func RenderFuncTable(functions []openapi_chaos_client.Function) string {

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		// Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 2)

	columns := []table.Column{
		{Title: "Name", Width: 17},
		{Title: "Image", Width: 18},
		{Title: "Description", Width: 17},
		{Title: "Secrets", Width: 15},
	}

	// Prepare rows
	var rows []table.Row
	for _, function := range functions {
		secrets := strconv.Itoa(len(function.Secrets))
		row := table.Row{
			function.Name,
			function.Image,
			function.Description,
			secrets,
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(functions)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	// Render the table
	return t.View()

}

func RenderDNSTable(records map[string]string) string {

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 2)
	cellStyle := lipgloss.NewStyle().
		Padding(0, 2)

	columns := []table.Column{
		{Title: "Record Type", Width: 15},
		{Title: "Value", Width: 40},
	}

	var rows []table.Row
	for recordType, value := range records {
		row := table.Row{
			recordType,
			value,
		}
		rows = append(rows, row)
	}

	// Create and configure the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(records)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	return t.View()
}

func RenderHostTable(hosts []openapi_chaos_client.Host) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 2)
	cellStyle := lipgloss.NewStyle().
		Padding(0, 2)

	columns := []table.Column{
		{Title: "Hostname", Width: 30},
		{Title: "Status", Width: 15},
		{Title: "Certificate", Width: 20},
		{Title: "State", Width: 10},
	}

	var rows []table.Row
	for _, host := range hosts {

		certStatus := "None"
		if host.Certificate != nil && *host.Certificate != "" {
			certStatus = "Configured"
		}

		state := "Enabled"
		if host.Disabled {
			state = "Disabled"
		}

		row := table.Row{
			host.Host,
			string(host.VerificationStatus),
			certStatus,
			state,
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(hosts)),
	)

	// Apply the styles
	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	return t.View()
}
