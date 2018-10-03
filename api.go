package servrlinkapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type linkType int

const (
	linkDiscord linkType = iota
	linkMinecraft
	linkInvalid
)

var (
	// EndpointAPI the URL of the server running the api
	EndpointAPI = "http://go.servr.link/api/"

	EndpointDiscord = EndpointAPI + "discord/"
	// EndpointDiscordIsRegistered returns whether the supplied discord ID is registered
	EndpointDiscordIsRegistered = func(discordId string) string {
		return EndpointDiscord + "isregistered?id=" + url.QueryEscape(discordId)
	}
	// EndpointDiscordGetID gets the UUID of the minecraft account with the discord ID provided
	EndpointDiscordGetUUID = func(discordId string) string { return EndpointDiscord + "getuuid?id=" + url.QueryEscape(discordId) }

	EndpointMinecraft = EndpointAPI + "minecraft/"
	// EndpointMCIsRegistered returns whether the supplied minecraft UUID is registered
	EndpointMCIsRegistered = func(uuid string) string { return EndpointMinecraft + "isregistered?uuid=" + url.QueryEscape(uuid) }
	// EndpointMCGetID retrieves the Discord user ID associated with a Minecraft account
	EndpointMCGetID = func(uuid string) string { return EndpointMinecraft + "getid?uuid=" + url.QueryEscape(uuid) }
)

// A result returned by the api.
type Result struct {
	Success    bool
	Registered bool
	Id         string
}

func (r *Result) UnmarshalJSON(b []byte) error {

	resultRaw := struct {
		Success    bool   `json:"success"`
		Registered bool   `json:"registered"`
		Id         uint64 `json:"id"`
	}{}

	if err := json.Unmarshal(b, &resultRaw); err != nil {
		return err
	}

	// Fucking aids because int/string
	if resultRaw.Success && resultRaw.Id == 0 {
		resultRaw2 := struct {
			Id string `json:"id"`
		}{}
		if err := json.Unmarshal(b, &resultRaw2); err != nil {
			return err
		}

		r.Id = resultRaw2.Id
	} else {
		r.Id = strconv.FormatUint(resultRaw.Id, 10)
	}

	r.Registered = resultRaw.Registered
	r.Success = resultRaw.Success
	return nil
}

// ApiClient supplies a client that has a timeout set to stop infinite requests.
var ApiClient = &http.Client{
	Timeout: 3 * time.Second,
}

// DoRequest does the request with the specified endpoint
func DoRequest(endpoint string) (res Result, err error) {
	resp, err := ApiClient.Get(endpoint)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New("invalid response status code: " + strconv.Itoa(resp.StatusCode) + " (" + resp.Status + ")")
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &res)
	return
}

func getLinkType(input string) (link linkType, err error) {
	link = linkInvalid
	if strings.ContainsRune(input, '-') {
		if len(input) == 36 {
			link = linkMinecraft
		} else {
			return link, errors.New("invalid input supplied, neither ID or UUID")
		}
	} else {
		_, err = strconv.ParseUint(input, 10, 64)
		if err != nil {
			return
		}

		link = linkDiscord
	}
	return
}

// Get returns whether the supplied input is registered
// A UUID provided will return whether the minecraft account is registered
// A discord ID provided will return whether the discord account is registered
func IsRegistered(input string) (registered bool, err error) {

	link, err := getLinkType(input)
	if err != nil {
		return
	}

	var res Result
	switch link {
	case linkMinecraft:
		res, err = DoRequest(EndpointMCIsRegistered(input))
		break
	case linkDiscord:
		res, err = DoRequest(EndpointDiscordIsRegistered(input))
	default:
		return
	}

	if !res.Success {
		err = errors.New("request was not success")
	}

	registered = res.Registered
	return
}

// Get returns a the linked and opposite output of the supplied input if it exists
// A UUID provided will return the linked discord ID
// A discord ID provided will return the linked UUID
func Get(input string) (output string, err error) {

	link, err := getLinkType(input)
	if err != nil {
		return
	}

	var res Result
	switch link {
	case linkMinecraft:
		res, err = doRequest(EndpointMCGetID(input))
		break
	case linkDiscord:
		res, err = doRequest(EndpointDiscordGetUUID(input))
	default:
		return
	}

	if !res.Success {
		err = errors.New("request was not success")
	}

	output = res.Id
	return
}
