package main 

import(
	//bulit-in
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request){

		//parse querydate
		r.ParseForm()

		//sQuestion := r.Form["q"][0]


		w.Write([]byte("我很帅！"))
	})

	oHttpServer := &http.Server{
		Addr : ":8521",
		Handler : nil,
		ReadTimeout: time.Duration(20) * time.Second,
		WriteTimeout: time.Duration(20) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//start listen and server
	oHttpServer.ListenAndServe()

}