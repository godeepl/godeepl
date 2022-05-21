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

// https://github.com/zu1k/deepl-api-rs

package proxyapi

import (
    "bytes"
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
    HTTClient *http.Client
    ctx        context.Context
    authKey    string
}

type ResponseBody struct {
    Data string `json:"data"`
    Code int64  `json:"code"`
}

// New get deepl free api
func New(apiURL string, ctx context.Context) (*Client, error) {
    baseURL, err := url.Parse(apiURL)
    if err != nil {
        return nil, err
    }

    if ctx == nil {
        ctx = context.Background()
    }

    return &Client{
        baseURL:    baseURL,
        HTTClient: http.DefaultClient,
        ctx:        ctx,
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

// Translate Do proxy translate
func (c *Client) Translate(requestBody base.RequestBody) (string, error) {
    requestURL := c.baseURL

    // set request body
    _body, _ := json.Marshal(requestBody)

    // new http request
    req, err := http.NewRequest(http.MethodPost, requestURL.String(), bytes.NewReader(_body))
    if err != nil {
        return "", err
    }

    // set http header
    req.Header.Set("User-Agent", "godeepl client")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Content-Type", "application/json")

    // set http context
    req.WithContext(c.ctx)

    // do http request
    resp, err := c.HTTClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    rowResponseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    switch resp.StatusCode {
    case http.StatusOK:
        var ResponseBody ResponseBody
        err = json.Unmarshal(rowResponseBody, &ResponseBody)
        return ResponseBody.Data, nil
    default:
        msg := fmt.Sprintf("request [%s], response status code [%d], response body [%v]", requestURL, resp.StatusCode, resp.Body)
        return "", errors.New(msg)
    }
}
