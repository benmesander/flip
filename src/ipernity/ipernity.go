// Package ipernity provides a basic golang API to the ipernity photo sharing site

package ipernity

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
)

// result of ipernity auth.getfrob api call
type Authgetfrob struct {
	Auth struct {
		Frob string `json:"frob"`
	}
	// XXX: coalesce
	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}

// result of ipernity auth.getToken api call
type Authgettoken struct {
	Auth struct {
		Token       string `json:"token"`
		Permissions struct {
			Doc     string `json:"doc"`
			Blog    string `json:"blog"`
			Profile string `json:"profile"`
			Network string `json:"network"`
			Post    string `json:"post"`
		}
		User struct {
			User_id  string `json:"user_id"`
			Username string `json:"username"`
			Realname string `json:"realname"`
			Lg       string `json:"lg"`
			Is_pro   string `json:"is_pro"`
		}
	}
	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}

// a parameter to an ipernity api call
type Parameter struct {
	name  string
	value string
}

// a slice of parameters
type pslice []Parameter

const (
	tokenfile = "ipernity_auth_token"              // Filename where ipernity token is cached
)

var (
	apikey    = "" // API Key, a constant obtained from ipernity
	apisecret = "" // API Secret, a constant obtained from ipernity
	token   = ""           // api auth token, either cached locally in a file or obtained from ipernity via api call
	user_id = ""           // your ipernity user id
	HttpClient http.Client // Persistent http client
)

// Len is part of sort.Interface
func (p pslice) Len() int {
	return len(p)
}

// Less is part of sort.Interface
func (p pslice) Less(i, j int) bool {
	return p[i].name < p[j].name
}

// Swap is part of sort.Interface
func (p pslice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Sign an ipernity api request with the md5 signature
func signRequest(request string) string {
	return fmt.Sprintf("%32x", md5.Sum([]byte(request)))
}

// call ipernity auth.getFrob api
func call_auth_getFrob() (Authgetfrob, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey})
	f, err := CallApiMethod(parms, "auth.getFrob")
	if err != nil {
		return Authgetfrob{}, err
	}
	jsonresult := &Authgetfrob{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting frob: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}

	return *jsonresult, nil
}

// call ipernity auth.getToken api
func call_auth_getToken(frob string) (Authgettoken, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey}, Parameter{"frob", frob})
	f, err := CallApiMethod(parms, "auth.getToken")
	if err != nil {
		return Authgettoken{}, err
	}
	jsonresult := &Authgettoken{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting frob: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}
	user_id = jsonresult.Auth.User.User_id

	return *jsonresult, nil
}

// call an ipernity api method
func CallApiMethod(parameters pslice, method string) ([]byte, error) {
	var (
		encodedval string
		signparams string
		urlparams  string
	)

	sort.Sort(parameters)

	for _, p := range parameters {
		encodedval = url.QueryEscape(p.value)
		signparams += p.name + encodedval
		urlparams += p.name + "=" + encodedval + "&"
	}
	urlparams += "api_sig=" + signRequest(signparams+method+apisecret)

	resp, err := HttpClient.Post("http://api.ipernity.com/api/"+method+"/json", "application/x-www-form-urlencoded", bytes.NewBufferString(urlparams))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// return URL user must visit to authorize this program to talk to ipernity
func getAuthUrl(frob string) string {
	var (
		encodedval string
		signparams string
		urlparams  string
		parameters pslice
	)

	parameters = append(parameters, Parameter{"api_key", apikey}, Parameter{"frob", frob}, Parameter{"api_secret", apisecret},
		Parameter{"perm_doc", "write"}, Parameter{"perm_blog", "write"})

	sort.Sort(parameters)

	for _, p := range parameters {
		encodedval = url.QueryEscape(p.value)
		signparams += p.name + encodedval
		urlparams += p.name + "=" + encodedval + "&"
	}

	urlparams += "api_sig=" + signRequest(signparams+apisecret)
	return "http://www.ipernity.com/apps/authorize?" + urlparams
}

// log in to ipernity
func Login() error {

	apikey = os.Getenv("IPERNITY_API_KEY")
	if apikey == "" {
		apikey = "XXX" // replace with value from http://www.ipernity.com/apps/key/0
	}

	apisecret = os.Getenv("IPERNITY_API_SECRET")	
	if apisecret == "" {
		apisecret = "YYY" // replace with value from http://www.ipernity.com/apps/key/0
	}

	// see if we have a token file
	data, err := ioutil.ReadFile(tokenfile)

	if err != nil {
		frob, err := call_auth_getFrob()
		if err != nil {
			return err
		}
		fmt.Println("go to " + getAuthUrl(frob.Auth.Frob))
		fmt.Println("and grant the permissions, then press <ENTER>")
		consolereader := bufio.NewReader(os.Stdin)
		input, err := consolereader.ReadString('\n')
		input = input

		tokenjson, err := call_auth_getToken(frob.Auth.Frob)
		if err != nil {
			return err
		}
		token = tokenjson.Auth.Token
		err = ioutil.WriteFile(tokenfile, []byte(token), 0644)
		if err != nil {
			return err
		}
	} else {
		token = string(data)
	}

	return nil
}
