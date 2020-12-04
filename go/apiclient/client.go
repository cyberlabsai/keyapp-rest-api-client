package apiclient

import (
	"bytes"
	"net/http"

	"github.com/cyberlabsai/exemplo-integracao-api/settings"

	"github.com/cyberlabsai/exemplo-integracao-api/actions"
)

// Client é uma estrutura com as configurações e o HTTP client utilizados na comunicação com a API
// O field Settings contém as chaves de API necessárias para comunicação com a API do KeyApp
type Client struct {
	Settings settings.Settings
	HttpClient
}

// ApiClient contrato com funções que é esperado a API expor
type ApiClient interface {
	VerifyActionToken(*actions.ActionToken) (*http.Response, error)
	ResponseBody(*http.Response) (string, error)
}

// HttpClient é o contrato usado para fazer requisições HTTP.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// VerifyActionToken adiciona os headers de autorização e o token da ação
func (client Client) VerifyActionToken(at *actions.ActionToken) (*http.Response, error) {
	request, _ := http.NewRequest("GET", client.Settings.ApiHost+"/apptokens", nil)
	request.Header.Add("authorization", "Basic "+client.Settings.ApiKey+":"+client.Settings.ApiSecret)
	request.Header.Add("token", at.Token)
	return client.HttpClient.Do(request)
}

// ResponseBody é responsável ler o buffer no body e retornar uma string
func (client Client) ResponseBody(res *http.Response) string {
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	responseBody := buf.String()

	return responseBody
}
