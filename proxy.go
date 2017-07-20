package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ProxyServer the web server
func ProxyServer(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "Mensaje") // to write message to client
	//fmt.Println(req.RequestURI) // to print request URI
	u, err := url.Parse(req.RequestURI) // parse to URL object form URI
	if err != nil {
		log.Println("Error Parse:", err)
	}

	q := u.Query() // query params of URL

	// If url gets /?form
	if len(q["form"]) > 0 { // if request by param form
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `
      <form action="/" method="get">
          URL: <input type="text" name="url">
          <input type="submit" value="Send">
      </form>
    `)
		return
	}

	// If url NOT gets /?url=your_url
	if len(q["url"]) == 0 {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "Please use http://localhost:12345?url=your_url encoded with encodeURI.\n")
		io.WriteString(w, "<script>setTimeout(function(){window.location='/?form'}, 5000);</script>")
		return
	}

	url := q["url"][0] // file url for proxy
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
	//serverValue := req.Header().Get("Server")

	// Search it is in whitelist (YOU CAN REMOVE)
	isInWhitelist := isInList(url)
	if isInWhitelist == false {
		log.Println("Error isNotInWhitelist:")
		io.WriteString(w, "The URL is not in whitelist.lst.\n")
		return
	}

	// assembly request of URL
	client := http.Client{} // request client
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error GET:", err)
	}

	// execute petition
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error Client DO:", err)
	}
	defer response.Body.Close()

	// get response data
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error responseData:", err)
	}

	// pass headers from request to response
	for k, v := range response.Header {
		w.Header().Set(k, v[0])
	}

	// w.Header().Set("AtEnd1", "value 1")
	// w.Header().Add("AtEnd1", "AtEnd3") // Not rewrite header
	// w.Header().Set("Saludo", "Hola Andres")

	// write data and exit
	w.Write(responseData)
}

func main() {
	http.HandleFunc("/", ProxyServer)
	fmt.Println("Listening on http://localhost:12345")
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func isInList(url string) bool {
	isIn := false
	file, err := os.Open("whitelist.lst")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line, url)
		i := strings.Index(url, line)
		if i == 0 { // Is in first position/index
			isIn = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return isIn
}
