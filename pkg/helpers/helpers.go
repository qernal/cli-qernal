package helpers

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	math_rand "math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
)

func CreateOrg() (string, string, error) {
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	organisationBody := *openapi_chaos_client.NewOrganisationBody(uuid.NewString())
	resp, r, err := client.OrganisationsAPI.OrganisationsCreate(context.Background()).OrganisationBody(organisationBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganisationsAPI.OrganisationsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Id, resp.Name, nil
}

func DeleteOrg(orgid string) {
	token, err := auth.GetQernalToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("obtaining token failed with:", err).Error())
	}

	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
	}
	_, r, err := client.OrganisationsAPI.OrganisationsDelete(context.Background(), orgid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
		fmt.Fprintf(os.Stderr, "Error when calling `OrganisationsAPI.OrganisationsDelete`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}

func CreateProj(orgid string) (string, string, error) {
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	projectBody := *openapi_chaos_client.NewProjectBody(orgid, uuid.NewString())
	resp, r, err := client.ProjectsAPI.ProjectsCreate(context.Background()).ProjectBody(projectBody).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Id, resp.Name, nil
}

func DeleteProj(projid string) {
	token, err := auth.GetQernalToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("obtaining token failed with:", err).Error())
	}

	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
	}

	_, r, err := client.ProjectsAPI.ProjectsDelete(context.Background(), projid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}

func GenerateSelfSignedCert() ([]byte, []byte, error) {
	// Generate a new ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Create a self-signed certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Example Corp"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create the self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode the public key (certificate) to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	// Encode the private key to PEM
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return certPEM, privateKeyPEM, nil
}

func CreateSecretEnv(projid string, secretname string) (string, string, error) {
	dek, dekRevision, err := FetchDek(projid)
	if err != nil {
		return "", "", err
	}

	encryptedSecret, err := client.EncryptLocalSecret(dek, secretname)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	secretEnvBody := *openapi_chaos_client.NewSecretBody(secretname, openapi_chaos_client.SECRETCREATETYPE_ENVIRONMENT, openapi_chaos_client.SecretCreatePayload{
		SecretEnvironment: &openapi_chaos_client.SecretEnvironment{
			EnvironmentValue: encryptedSecret,
		},
	}, fmt.Sprintf("keys/dek/%d", dekRevision))
	resp, r, err := qc.SecretsAPI.ProjectsSecretsCreate(context.Background(), projid).SecretBody(secretEnvBody).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsSecretsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Name, fmt.Sprintf("projects:%s/%s@%d", projid, resp.Name, resp.Revision), nil
}

func FetchDek(projectID string) (string, int32, error) {
	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", 0, err
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", 0, err
	}

	resp, r, err := qc.SecretsAPI.ProjectsSecretsGet(context.Background(), projectID, "dek").Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsSecretsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", 0, err
	}

	return resp.Payload.SecretMetaResponseDek.Public, resp.Revision, nil
}

func GetDefaultHost(projid string) (string, error) {
	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", err
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", err
	}
	resp, r, err := qc.HostsAPI.ProjectsHostsList(context.Background(), projid).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", err
	}

	for _, host := range resp.Data {
		if host.ReadOnly {
			return host.Host, nil
		}
	}

	return "", errors.New("no default host on project")
}

// specific to the secret type
func RandomSecretName() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[math_rand.Intn(len(charset))]
	}
	return fmt.Sprintf("TERRA_%s", string(b))
}

