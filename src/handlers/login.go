package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type LoginFetcher struct{}

type Session struct {
	PostResponse struct {
		StatusCode int `json:"status_code"`
		Lookup     struct {
			Identifier string `json:"identifier"`
			Digest     string `json:"digest"`
		} `json:"lookup"`
	} `json:"postResponse"`
	PassResponse struct {
		StatusCode int `json:"status_code"`
	} `json:"passResponse"`
	Cookies string `json:"Cookies"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type LoginResponse struct {
	Authenticated bool                   `json:"authenticated"`
	Session       map[string]interface{} `json:"session"`
	Lookup        any                    `json:"lookup"`
	Cookies       string                 `json:"cookies"`
	Status        int                    `json:"status"`
	Message       any                    `json:"message"`
	Errors        []string               `json:"errors"`
}

func (lf *LoginFetcher) Logout(token string) (map[string]interface{}, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://academia.srmist.edu.in/accounts/p/10002227248/logout?servicename=ZohoCreator&serviceurl=https://academia.srmist.edu.in")
	req.Header.SetMethod("GET")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Referer", "https://academia.srmist.edu.in/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cookie", token)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}

	bodyText := resp.Body()

	result := map[string]interface{}{
		"status": resp.StatusCode(),
		"result": string(bodyText),
	}
	return result, nil
}

func (lf *LoginFetcher) Login(username, password string) (*LoginResponse, error) {
	user := strings.Replace(username, "@srmist.edu.in", "", 1)

	url := fmt.Sprintf("https://academia.srmist.edu.in/accounts/p/40-10002227248/signin/v2/lookup/%s@srmist.edu.in", user)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-zcsrf-token", "iamcsrcoo=3c59613cb190a67effa5b17eaba832ef1eddaabeb7610c8c6a518b753bc73848b483b007a63f24d94d67d14dda0eca9f0c69e027c0ebd1bb395e51b2c6291d63")
	req.Header.Set("cookie", "npfwg=1; npf_r=; npf_l=www.srmist.edu.in; npf_u=https://www.srmist.edu.in/faculty/dr-g-y-rajaa-vikhram/; zalb_74c3a1eecc=44130d4069ebce16724b1740d9128cae; ZCNEWUIPUBLICPORTAL=true; zalb_f0e8db9d3d=93b1234ae1d3e88e54aa74d5fbaba677; stk=efbb3889860a8a5d4a9ad34903359b4e; zccpn=3c59613cb190a67effa5b17eaba832ef1eddaabeb7610c8c6a518b753bc73848b483b007a63f24d94d67d14dda0eca9f0c69e027c0ebd1bb395e51b2c6291d63; zalb_3309580ed5=2f3ce51134775cd955d0a3f00a177578; CT_CSRF_TOKEN=9d0ab1e6-9f71-40fd-826e-7229d199b64d; iamcsr=3c59613cb190a67effa5b17eaba832ef1eddaabeb7610c8c6a518b753bc73848b483b007a63f24d94d67d14dda0eca9f0c69e027c0ebd1bb395e51b2c6291d63; _zcsr_tmp=3c59613cb190a67effa5b17eaba832ef1eddaabeb7610c8c6a518b753bc73848b483b007a63f24d94d67d14dda0eca9f0c69e027c0ebd1bb395e51b2c6291d63; npf_fx=1; _ga_QNCRQG0GFE=GS1.1.1737645192.5.0.1737645194.58.0.0; TS014f04d9=0190f757c98d895868ec35d391f7090a39080dd8e7be840ed996d7e2827e600c5b646207bb76666e56e22bfaf8d2c06ec3c913fe80; cli_rgn=IN; JSESSIONID=E78E4C7013F0D931BD251EBA136D57AE; _ga=GA1.3.1900970259.1737341486; _gid=GA1.3.1348593805.1737687406; _gat=1; _ga_HQWPLLNMKY=GS1.3.1737687405.1.0.1737687405.0.0.0")
	req.Header.Set("Referer", "https://academia.srmist.edu.in/accounts/p/10002227248/signin?hide_fp=true&servicename=ZohoCreator&service_language=en&css_url=/49910842/academia-academic-services/downloadPortalCustomCss/login&dcc=true&serviceurl=https%3A%2F%2Facademia.srmist.edu.in%2Fportal%2Facademia-academic-services%2FredirectFromLogin")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	req.SetBody([]byte("mode=primary&cli_time=1737687406853&servicename=ZohoCreator&service_language=en&serviceurl=https%3A%2F%2Facademia.srmist.edu.in%2Fportal%2Facademia-academic-services%2FredirectFromLogin"))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		fmt.Println("ERR", err)
		return nil, err
	}

	if errors, ok := data["errors"].([]interface{}); ok && len(errors) > 0 {
		lookupMsg := errors[0].(map[string]interface{})["message"].(string)
		statusCode := int(data["status_code"].(float64))

		if statusCode == 400 {
			return &LoginResponse{
				Authenticated: false,
				Session:       nil,
				Lookup:        nil,
				Cookies:       "",
				Status:        statusCode,
				Message: func() string {
					if strings.Contains(data["message"].(string), "HIP") {
						return ">_ Captcha required, We don't support yet"
					}
					return data["message"].(string)
				}(),
				Errors: []string{lookupMsg},
			}, nil
		}
	}

	exists := strings.Contains(data["message"].(string), "User exists")

	if !exists {
		return &LoginResponse{
			Authenticated: false,
			Session:       nil,
			Lookup:        nil,
			Cookies:       "",
			Status:        int(data["status_code"].(float64)),
			Message: func() string {
				if strings.Contains(data["message"].(string), "HIP") {
					return data["localized_message"].(string)
				}
				return data["message"].(string)
			}(),
			Errors: nil,
		}, nil
	}

	lookup, ok := data["lookup"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid lookup data")
	}

	session, err := lf.GetSession(password, lookup)
	if err != nil {
		return nil, err
	}

	sessionBody := map[string]interface{}{
		"success": true,
		"code":    session["passwordauth"].(map[string]interface{})["code"],
		"message": session["message"],
	}

	if strings.Contains(strings.ToLower(session["message"].(string)), "invalid") || strings.Contains(session["cookies"].(string), "undefined") {
		sessionBody["success"] = false
		return &LoginResponse{
			Authenticated: false,
			Session:       sessionBody,
			Lookup: map[string]string{
				"identifier": lookup["identifier"].(string),
				"digest":     lookup["digest"].(string),
			},
			Cookies: session["cookies"].(string),
			Status:  int(data["status_code"].(float64)),
			Message: session["message"].(string),
			Errors:  nil,
		}, nil
	}

	return &LoginResponse{
		Authenticated: true,
		Session:       sessionBody,
		Lookup:        lookup,
		Cookies:       session["cookies"].(string),
		Status:        int(data["status_code"].(float64)),
		Message:       data["message"],
		Errors:        nil,
	}, nil
}

func (lf *LoginFetcher) GetSession(password string, lookup map[string]interface{}) (map[string]interface{}, error) {
	identifier := lookup["identifier"].(string)
	digest := lookup["digest"].(string)
	body := fmt.Sprintf(`{"passwordauth":{"password":"%s"}}`, password)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("https://academia.srmist.edu.in/accounts/p/40-10002227248/signin/v2/primary/%s/password?digest=%s&cli_time=1713533853845&servicename=ZohoCreator&service_language=en&serviceurl=https://academia.srmist.edu.in/portal/academia-academic-services/redirectFromLogin", identifier, digest))
	req.Header.SetMethod("POST")
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("x-zcsrf-token", "iamcsrcoo=884b99c7-829b-4ddf-8344-ce971784bbe8")
	req.Header.Set("cookie", "f0e8db9d3d=7ad3232c36fdd9cc324fb86c2c0a58ad; bdb5e23bb2=3fe9f31dcc0a470fe8ed75c308e52278; zccpn=221349cd-fad7-4b4b-8c16-9146078c40d5; ZCNEWUIPUBLICPORTAL=true; cli_rgn=IN; iamcsr=884b99c7-829b-4ddf-8344-ce971784bbe8; _zcsr_tmp=884b99c7-829b-4ddf-8344-ce971784bbe8; 74c3a1eecc=d06cba4b90fbc9287c4162d01e13c516;")
	req.SetBody([]byte(body))

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		println("SESSIONERR", err)
		return nil, err
	}

	cookies := resp.Header.Peek("Set-Cookie")
	data["cookies"] = string(cookies)
	return data, nil
}

func (lf *LoginFetcher) Cleanup(cookie string) (int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://academia.srmist.edu.in/accounts/p/10002227248/webclient/v1/account/self/user/self/activesessions")
	req.Header.SetMethod("DELETE")
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("x-zcsrf-token", "iamcsrcoo=8cbe86b2191479b497d8195837181ee152bcfd3d607f5a15764130d8fd8ebef9d8a22c03fd4e418d9b4f27a9822f9454bb0bf5694967872771e1db1b5fbd4585")
	req.Header.Set("Referer", "https://academia.srmist.edu.in/accounts/p/10002227248/announcement/sessions-reminder?servicename=ZohoCreator&serviceurl=https://academia.srmist.edu.in/portal/academia-academic-services/redirectFromLogin&service_language=en")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	req.Header.Set("cookie", cookie)

	if err := fasthttp.Do(req, resp); err != nil {
		return 0, err
	}

	return resp.StatusCode(), nil
}
