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

package cmd

import (
	"bufio"
	"fmt"
	"github.com/godeepl/godeepl/deepl/base"
	"os"
	"strings"

	"github.com/godeepl/godeepl/deepl"
	"github.com/godeepl/godeepl/utils"

	"github.com/spf13/cobra"
)

var file string
var output string
var text string
var sourceLang string
var targetLang string
var proxyAPI string

func translateText(client deepl.Client) {
	req := base.RequestBody{
		Text:       text,
		SourceLang: base.LangCode(sourceLang),
		TargetLang: base.LangCode(targetLang),
	}
	resp, err := client.Translate(req)
	if err != nil {
		fmt.Println("translate error:", err.Error())
		os.Exit(1)
	}

	if output == "" {
		fmt.Println(fmt.Sprintf("%v", resp))
	} else {
		_ = utils.Write(output, []byte(resp), 0644)
	}
}

func translateFile(client deepl.Client) {
	sourceContentItems, err := utils.ReadMarkdown(file)
	if err != nil {
		fmt.Println("read from file", file, "error", err.Error())
	}

	stringChan := make(chan string)
	quit := make(chan bool)

	if output == "" {
		output = strings.Replace(file, ".md", ".deepl.md", 1)
	}

	go func() {
		fmt.Println("begin to write result to", output)
		var file *os.File
		var err error
		if utils.IsFileExist(output) {
			file, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else {
			file, err = os.Create(output)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		defer file.Close()

		w := bufio.NewWriter(file)

		// channel
		for {
			select {
			case v := <-stringChan:
				n, err := w.WriteString(v)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				fmt.Printf("write %d bytes.\n", n)

				err = w.Flush()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			case <-quit:
				fmt.Println("write goroutine exit.")
				return
			}
		}
	}()

	for _, sourceContentItem := range sourceContentItems {
		if sourceContentItem.Type == utils.Code {
			stringChan <- strings.Join(sourceContentItem.Content, "\n") + "\n"
			continue
		} else if sourceContentItem.Type == utils.Text {
			if len(sourceContentItem.Content) == 1 && sourceContentItem.Content[0] == "" {
				stringChan <- "\n"
			} else if len(sourceContentItem.Content) == 1 && strings.HasPrefix(sourceContentItem.Content[0], "![](") {
				stringChan <- strings.Join(sourceContentItem.Content, "\n") + "\n"
			} else {
				var targetContent []string
				for _, i := range sourceContentItem.Content {
					req := base.RequestBody{
						Text:       i,
						SourceLang: base.LangCode(sourceLang),
						TargetLang: base.LangCode(targetLang),
					}
					resp, err := client.Translate(req)
					if err != nil {
						fmt.Println("translate error:", err.Error())
						targetContent = append(targetContent, i)
					} else {
						fmt.Println("source text:", i, ", target text:", resp)
						targetContent = append(targetContent, resp)
					}

					// rand sleep for api limit
					utils.Sleep()
				}
				stringChan <- strings.Join(targetContent, "\n") + "\n"
			}
		}
	}
	quit <- true
	fmt.Println("translator finished.")
}

// translatorCmd represents the translator command
var translatorCmd = &cobra.Command{
	Use:   "translator",
	Short: "translate text or file to target lang.",
	Long: `Translate text or file to target lang. For example:

# use proxy api
godeepl translator --text "hello world!" -s EN -t ZH -p http://127.0.0.1/v2/translate

# use Deepl free API
godeepl translator --text "hello world!" -s EN -t ZH

# from file
godeepl translator --file "/<path>/<file>.md" -s EN -t ZH
`,
	Run: func(cmd *cobra.Command, args []string) {
		if text == "" && file == "" {
			fmt.Println("text or file must exits one.")
			os.Exit(1)
		}

		if file != "" {
			if utils.IsFileExist(file) == false {
				fmt.Println("file", file, "is not exist.")
				os.Exit(1)
			}
		}

		if sourceLang == "" && targetLang == "" {
			fmt.Println("sourceLang and targetLang must exits.")
			os.Exit(1)
		}

		var client deepl.Client
		if proxyAPI != "" {
			client = deepl.Factory(base.ProxyAPI, "", proxyAPI)
		} else {
			client = deepl.Factory(base.FreeAPI, "", "")
		}

		if text != "" {
			translateText(client)
		} else if file != "" {
			translateFile(client)
		}
	},
}

func init() {
	rootCmd.AddCommand(translatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// translatorCmd.PersistentFlags().String("foo", "", "A help for foo")
	translatorCmd.Flags().StringVarP(&text, "text", "c", "", "text to be translate.")
	translatorCmd.Flags().StringVarP(&file, "file", "f", "", "file path to be translate.")
	translatorCmd.Flags().StringVarP(&output, "output", "o", "", "output file path.")
	translatorCmd.Flags().StringVarP(&sourceLang, "source-lang", "s", "EN", "source language, like EN.")
	translatorCmd.Flags().StringVarP(&targetLang, "target-lang", "t", "ZH", "target language, like ZH.")
	translatorCmd.Flags().StringVarP(&proxyAPI, "proxy-api", "p", "", "proxy API url, like http://127.0.0.1/v2/translate.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// translatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
