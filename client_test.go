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
		  "value": "token",
		  "path": "/",
		  "domain": "ai-prompt.guokr.net",
		  "expires": "",
		  "httpOnly": true,
		  "secure": true,
		  "sameSite": "Lax"
		},
		{
		  "name": "__Secure-uid",
		  "value": "uid",
		  "path": "/",
		  "domain": "ai-prompt.guokr.net",
		  "expires": "",
		  "httpOnly": false,
		  "secure": true,
		  "sameSite": "Lax"
		}
	  ]
	`), &Cookies)
	client := aimagician.NewClient(Cookies)
	for i := 0; i < 100; i++ {
		reply := client.Chat("基金")
		fmt.Println(reply)
	}

}
