package charm

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
)

func RenderProjectTable(projects []openapi_chaos_client.ProjectResponse) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	columns := []table.Column{
		{Title: "ID", Width: 36},
		{Title: "Org ID", Width: 36},
		{Title: "Name", Width: 20},
		{Title: "Date Created", Width: 16},
	}

	var rows []table.Row
	for _, proj := range projects {
		date, err := time.Parse(time.RFC3339Nano, proj.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		formattedDate := date.Format("2006-01-02 15:04")

		row := table.Row{
			proj.Id,
			proj.OrgId,
			proj.Name,
			formattedDate,
		}
		rows = append(rows, row)
	}

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

	return t.View()
}

func RenderOrgTable(orgs []openapi_chaos_client.OrganisationResponse) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	columns := []table.Column{
		{Title: "Org ID", Width: 36},
		{Title: "Name", Width: 20},
		{Title: "Date Created", Width: 16},
		{Title: "User ID", Width: 20},
	}

	var rows []table.Row
	for _, org := range orgs {
		date, err := time.Parse(time.RFC3339Nano, org.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		formattedDate := date.Format("2006-01-02 15:04")

		row := table.Row{
			org.Id,
			org.Name,
			formattedDate,
			org.UserId,
		}
		rows = append(rows, row)
	}

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

	return t.View()
}

func RenderSecretsTable(secrets []openapi_chaos_client.SecretMetaResponse) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Type", Width: 60},
		{Title: "Revison", Width: 10},
		{Title: "Date Created", Width: 16},
	}

	var rows []table.Row
	for _, secret := range secrets {
		date, err := time.Parse(time.RFC3339Nano, secret.Date.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		formattedDate := date.Format("2006-01-02 15:04")

		certSNIName := ""
		certExpiry := ""
		if secret.Type == "certificate" {
			certPEM := secret.Payload.SecretMetaResponseCertificatePayload.Certificate
			certBlock, _ := pem.Decode([]byte(certPEM))

			if certBlock == nil {
				break
			}

			cert := certBlock.Bytes
			x509Cert, err := x509.ParseCertificate(cert)
			if err != nil {
				break
			}

			certSNIName = x509Cert.Subject.CommonName
			certExpiry = x509Cert.NotAfter.Format("2006-01-02 15:04")
		}

		secretType := ""
		if certSNIName != "" {
			secretType = fmt.Sprintf("%s (%s, %s)", secret.Type, certSNIName, certExpiry)
		} else {
			secretType = string(secret.Type)
		}

		row := table.Row{
			secret.Name,
			secretType,
			string(secret.Revision),
			formattedDate,
		}
		rows = append(rows, row)
	}

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

	return t.View()
}

func RenderFuncTable(functions []openapi_chaos_client.Function) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 2)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 2)

	columns := []table.Column{
		{Title: "Name", Width: 17},
		{Title: "Image", Width: 18},
		{Title: "Description", Width: 17},
		{Title: "Secrets", Width: 15},
	}

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
		{Title: "Verification Status", Width: 10},
		{Title: "Certificate", Width: 30},
		{Title: "State", Width: 20},
	}

	var rows []table.Row
	for _, host := range hosts {
		certName := "None"
		if host.Certificate != nil && *host.Certificate != "" {
			certRefParts := strings.Split(*host.Certificate, "/")
			certName = certRefParts[len(certRefParts)-1]
		}

		routeable := ""
		if host.VerificationStatus != "completed" && !host.Disabled {
			routeable = " (unroutable, not verified)"
		}

		state := "Enabled"
		if host.Disabled {
			state = "Disabled"
		}

		row := table.Row{
			host.Host,
			string(host.VerificationStatus),
			certName,
			fmt.Sprintf("%s%s", state, routeable),
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(hosts)),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Cell = cellStyle
	t.SetStyles(s)

	return t.View()
}

func RenderProviderTable(providers []openapi_chaos_client.Provider) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 0)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 0)

	columns := []table.Column{
		{Title: "Name", Width: 34},
		{Title: "Countries", Width: 20},
		{Title: "Cities", Width: 20},
		{Title: "Continents", Width: 20},
	}

	var rows []table.Row
	for _, provider := range providers {
		row := table.Row{
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Render(provider.Name),
			strings.Join(provider.Locations.Countries, ", "),
			strings.Join(provider.Locations.Cities, ", "),
			strings.Join(provider.Locations.Continents, ", "),
		}
		rows = append(rows, row)
	}

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
