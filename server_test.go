package restful

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go func() {
		server := NewServer(":801")
		server.AddRestfulHandler("/servers", http.MethodPost, "/{id}/info", func(w http.ResponseWriter, r *http.Request, uriParams map[string]string) {
			w.Write([]byte(fmt.Sprintf("server id is %s", uriParams["id"])))
		})
		server.AddRestfulHandler("/servers", http.MethodGet, "/{id}/info/{field}", func(w http.ResponseWriter, r *http.Request, uriParams map[string]string) {
			w.Write([]byte(fmt.Sprintf("get %s for %s", uriParams["field"], uriParams["id"])))
		})
		server.Start()
	}()

	time.Sleep(time.Second)
	resp, err := http.Post("http://localhost:801/servers/1/info", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("request server fail, err=%s", err)
		t.Fail()
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	t.Logf("response body is: %s", respBody)
	if string(respBody) != "server id is 1" {
		t.Fail()
	}

	resp, err = http.Get("http://localhost:801/servers/1/info/f1")
	if err != nil {
		t.Errorf("request server fail, err=%s", err)
		t.Fail()
	}
	respBody, _ = ioutil.ReadAll(resp.Body)
	t.Logf("response body is: %s", respBody)
	if string(respBody) != "get f1 for 1" {
		t.Fail()
	}

	resp, err = http.Get("http://localhost:801/servers/1/info")
	if err != nil {
		t.Errorf("request server fail, err=%s", err)
		t.Fail()
	}
	if resp.StatusCode != 404 {
		t.Errorf("status code %d is not equal to 404", resp.StatusCode)
		t.Fail()
	}
}
