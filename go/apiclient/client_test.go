package apiclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cyberlabsai/exemplo-integracao-api/actions"
	"github.com/cyberlabsai/exemplo-integracao-api/settings"
)

const fakeToken = "fake-token"
const apiHost = "api.host"
const apiKey = "foo"
const apiSecret = "bar"
const fakeBody = "fake_body"

func TestVerifyActionToken(t *testing.T) {
	actionToken := &actions.ActionToken{
		Token: fakeToken,
	}

	fakeClient := &FakeHTTPClient{}
	cli := &Client{
		Settings: settings.Settings{
			ApiHost:   apiHost,
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		},
		HttpClient: fakeClient,
	}

	res, err := cli.VerifyActionToken(actionToken)
	if err != nil {
		t.Errorf(err.Error())
	}

	body := cli.ResponseBody(res)
	if body != fakeBody {
		t.Errorf("Expecting %s as body response. Got %s", fakeBody, body)
	}
}

type FakeHTTPClient struct{}

func (fake FakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	expectedAuth := "Basic " + apiKey + ":" + apiSecret
	currentAuth := req.Header.Get("authorization")

	if currentAuth != expectedAuth {
		return &http.Response{}, fmt.Errorf("Missing or wrong authorization header. Expected: %s . Got: %s", expectedAuth, currentAuth)
	}

	token := req.Header.Get("token")
	if token != fakeToken {
		return &http.Response{}, fmt.Errorf("Missing or wrong token header. Expected: %s . Got: %s", fakeToken, token)
	}

	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(fakeBody)),
	}, nil
}
