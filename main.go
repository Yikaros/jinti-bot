// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    	"bufio"
	"fmt"
	"log"
    	"io"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	fi, err := os.Open("buffer/list.txt")
    	if err != nil {
        	fmt.Printf("Error: %s\n", err)
        	return
    	}
    	defer fi.Close()

    	br := bufio.NewReader(fi)
	var list string
	var work string
	var stock string
	var msg string
    	for {
        	a, _, c := br.ReadLine()
        	if c == io.EOF {
            	break
        	}
		list = list + "&" + string(a)
    	}
	
	var list_array []string
	list_array = strings.Split(list, "&")
	
	fi2, err2 := os.Open("buffer/ppl.txt")
    	if err2 != nil {
        	fmt.Printf("Error: %s\n", err2)
        	return
    	}
    	defer fi2.Close()
	list = ""
    	br2 := bufio.NewReader(fi2)
    	for {
        	a, _, c := br2.ReadLine()
        	if c == io.EOF {
            	break
        	}
		list = list + "&" + string(a)
    	}
	
	var list_array2 []string
	list_array2 = strings.Split(list, "&")
	
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch {
					case Index(message.Text,"@") == 1:
						i:=0
						itemname := message.Text
						for i<=len(list_array2)-1{
							var menu []string
							menu = strings.Split(list_array2[i], "$")
							if menu[0] == itemname{
								msg = menu[1]
								break
							}
							i++
						}
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do()
//查庫存的code
					default:
						i:=0
						itemname := message.Text
						for i<=len(list_array){
							var menu []string
							menu = strings.Split(list_array[i], "@")
							if menu[0] == itemname{
								stock=menu[1]
								work=menu[2]
								break
							}
							i++
						}
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(itemname + "還有庫存" + stock + "支，在製品" + work + "支i")).Do() 
				}
					
			}

		}
	}
}

func Contains(s, substr string) bool {
     return Index(s, substr) != -1
}

func Index(s string, sep string) int {
    n := len(sep)
    if n == 0 {
        return 0
    }
    c := sep[0]
    if n == 1 {
        // special case worth making fast
        for i := 0; i < len(s); i++ {
            if s[i] == c {
                return i
            }
        }
        return -1
    }
    // n > 1
    for i := 0; i+n <= len(s); i++ {
        if s[i] == c && s[i:i+n] == sep {
            return i
        }
    }
    return -1
}
