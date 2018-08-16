package HTTPShared

import (
	"net/http"
	"log"
	"context"
	"encoding/json"
	"path"
	"io/ioutil"
	"fmt"
)

type WebService struct {
	shared *Shared

	Server http.Server
	stop chan bool

	Token *string

	Route map[string] func(http.ResponseWriter,*http.Request)
}

func NewWebService(addr string, Token string) *WebService {
	var web WebService
	web.shared = NewShared()
	web.Server.Addr = addr
	web.Route = make(map[string] func(http.ResponseWriter,*http.Request))
	web.stop  = make(chan bool)
	web.Token = nil
	if Token != "" {
		web.Token = &Token
	}

	go web.start()
	return &web
}

func (this *WebService) WaitExit()  {
	select {
	case <- this.stop:
	}
}

func (this *WebService) Stop() {
	go func() {
		this.stop <- true
		if err := this.Server.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()
}

func (this *WebService) start() {
	http.HandleFunc("/", this.handleHTTP)
	this.RegisterRoute()
	if err := this.Server.ListenAndServe(); err != nil {
		log.Println("Http Server ListenAndServe Error:", err)
		this.Stop()
		return
	}
}


func ResponseError(w http.ResponseWriter, code int) {
	if data, err := json.Marshal(ErrorData{Code:code, Msg:http.StatusText(code)}); err == nil {
		w.WriteHeader(code)
		w.Write([]byte(data))
	}
}

func ResponseInternalError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, 500)
}

func ResponseJson(w http.ResponseWriter, i interface{}) (error) {
	data, err := json.Marshal(i);
	if  err == nil {
		w.Write(data)
		return nil
	}
	return err
}

func (this *WebService) handleHTTP(w http.ResponseWriter, r *http.Request) {
	route := path.Dir(r.URL.String())
	if cb, ok := this.Route[route]; ok {
		cb(w, r)
		return
	} else {
		ResponseError(w, http.StatusBadRequest)
	}
}

func (this *WebService) RegisterRoute() {
	this.Route["/v1/keys"] = this.v1Keys
}

func (this *WebService) v1Keys(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body []byte
	var err error
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		ResponseError(w, http.StatusBadRequest)
		return
	}

	if this.Token != nil {
		args := r.URL.Query()
		if args != nil && args.Get("token") !=  *this.Token {
			ResponseError(w, http.StatusNetworkAuthenticationRequired)
			return
		}
	}

	if r.Method == "PUT" {
		key := path.Base(r.URL.Path)
		value := string(body)
		version := this.shared.Set(key, value)
		if err := ResponseJson(w, &PutActionData{Version:version}); err != nil {
			ResponseInternalError(w, err.Error())
		}
		return
	} else {
		// get value
		key := path.Base(r.URL.Path)
		args := r.URL.Query()
		isWaiting := false
		isRepeat := false
		if args != nil && args.Get("wait") == "true" {
			isWaiting = true
			if args.Get("r") == "true" {
				isRepeat = true
			}
		}
		if isWaiting == false {
			result := this.shared.Get(key)
			if result != nil {
				v := GetActionData{Key: result.Key, Value: result.Value, Version: result.Version}
				ResponseJson(w, &v)
			} else {
				w.Write([]byte("null"))
			}
		} else {
			if isRepeat {
				this.waitRepeat(w, key)
			} else {
				this.waitOneTime(w, key)
			}
		}
	}
}

func (this* WebService) waitRepeat(w http.ResponseWriter, key string) {
	ch := make(chan bool)
	done := make(chan bool)

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	Id, _ := this.shared.Watch2(key, func(result *Result) bool {
		v := GetActionData{Key: result.Key, Value: result.Value, Version: result.Version}
		ResponseJson(w, &v)
		if flusher, ok := w.(http.Flusher); ok {
			fmt.Fprint(w, "\n") // 需要写入换行符才会真正刷新
			flusher.Flush()
		}

		return true
	})
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		case <-notify:
			this.shared.UnWatch(key, Id)
		case <-done:
		}
	}()

	<-ch
}


func (this* WebService) waitOneTime(w http.ResponseWriter, key string) {
	ch := make(chan bool)
	done := make(chan bool)

	Id, _ := this.shared.Watch2(key, func(result *Result) bool {
		v := GetActionData{Key: result.Key, Value: result.Value, Version: result.Version}
		ResponseJson(w, &v)
		ch <- true
		return false
	})
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		case <-notify:
			this.shared.UnWatch(key, Id)
		case <-done:
		}
	}()

	<-ch
}