func PaginateProjects(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32) ([]openapi_chaos_client.ProjectResponse, error) {

	initialResp, httpRes, err := qc.ProjectsAPI.ProjectsList(ctx).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list projects, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list projects", err)
	}
	allProjects := initialResp.GetData()
	if initialResp.Meta.Results <= 20 {
		return allProjects, nil
	}

	pageSize := int32(20)
	var currentPage int32
	for currentPage < initialResp.Meta.Pages {
		if maxResults > 0 && len(allProjects) >= int(maxResults) {
			break
		}

		currentPage++
		previous := currentPage - 1
		next, httpRes, err := qc.ProjectsAPI.ProjectsList(ctx).
			Page(openapi_chaos_client.OrganisationsListPageParameter{
				Size:   &pageSize,
				Before: &previous,
				After:  &currentPage,
			}).Execute()
		if err != nil {
			resData, _ := client.ParseResponseData(httpRes)
			if data, ok := resData.(map[string]interface{}); ok {
				if innerData, ok := data["data"].(map[string]interface{}); ok {
					if nameErr, ok := innerData["name"].(string); ok {
						return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
					}
				}
			}
			printer.Logger.Debug("unable to list projects, request failed",
				slog.String("error", err.Error()),
				slog.Any("response", resData))
			return nil, printer.RenderError("unable to list projects", err)
		}

		allProjects = append(allProjects, next.GetData()...)
	}

	if maxResults > 0 && len(allProjects) > int(maxResults) {
		allProjects = allProjects[:maxResults]
	}

	return allProjects, nil
}

func PaginateOrganisations(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32) ([]openapi_chaos_client.OrganisationResponse, error) {

	initialResp, httpRes, err := qc.OrganisationsAPI.OrganisationsList(ctx).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list organisations", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list organisations, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list organisations", err)
	}

	allProjects := initialResp.GetData()
	if initialResp.Meta.Results <= 20 {
		return allProjects, nil
	}

	pageSize := int32(20)
	var currentPage int32
	for currentPage < initialResp.Meta.Pages {
		if maxResults > 0 && len(allProjects) >= int(maxResults) {
			break
		}

		currentPage++
		previous := currentPage - 1
		next, httpRes, err := qc.OrganisationsAPI.OrganisationsList(ctx).
			Page(openapi_chaos_client.OrganisationsListPageParameter{
				Size:   &pageSize,
				Before: &previous,
				After:  &currentPage,
			}).Execute()
		if err != nil {
			resData, _ := client.ParseResponseData(httpRes)
			if data, ok := resData.(map[string]interface{}); ok {
				if innerData, ok := data["data"].(map[string]interface{}); ok {
					if nameErr, ok := innerData["name"].(string); ok {
						return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
					}
				}
			}
			printer.Logger.Debug("unable to list projects, request failed",
				slog.String("error", err.Error()),
				slog.Any("response", resData))
			return nil, printer.RenderError("unable to list projects", err)
		}

		allProjects = append(allProjects, next.GetData()...)
	}

	if maxResults > 0 && len(allProjects) > int(maxResults) {
		allProjects = allProjects[:maxResults]
	}

	return allProjects, nil
}

// PaginateSecrets
func PaginateSecrets(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32, projectID string) ([]openapi_chaos_client.SecretMetaResponse, error) {

	initialResp, httpRes, err := qc.SecretsAPI.ProjectsSecretsList(ctx, projectID).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list projects, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list projects", err)
	}

	allProjects := initialResp.Data
	if initialResp.Meta.Results <= 20 {
		return allProjects, nil
	}

	pageSize := int32(20)
	var currentPage int32
	for currentPage < initialResp.Meta.Pages {
		if maxResults > 0 && len(allProjects) >= int(maxResults) {
			break
		}

		currentPage++
		previous := currentPage - 1
		next, httpRes, err := qc.SecretsAPI.ProjectsSecretsList(ctx, projectID).
			Page(openapi_chaos_client.OrganisationsListPageParameter{
				Size:   &pageSize,
				Before: &previous,
				After:  &currentPage,
			}).Execute()
		if err != nil {
			resData, _ := client.ParseResponseData(httpRes)
			if data, ok := resData.(map[string]interface{}); ok {
				if innerData, ok := data["data"].(map[string]interface{}); ok {
					if nameErr, ok := innerData["name"].(string); ok {
						return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
					}
				}
			}
			printer.Logger.Debug("unable to list projects, request failed",
				slog.String("error", err.Error()),
				slog.Any("response", resData))
			return nil, printer.RenderError("unable to list projects", err)
		}

		allProjects = append(allProjects, next.GetData()...)
	}

	if maxResults > 0 && len(allProjects) > int(maxResults) {
		allProjects = allProjects[:maxResults]
	}

	return allProjects, nil
}
