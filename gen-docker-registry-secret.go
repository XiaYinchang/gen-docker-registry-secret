package gendockerregistrysecrettool

import (
	"encoding/base64"
	"encoding/json"
	"k8s.io/api/core/v1"
)

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, server string) ([]byte, error) {
	dockercfgAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}

	dockerCfgJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockercfgAuth},
	}

	return json.Marshal(dockerCfgJSON)
}

func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}

func StructuredGenerate(secretName, username, password, server string) (*v1.Secret, error) {
	secret := &v1.Secret{}
	secret.Name = secretName
	secret.Type = v1.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{}
	dockercfgJSONContent, err := handleDockerCfgJSONContent(username, password, server)
	if err != nil {
		return nil, err
	}
	secret.Data[v1.DockerConfigJsonKey] = dockercfgJSONContent
	return secret, nil
}
