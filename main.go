package main

import (
	"bytes"
	"net/http"
	"strconv"
)

func main() {
	router.Start()

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
			http.Error(w, `{"status_code":"`+strconv.Itoa(statusCode)+`", "description":"`+err.Error()+`"}`, statusCode)
			return
		}
		w.WriteHeader(statusCode)
		w.Write(jsonData)
	})
	http.ListenAndServe(":3000", nil)
}
