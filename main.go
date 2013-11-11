package main

import (
	"bytes"
	"encoding/json"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-UserService/models"
	"github.com/wurkhappy/mdp"
	"net/http"
	"net/url"
	"strconv"
)

type ServiceReq struct {
	Method string
	Path   string
	Body   []byte
}

func main() {
	config.Prod()
	models.Setup()
	router.Start()

	gophers := 10

	for i := 0; i < gophers; i++ {
		worker := mdp.NewWorker("tcp://localhost:5555", config.UserService, false)
		defer worker.Close()
		go route(worker)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//route to function based on the path and method
		route, pathParams, _ := router.FindRoute(r.URL.String())
		routeMap := route.Dest.(map[string]interface{})
		handler := routeMap[r.Method].(func(map[string]interface{}, []byte) ([]byte, error, int))

		//parse the request
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		//add url params to params var
		params := make(map[string]interface{})
		for key, value := range pathParams {
			params[key] = value
		}
		//add url query params
		values := r.URL.Query()
		for key, value := range values {
			params[key] = value
		}

		//run handler and do standard http stuff(write JSON, return err, set status code)
		jsonData, err, statusCode := handler(params, buf.Bytes())
		if err != nil {
			http.Error(w, `{"status_code":`+strconv.Itoa(statusCode)+`, "description":"`+err.Error()+`"}`, statusCode)
			return
		}
		w.WriteHeader(statusCode)
		w.Write(jsonData)
	})
	http.ListenAndServe(":3000", nil)
}

type Resp struct {
	Body       []byte `json:"body"`
	StatusCode int    `json:"status_code"`
}

func route(worker mdp.Worker) {
	for reply := [][]byte{}; ; {
		request := worker.Recv(reply)
		if len(request) == 0 {
			break
		}
		var req *ServiceReq
		json.Unmarshal(request[0], &req)

		//route to function based on the path and method
		route, pathParams, _ := router.FindRoute(req.Path)
		routeMap := route.Dest.(map[string]interface{})
		handler := routeMap[req.Method].(func(map[string]interface{}, []byte) ([]byte, error, int))

		//add url params to params var
		params := make(map[string]interface{})
		for key, value := range pathParams {
			params[key] = value
		}
		//add url query params
		uri, _ := url.Parse(req.Path)
		values := uri.Query()
		for key, value := range values {
			params[key] = value
		}

		//run handler and do standard http stuff(write JSON, return err, set status code)
		jsonData, err, statusCode := handler(params, req.Body)
		if err != nil {
			resp := &Resp{[]byte(`{"description":"` + err.Error() + `"}`), statusCode}
			d, _ := json.Marshal(resp)
			reply = [][]byte{d}
			continue
		}
		resp := &Resp{jsonData, statusCode}
		d, _ := json.Marshal(resp)
		reply = [][]byte{d}
	}
}
