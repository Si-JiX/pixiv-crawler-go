package request

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/input"
	"math/rand"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

//func genClientHash(clientTime string) string {
//	h := md5.New()
//	_, _ = io.WriteString(h, clientTime)
//	_, _ = io.WriteString(h, utils.ClientHashSecret)
//	return hex.EncodeToString(h.Sum(nil))
//}

//func RefreshAuth() bool {
//	req := &Request{Params: url.Values{}, Header: map[string]string{}}
//	req.AddParams("get_secure_url", "1")
//	req.AddParams("client_id", utils.ClientID)
//	req.AddParams("client_secret", utils.ClientSecret)
//	req.AddParams("grant_type", "refresh_token")
//	req.AddParams("refresh_token", config.Vars.PixivRefreshToken)
//	clientTime := time.Now().Format(time.RFC3339)
//	req.AddHeader("X-Client-Time", clientTime)
//	req.AddHeader("X-Client-Hash", genClientHash(clientTime))
//	accessToken := &AccessToken{}
//	//Post("https://oauth.secure.pixiv.net/auth/token", req).Json(accessToken)
//
//	fmt.Println(Post("https://oauth.secure.pixiv.net/auth/token", req).Text())
//	if accessToken.AccessToken == "" {
//		fmt.Println("refresh auth error  ", accessToken.AccessToken)
//		return false
//	} else {
//		config.VarsFile.Vipers.Set("PIXIV_TOKEN", accessToken.AccessToken)
//		config.VarsFile.SaveConfig()
//		fmt.Println("refresh auth success, new token: ", accessToken.AccessToken)
//	}
//	return true
//
//}

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

func get_pixiv_login_url() (string, string) {
	codeVerifier, codeChallenge := oauthPkce()
	urlValues := url.Values{
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"client":                {"pixiv-android"},
	}
	return codeVerifier, "https://app-api.pixiv.net/web/v1/login" + "?" + urlValues.Encode()
}

func loginPixiv(codeVerifier, code string) (*AccessToken, error) {
	req := &Request{
		Mode:   "POST",
		Params: url.Values{},
		Header: map[string]string{},
		Path:   "https://oauth.secure.pixiv.net/auth/token",
	}
	req.AddParams("client_id", "MOBrBDS8blbauoSck0ZfDbtuzpyT")
	req.AddParams("client_secret", "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj")
	req.AddParams("code", code)
	req.AddParams("code_verifier", codeVerifier)
	req.AddParams("grant_type", "authorization_code")
	req.AddParams("include_policy", "true")
	req.AddParams("redirect_uri", "https://app-api.pixiv.net/web/v1/users/auth/pixiv/callback")

	accessToken := &AccessToken{}
	response := Post(req)
	response.Json(accessToken)
	if accessToken.AccessToken == "" {
		return nil, fmt.Errorf("login error %s", response.Text())
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
func ChromeDriverLogin() (*AccessToken, error) {
	codeVerifier, loginURL := get_pixiv_login_url() // Get the login URL and code verifier
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

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ProfileImageUrls struct {
			Px16X16   string `json:"px_16x16"`
			Px50X50   string `json:"px_50x50"`
			Px170X170 string `json:"px_170x170"`
		} `json:"profile_image_urls"`
		ID                     string `json:"id"`
		Name                   string `json:"name"`
		Account                string `json:"account"`
		MailAddress            string `json:"mail_address"`
		IsPremium              bool   `json:"is_premium"`
		XRestrict              int    `json:"x_restrict"`
		IsMailAuthorized       bool   `json:"is_mail_authorized"`
		RequirePolicyAgreement bool   `json:"require_policy_agreement"`
	} `json:"user"`
	Response struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ProfileImageUrls struct {
				Px16X16   string `json:"px_16x16"`
				Px50X50   string `json:"px_50x50"`
				Px170X170 string `json:"px_170x170"`
			} `json:"profile_image_urls"`
			ID                     string `json:"id"`
			Name                   string `json:"name"`
			Account                string `json:"account"`
			MailAddress            string `json:"mail_address"`
			IsPremium              bool   `json:"is_premium"`
			XRestrict              int    `json:"x_restrict"`
			IsMailAuthorized       bool   `json:"is_mail_authorized"`
			RequirePolicyAgreement bool   `json:"require_policy_agreement"`
		} `json:"user"`
	} `json:"response"`
}
