package pixiv

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dghubble/sling"
)

var (
	TokenVariable        string
	RefreshTokenVariable string
	_tokenDeadline       time.Time
	authHook             func(string, string, time.Time) error
)

type AccountProfileImages struct {
	Px16  string `json:"px_16x16"`
	Px50  string `json:"px_50x50"`
	Px170 string `json:"px_170x170"`
}

type Account struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Account          string `json:"account"`
	MailAddress      string `json:"mail_address"`
	IsPremium        bool   `json:"is_premium"`
	XRestrict        int    `json:"x_restrict"`
	IsMailAuthorized bool   `json:"is_mail_authorized"`

	ProfileImage AccountProfileImages `json:"profile_image_urls"`
}

type authInfo struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
	User         *Account `json:"user"`
	DeviceToken  string   `json:"device_token"`
}

type authParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

type loginResponse struct {
	Response *authInfo `json:"response"`
}
type loginError struct {
	HasError bool              `json:"has_error"`
	Errors   map[string]Perror `json:"errors"`
}
type Perror struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func genClientHash(clientTime string) string {
	h := md5.New()
	_, _ = io.WriteString(h, clientTime)
	_, _ = io.WriteString(h, utils.ClientHashSecret)
	return hex.EncodeToString(h.Sum(nil))
}

func auth(params *authParams) (*authInfo, error) {
	clientTime := time.Now().Format(time.RFC3339)
	s := sling.New().Base("https://oauth.secure.pixiv.net/").Set("User-Agent", "PixivAndroidApp/5.0.115 (Android 6.0)").Set("X-Client-Time", clientTime).Set("X-Client-Hash", genClientHash(clientTime))

	res := &loginResponse{
		Response: &authInfo{
			User: &Account{},
		},
	}
	loginErr := &loginError{
		Errors: map[string]Perror{},
	}
	_, err := s.New().Post("auth/token").BodyForm(params).Receive(res, loginErr)
	if err != nil {
		return nil, err
	}
	if loginErr.HasError {
		for k, v := range loginErr.Errors {
			return nil, fmt.Errorf("login %s error: %s", k, v.Message)
		}
	}
	TokenVariable = res.Response.AccessToken
	RefreshTokenVariable = res.Response.RefreshToken
	_tokenDeadline = time.Now().Add(time.Duration(res.Response.ExpiresIn) * time.Second)

	if authHook != nil {
		err = authHook(TokenVariable, RefreshTokenVariable, _tokenDeadline)
	}

	return res.Response, err
}

// HookAuth add a hook with (token, refreshToken, tokenDeadline) after a successful auth.
// Prividing a way to store the latest token.
func _(f func(string, string, time.Time) error) {
	authHook = f
}

// Login pixiv has deprecated login api, so this function is useless
//func _(username, password string) (*Account, error) {
//	params := &authParams{
//		GetSecureURL: 1,
//		ClientID:     utils.ClientID,
//		ClientSecret: utils.ClientSecret,
//		GrantType:    "password",
//		Username:     username,
//		Password:     password,
//	}
//	a, err := auth(params)
//	if err != nil {
//		return nil, err
//	}
//	return a.User, nil
//}

//func LoadAuth(token, refreshToken string, tokenDeadline time.Time) (*Account, error) {
//	_token = token
//	_refreshToken = refreshToken
//	_tokenDeadline = tokenDeadline
//	return refreshAuth()
//}

func RefreshAuth() bool {
	params := &authParams{
		GetSecureURL: 1,
		ClientID:     utils.ClientID,
		ClientSecret: utils.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: config.Vars.PixivRefreshToken,
	}
	if a, err := auth(params); err != nil {
		fmt.Println("refresh auth error: ", err)
		return false
	} else {
		config.VarsFile.Vipers.Set("PIXIV_TOKEN", a.AccessToken)
		config.VarsFile.SaveConfig()
		TokenVariable = a.AccessToken
		fmt.Println("refresh auth success, new token: ", a.AccessToken)
	}
	return true

}

// download image to file (use 6.0 app-api)
func download(client *http.Client, url, path, name string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("download path needed")
	}
	if name == "" {
		name = filepath.Base(url)
	}
	fullPath := filepath.Join(path, name)

	if _, err := os.Stat(fullPath); err == nil {
		return 0, nil
	}

	output, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer func(output *os.File) {
		err = output.Close()
		if err != nil {
			log.Println(err)
		}
	}(output)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Referer", API_BASE)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download failed: %s", resp.Status)
	}

	n, err := io.Copy(output, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}
