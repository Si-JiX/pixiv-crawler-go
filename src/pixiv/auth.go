package pixiv

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"pixiv-cil/pkg/input"
	"pixiv-cil/src/pixivstruct"
	"runtime"
	"strings"
)

// Generate a random token
func generateURLSafeToken(length int) string {

	str := "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(str[rand.Intn(len(str))])
	}
	return sb.String()
}

// S256 transformation method.
func s256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// Proof Key for Code Exchange by OAuth Public Clients (RFC7636).
func oauthPkce() (string, string) {
	codeVerifier := generateURLSafeToken(32)
	codeChallenge := s256(codeVerifier)
	return codeVerifier, codeChallenge
}

func getLoginURL() (string, string) {
	codeVerifier, codeChallenge := oauthPkce()
	urlValues := url.Values{
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"client":                {"pixiv-android"},
	}
	return codeVerifier, "https://app-api.pixiv.net/web/v1/login" + "?" + urlValues.Encode()
}

func loginPixiv(codeVerifier, code string) (*pixivstruct.AccessToken, error) {
	urlValues := url.Values{
		"client_id":      {"MOBrBDS8blbauoSck0ZfDbtuzpyT"},
		"client_secret":  {"lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"},
		"code":           {code},
		"code_verifier":  {codeVerifier},
		"grant_type":     {"authorization_code"},
		"include_policy": {"true"},
		"redirect_uri":   {"https://app-api.pixiv.net/web/v1/users/auth/pixiv/callback"},
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	req, err := http.NewRequest("POST", "https://oauth.secure.pixiv.net/auth/token", strings.NewReader(urlValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 (Android 11; Pixel 5)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accessToken *pixivstruct.AccessToken
	if err = json.Unmarshal(all, &accessToken); err != nil {
		return nil, err
	} else {
		return accessToken, nil
	}
}

func openbrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
func ChromeDriverLogin() (*pixivstruct.AccessToken, error) {
	codeVerifier, loginURL := getLoginURL() // Get the login URL and code verifier
	fmt.Println("please open the following link in your browser:", loginURL)
	fmt.Println("please press f12 to open the developer console, and switch to the network tab.")
	fmt.Println("now, please send the value of the code parameter in the request url of the remaining request.")
	fmt.Println("after logging in, please enter the code value:")
	fmt.Println("note that the code has a very short lifetime, please make sure that the previous step is completed quickly.")
	if err := openbrowser(loginURL); err != nil {
		fmt.Println("failed to open browser, please open the following link in your browser:", loginURL)
	} else {
		fmt.Printf("browser opened successfully,please input the code value:")
	}
	return loginPixiv(codeVerifier, input.Input("please input the code value:", ">"))

}
