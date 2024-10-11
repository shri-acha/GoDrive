package handlers

import (
	"net/http"
	"os"
	"fmt"
)

func HandleDelete(w http.ResponseWriter , r* http.Request){
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid Method[CS]", 405)
	}
	
	fileLocation := r.Header.Get("X-File-Location")

	err := os.Remove("./uploads/"+fileLocation)
	
	if err != nil {
		fmt.Fprintf(w, "Error in deleting the file!")
	}

	fmt.Fprintf(w,"File deleted successfully!");
}
