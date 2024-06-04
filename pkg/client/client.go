package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/qernal/cli-qernal/pkg/oauth"

	openapiclient "github.com/qernal/openapi-chaos-go-client"
	"golang.org/x/crypto/nacl/box"
)

var (
	hostHydra = getEnv("HOST_HYDRA", "https://hydra.qernal.dev")
	hostChaos = getEnv("HOST_CHAOS", "https://chaos.qernal.dev")
)

type QernalAPIClient struct {
	openapiclient.APIClient
}

func New(ctx context.Context, token string) (client QernalAPIClient, err error) {


	oauthClient := oauth.NewOauthClient(hostHydra)

	err = oauthClient.ExtractClientIDAndClientSecretFromToken(token)
	if err != nil {
		return QernalAPIClient{}, err
	}

	accessToken, err := oauthClient.GetAccessTokenWithClientCredentials()
	if err != nil {
		return QernalAPIClient{}, err
	}

	configuration := &openapiclient.Configuration{
		Servers: openapiclient.ServerConfigurations{
			{

				URL: hostChaos + "/v1",

			},
		},
		DefaultHeader: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
	}
	apiClient := openapiclient.NewAPIClient(configuration)

	return QernalAPIClient{
		APIClient: *apiClient,
	}, nil
}

func (qc *QernalAPIClient) FetchDek(ctx context.Context, projectID string) (*openapiclient.SecretMetaResponse, error) {
	keyRes, httpres, err := qc.SecretsAPI.ProjectsSecretsGet(ctx, projectID, "dek").Execute()
	slog.Info(httpres.Status)
	if err != nil {
		resData, httperr := ParseResponseData(httpres)
		if httperr != nil {
			return nil, fmt.Errorf("failed to fetch DEK key: unexpected HTTP error: %w", httperr)

		}
		return nil, fmt.Errorf("failed to fetch DEK key: unexpected error: %w, detail: %v", err, resData)
	}
	return keyRes, nil
}
func ParseResponseData(res *http.Response) (resData interface{}, err error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	return data, nil
}

type ResponseData struct {
	Data string `json:"data"`
}

func EncryptLocalSecret(pk, secret string) (string, error) {
	secretBytes := []byte(secret)
	pubKey, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return "", err
	}

	// Create a slice with enough capacity for both secret and public key
	privateKey := make([]byte, 0, len(secretBytes)+len(pubKey))
	privateKey = append(privateKey, secretBytes...)
	privateKey = append(privateKey, pubKey...)
	plaintextBytes := []byte(secret)

	var privateKeyArray [32]byte
	copy(privateKeyArray[:], privateKey)

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	encrypted := box.Seal(nonce[:], plaintextBytes, &nonce, &privateKeyArray, new([32]byte))

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
func getEnv(key, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		slog.Debug("falling back to default ")
	}

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

