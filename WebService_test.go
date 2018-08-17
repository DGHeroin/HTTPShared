package HTTPShared

import (
	"testing"
	"time"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
)

func TestWebService(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	web := NewWebService(":9999", "")
	go func() {
		web.WaitExit()
	}()
	time.Sleep(time.Second)

	testGet()

	testSet()

	testGet()

	testWaitOneTime()

	testSet()

	testWaitRepeat()
	testSet()
	testSet()

	testWaitRepeatHeader()
	testSet()
	testSet()

	time.Sleep(time.Second * 3)
	web.Stop()
}

func testGet()  {
	url := "http://127.0.0.1:9999/v1/keys/message"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Println(string(body))
	}
}

func testSet()  {
	url := "http://127.0.0.1:9999/v1/keys/message"
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, strings.NewReader("GolangIsNice"))
	if err != nil {
		log.Println(err)
		return
	}
	defer req.Body.Close()
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Println(string(body))
	}
}

func testWaitOneTime() {
	go func() {
		url := "http://127.0.0.1:9999/v1/keys/message?wait=true"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Println("wait one time", string(data))
			return
		}
	}()
	time.Sleep(time.Second)
}

func testWaitRepeat() {
	go func() {
		url := "http://127.0.0.1:9999/v1/keys/message?wait=true&r=true"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}

		defer resp.Body.Close()
		for {
			buf := make([]byte, 4096)
			n, err := resp.Body.Read(buf)
			if err != nil {
				log.Println(err)
				break
			}
			data := buf[:n]
			log.Println("wait with header:", n, string(data))
		}

	}()
	time.Sleep(time.Second)
}

func testWaitRepeatHeader() {
	go func() {
		url := "http://127.0.0.1:9999/v1/keys/message?wait=true&r=true&h=true"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		readHead := true
		var body     []byte
		var bodySize int
		var content  string
		for {
			buf := make([]byte, 4096)
			n, err := resp.Body.Read(buf)
			if err != nil {
				log.Println(err)
				break
			}
			body = append(body, buf[:n]...)
			if readHead {
				if len(body) < 4 {
					continue
				}
				bodySize =
					int(body[0] << 24) |
					int(body[1] << 16) |
					int(body[2] << 8) |
					int(body[3] << 0)
				body = body[4:]
				bodySize = bodySize + 1 // 刷新数据需要写入  '\n'
				readHead= false
			} else {
				if len(body) < bodySize {
					continue
				}
				content = string(body[:bodySize])
				body = body[:bodySize]
				log.Println("data:", n, string(content))
			}


		}
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Println("wait", string(body))
		}
	}()
	time.Sleep(time.Second)
}