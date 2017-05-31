package jira

import (
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"encoding/json"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"fmt"
	"errors"
	"net/url"
	"net/http"
)

type Client struct {
	*oauth.Client
	host string
}

type AuthorizedClient struct {
	*Client
	userCredentials *oauth.Credentials
}

var FIELDS = []string{"status", "duedate", "customfield_22510", "customfield_22511", "assignee", "votes",
  "summary", "description", "attachment"}

func NewClient(host string) (*Client, error) {
	credentials, err := loadCredentials()
	if err != nil {
		return nil, fmt.Errorf("Error reading credentials, %v", err)
	}

	key, err := loadKey()
	if err != nil {
		return nil, fmt.Errorf("Could not load key: %v", err)
	}

	var client = &Client{
		Client: &oauth.Client{
			TemporaryCredentialRequestURI: buildURL("https", host, "/plugins/servlet/oauth/request-token"),
			ResourceOwnerAuthorizationURI: buildURL("https", host, "/plugins/servlet/oauth/authorize"),
			TokenRequestURI:               buildURL("https", host, "/plugins/servlet/oauth/access-token"),
		},
	}

	client.Credentials = *credentials
	client.PrivateKey = key
	client.SignatureMethod = oauth.RSASHA1
	client.host = host

	return client, nil
}

func buildURL(scheme string, host string, path string) string {
	result := url.URL{
		Scheme: scheme,
		Host: host,
		Path: path,
	}
	return result.String()
}

func loadCredentials() (*oauth.Credentials, error) {
	b, err := ioutil.ReadFile("etc/config.json")
	if err != nil {
		return nil, err
	}
	credetials := new(oauth.Credentials)

	if err = json.Unmarshal(b, credetials); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall credentials: %v",err)
	}

	return credetials, nil

}

func parseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}

func loadKey() (*rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile("etc/key")
	if err != nil {
		return nil, fmt.Errorf("Error reading key: %v", err)
	}
	key, err := parseKey(file)
	if err != nil {
		return nil, fmt.Errorf("Error parsing key: %v", err)
	}
	return key, nil
}

type httpCall func(client *http.Client, credentials *oauth.Credentials, urlStr string, form url.Values) (*http.Response, error)

func NewAuthorizedClient(client *Client, credentials *oauth.Credentials) *AuthorizedClient {
	result := new(AuthorizedClient)
	result.Client = client
	result.userCredentials = credentials
	return result
}

func (authorizedClient AuthorizedClient) GetMyself() (*User, error) {
  result := new(User)

  resp, err := authorizedClient.Get(
    nil,
    authorizedClient.userCredentials,
    buildURL("https", authorizedClient.host, "/rest/api/2/myself"),
    nil,
  )

  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    p, _ := ioutil.ReadAll(resp.Body)
    return nil, fmt.Errorf("GET %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
  }

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    return nil, err
  }

  json.Unmarshal(body, result)

  return result, nil
}

func (authorizedClient AuthorizedClient) GetTalk(key string) (*Talk, error) {
  result := new(Talk)

  resp, err := authorizedClient.Get(
    nil,
    authorizedClient.userCredentials,
    buildURL("https", authorizedClient.host, fmt.Sprintf("/rest/api/2/issue/%s", key)),
    url.Values{
      "fields": FIELDS,
    },
  )

  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    p, _ := ioutil.ReadAll(resp.Body)
    return nil, fmt.Errorf("GET %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
  }

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    return nil, err
  }

  json.Unmarshal(body, result)

  return result, nil
}

func (authorizedClient AuthorizedClient) GetTalks(query string) (*Talks, error) {
	result := new(Talks)

	resp, err := authorizedClient.Get(
		nil,
		authorizedClient.userCredentials,
		buildURL("https", authorizedClient.host, "/rest/api/2/search"),
		url.Values{
			"jql": {query},
			"fields": FIELDS,
      "maxResults": {"10000000"},
		},
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, result)

	return result, nil
}

func (authorizedClient AuthorizedClient) setVotedFlag(talk string, voted bool) error {
	call := authorizedClient.Delete

	if voted {
		call = authorizedClient.Post
	}

	resp, err := call(
		nil,
		authorizedClient.userCredentials,
		buildURL("https", authorizedClient.host, fmt.Sprintf("/rest/api/2/issue/%s/votes", talk)),
		url.Values{},
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("post %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}

	return nil
}

func (authorizedClient AuthorizedClient) Like(talk string) error {
	return authorizedClient.setVotedFlag(talk, true)
}
func (authorizedClient AuthorizedClient) Dislike(talk string) error {
	return authorizedClient.setVotedFlag(talk, false)
}
