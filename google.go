package goauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type OnGoogleAuth func(authToken GoogleOAuthToken)

type GoogleOAuth struct {
	OAuth *OAuth
}

type GoogleOAuthToken struct {
	OAuthToken   `json:"-"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type googleOAuthRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

func (g *GoogleOAuth) Auth() GoogleOAuthToken {
	oAuthToken := GoogleOAuthToken{}
	if g.OAuth.error != nil {
		oAuthToken.Error = g.OAuth.error
		return oAuthToken
	}

	g.OAuth.ServerConfig.oAuthUrl = genGoogleOauthUrl(g.OAuth)

	code := fetchAuthCode(&g.OAuth.ServerConfig)
	if code == "" {
		oAuthToken.Error = AuthFailedError{Msg: "Failed to retrieve auth code."}
	} else {
		var err error
		oAuthToken, err = fetchOAuthToken(code, g.OAuth)
		oAuthToken.Error = err
	}
	return oAuthToken
}

func (g *GoogleOAuth) AsyncAuth(onGoogleAuth OnGoogleAuth) error {
	if g.OAuth.error != nil {
		return g.OAuth.error
	}

	g.OAuth.ServerConfig.oAuthUrl = genGoogleOauthUrl(g.OAuth)

	onAuth := func(authCode string) {
		token, err := fetchOAuthToken(authCode, g.OAuth)
		token.Error = err
		onGoogleAuth(token)
	}

	return fetchAuthCodeAsync(&g.OAuth.ServerConfig, onAuth)
}

func fetchOAuthToken(code string, auth *OAuth) (GoogleOAuthToken, error) {
	url := "https://www.googleapis.com/oauth2/v4/token"

	body := googleOAuthRequest{
		Code:         code,
		ClientId:     auth.ClientId,
		ClientSecret: auth.ClientSecret,
		RedirectUri:  auth.ServerConfig.redirectUri,
		GrantType:    "authorization_code",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return GoogleOAuthToken{}, nil
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return GoogleOAuthToken{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GoogleOAuthToken{}, err
	}

	var oAuthToken GoogleOAuthToken
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GoogleOAuthToken{}, err
	}
	err = json.Unmarshal(respBodyData, &oAuthToken)
	if err != nil {
		return GoogleOAuthToken{}, err
	}

	return oAuthToken, nil
}

func genGoogleOauthUrl(oAuth *OAuth) string {
	googleUrl := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=%s&" +
		"redirect_uri=%s&" +
		"response_type=code&" +
		"scope=%s&" +
		"prompt=select_account"

	redirectUri := "http://localhost:" + strconv.Itoa(oAuth.ServerConfig.Port)

	return fmt.Sprintf(googleUrl, oAuth.ClientId, redirectUri, oAuth.Scope)
}
