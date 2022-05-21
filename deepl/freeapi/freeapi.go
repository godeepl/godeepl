/*
Copyright © 2022 xiexianbin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// https://www.deepl.com/zh/docs-api/translating-text/example/

package freeapi

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/godeepl/godeepl/deepl/base"
    "io/ioutil"
    "net/http"
    "net/url"
)

// Client deepl http client
type Client struct {
    baseURL    *url.URL
    HTTPClient *http.Client
    ctx        context.Context
    authKey    string
}

type Translation struct {
    DetectedSourceLanguage string `json:"detected_source_language"`
    Text                   string `json:"text"`
}

type ResponseBody struct {
    Translations []Translation `json:"translations"`
}

// New get deepl free api
func New(authKey string, ctx context.Context) (*Client, error) {
    rowURL := fmt.Sprintf("https://%s/%s/translate", DomainFree, APIVersion)
    baseURL, err := url.Parse(rowURL)
    if err != nil {
        return nil, err
    }

    if ctx == nil {
        ctx = context.Background()
    }

    return &Client{
        baseURL:    baseURL,
        HTTPClient: http.DefaultClient,
        ctx:        ctx,
        authKey:    authKey,
    }, nil
}

func (c *Client)  SetApiURL(baseURL string) error {
    if baseURL != "" {
        _baseURL, err := url.Parse(baseURL)
        if err != nil {
            return err
        }
        c.baseURL = _baseURL
    }
    return nil
}

// Translate Do translate
func (c *Client) Translate(requestBody base.RequestBody) (string, error) {
    requestURL := c.baseURL

    // set request query
    query := requestURL.Query()
    query.Add("auth_key", c.authKey)
    query.Add("text", requestBody.Text)
    query.Add("source_lang", string(requestBody.SourceLang))
    query.Add("target_lang", string(requestBody.TargetLang))
    requestURL.RawQuery = query.Encode()

    // new http request
    req, err := http.NewRequest(http.MethodPost, requestURL.String(), nil)
    if err != nil {
        return "", err
    }

    // set http header
    req.Header.Set("User-Agent", "godeepl client")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // set http context
    req.WithContext(c.ctx)

    // do http request
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    switch resp.StatusCode {
    case http.StatusOK:
        var responseBody ResponseBody
        err = json.Unmarshal(body, &responseBody)
        return responseBody.Translations[0].Text, nil
    default:
        return "", errors.New(fmt.Sprintf("request [%s], response status code [%d], response body [%s]", requestURL, resp.StatusCode, body))
    }
}
