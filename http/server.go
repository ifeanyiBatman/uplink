package main

import (
	"net/http"
)

func main(){
	serveMux := http.NewServeMux()
	srv := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}
	serveMux.HandleFunc("/home", handlerHome )
	
	srv.ListenAndServe()
}

func handlerHome (w http.ResponseWriter, req *http.Request){
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Hello Fucking World"))
}
