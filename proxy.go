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

	conf "./conf"
)

var (
	// Address is a complete service URL http://hostname:port/context
	Address string
)

func main() {
	p := conf.Parameters
	Address = "http://" + p.HOSTNAME + ":" + p.PORT + p.CONTEXT
	http.HandleFunc(p.CONTEXT, ProxyServer)
	fmt.Println("Listening on " + Address)
	log.Fatal(http.ListenAndServe(p.HOSTNAME+":"+p.PORT, nil))
}

// ProxyServer the web server
func ProxyServer(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	//io.WriteString(w, "Mensaje") // to write message to client
	//fmt.Println(req.RequestURI) // to print request URI
	u, err := url.Parse(req.RequestURI) // parse to URL object form URI
	if err != nil {
		log.Println("Error Parse: ", err)
	}

	q := u.Query() // query params of URL

	// If url gets /?form
	if len(q["form"]) > 0 { // if request by param form
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `
      <form action="" method="get">
          URL: <input type="text" name="url">
          <input type="submit" value="Send">
      </form>
    `)
		return
	}

	// If url NOT gets /?url=your_url
	if len(q["url"]) == 0 {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "Please use "+Address+"?url=your_url encoded with encodeURI.\n")
		io.WriteString(w, "<script>setTimeout(function(){window.location='"+Address+"?form'}, 5000);</script>")
		return
	}

	url := q["url"][0] // file url for proxy
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
	//serverValue := req.Header().Get("Server")

	// Search it is in whitelist (YOU CAN REMOVE)
	isInWhitelist := isInList(url)
	if isInWhitelist == false {
		log.Println("Error isNotInWhitelist: " + url)
		io.WriteString(w, "The URL is not in whitelist.lst.\n")
		return
	}

	// assembly request of URL
	client := http.Client{} // request client
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error GET: ", err)
	}

	if len(q["headers"]) > 0 {
		headers := q["headers"][0]
		// would be oriented to headers separed by \r\n literally strings for LF CR
		// example headers=Cookie:####\r\nOrigin:OtherOriginValue
		// im spect this is sufficient
		headersList := strings.Split(headers, "\\r\\n")
		for index := 0; index < len(headersList); index++ {
			header := strings.Split(headersList[index], ":")
			if len(header) < 2 {
				log.Println("Error HEADER: ", header)
			} else {
				request.Header.Set(header[0], header[1])
			}
		}

	}

	// execute petition
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error Client DO: ", err)
	}
	defer response.Body.Close()

	// get response data
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error responseData: ", err)
	}

	renameheaders := false
	if len(q["renameheaders"]) > 0 { // si viene el parametro con cualquier valor
		renameheaders = true
	}

	// pass headers from request to response
	headerList := ""
	for k, v := range response.Header {
		w.Header().Set(k, v[0])
		if renameheaders {
			w.Header().Set("_"+k, v[0])
			headerList += "_" + k + ","
		}
		headerList += k + ","
	}
	headerList = headerList[0 : len(headerList)-1]

	// w.Header().Set("AtEnd1", "value 1")
	// w.Header().Add("AtEnd1", "AtEnd3") // Not rewrite header
	// w.Header().Set("Saludo", "Hola Andres")

	// add CORS headers
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS#Access-Control-Expose-Headers
	w.Header().Set("Access-Control-Allow-Origin", origin)
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", headerList)
	w.Header().Set("Access-Control-Expose-Headers", headerList)
	// w.Header().Set("Access-Control-Allow-Headers", "Set-Cookie")
	// w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
	// w.Header().Set("Access-Control-Allow-Methods", "GET")
	//w.Header().Set("Access-Control-Expose-Headers", "ETag")

	// write data and exit
	w.Write(responseData)
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
		i := strings.Index(url, line)
		if i == 0 { // Is in first position/index
			isIn = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return isIn
}
