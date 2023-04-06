package aimagician_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/wztzt/aimagician"
)

func TestClient(t *testing.T) {
	var Cookies []http.Cookie
	json.Unmarshal([]byte(
		`[
		{
		  "name": "__Secure-session-token",
		  "value": "",
		  "path": "/",
		  "domain": "ai-prompt.guokr.net",
		  "expires": "2023-04-22T05:32:49.158Z",
		  "httpOnly": true,
		  "secure": true,
		  "sameSite": "Lax"
		},
		{
		  "name": "__Secure-uid",
		  "value": "",
		  "path": "/",
		  "domain": "ai-prompt.guokr.net",
		  "expires": "2023-04-22T05:32:49.158Z",
		  "httpOnly": false,
		  "secure": true,
		  "sameSite": "Lax"
		}
	  ]
	`), &Cookies)
	asks := []string{"讲个故事", "讲个笑话", "五粮液", "茅台"}
	client := aimagician.NewClient(Cookies)
	stream := client.ChatStream(asks[0])
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(*response)
	}
	/*for i := 0; i < 100; i++ {
		reply := client.Chat(asks[i%4])
		fmt.Println(reply)
	}*/

}
