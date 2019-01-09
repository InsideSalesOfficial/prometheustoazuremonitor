package azuremonitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const strURL = "https://login.microsoftonline.com/##TENANTID##/oauth2/token"
const strResourceName = "https://monitoring.azure.com/"

type api struct {
	url          string
	tenantID     string
	clientID     string
	clientSecret string
	token        *Token
}

//New creates the azuremonitor api client
func New(tenantID, clientID, clientSecret string) api {
	api := api{
		strings.Replace(strURL, "##TENANTID##", tenantID, 1),
		tenantID,
		clientID,
		clientSecret,
		nil}
	return api
}

//GetAccessToken Generates a token to access the resource
func (api api) GetAccessToken() (Token, error) {
	if api.token == nil || api.token.IsExpired() {
		body := url.Values{}
		body.Set("grant_type", "client_credentials")
		body.Add("resource", strResourceName)
		body.Add("client_id", api.clientID)
		body.Add("client_secret", api.clientSecret)

		client := &http.Client{}
		r, err := http.NewRequest("POST", api.url, strings.NewReader(body.Encode())) // URL-encoded payload
		if err != nil {
			return Token{}, err
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(r)
		if err != nil {
			return Token{}, err
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var t, errToken = UnmarshalToken(bodyBytes)
		if errToken != nil {
			panic(err)
		}
		api.token = &t
		return *api.token, nil
	}
	return *api.token, nil
}

// SaveCustomAzureData Save data to the Azure Monitor Custom API
func (api api) SaveCustomAzureData(region, resourceID, postData string) (int, error) {

	urlStr := fmt.Sprintf("https://%s.monitoring.azure.com%s/metrics", region, resourceID)
	accessToken, err := api.GetAccessToken()
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(postData)) // URL-encoded payload
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)

	resp, err := client.Do(r)

	if err != nil {
		return 0, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(fmt.Sprintf("URL: %s \n AccessToken: %s, StatusCode: %d \n ResponseBody: %s", urlStr, accessToken.AccessToken, resp.StatusCode, string(bodyBytes)))

	return resp.StatusCode, nil
}
