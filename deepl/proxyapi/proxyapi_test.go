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

package proxyapi

import (
    "fmt"
    "github.com/godeepl/godeepl/deepl/base"
    "testing"
)

func TestNewProxyAPI(t *testing.T) {
    client, _ := New("http://47.243.244.86:8080/translate", nil)
    req := base.RequestBody{
        Text:       "hello world!",
        SourceLang: base.EN,
        TargetLang: base.ZH,
    }
    resp, err := client.Translate(req)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(fmt.Sprintf("%v", resp))
}
