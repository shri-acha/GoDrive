package main

import (
	"net/http"
	"GoDrive/handlers"
)


func main(){
	
	http.HandleFunc("/",handlers.HandleHome)	
	http.HandleFunc("/upload",handlers.HandleUpload)	
	http.HandleFunc("/delete",handlers.HandleDelete)	

	http.ListenAndServe(":8080",nil) // In port 8080
}

