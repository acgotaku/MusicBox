package xiami

func SearchHandler(w http.ResponseWriter, r *http.Request) {

}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write("Hello World")

}
