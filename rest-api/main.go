package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

var appConfig = &Config{}

func main() {
	appConfig.LoadConfig("/etc/dyndns.json")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/update", Update).Methods("GET")
	router.HandleFunc("/resolve", Resolve).Methods("GET") // AFEGIM PER FER PSEUDO-RESOLUCIONS DNS

	/* DynDNS compatible handlers. Most routers will invoke /nic/update */
	router.HandleFunc("/nic/update", DynUpdate).Methods("GET")
	router.HandleFunc("/v2/update", DynUpdate).Methods("GET")
	router.HandleFunc("/v3/update", DynUpdate).Methods("GET")

	log.Println(fmt.Sprintf("Serving dyndns REST services on 0.0.0.0:8053..."))
	log.Fatal(http.ListenAndServe(":8053", router))
}

func DynUpdate(w http.ResponseWriter, r *http.Request) {
	extractor := RequestDataExtractor{
		Address: func(r *http.Request) string { return r.URL.Query().Get("myip") },
		Secret: func(r *http.Request) string {
			_, sharedSecret, ok := r.BasicAuth()
			if !ok || sharedSecret == "" {
				sharedSecret = r.URL.Query().Get("password")
			}

			return sharedSecret
		},
		Domain: func(r *http.Request) string { return r.URL.Query().Get("hostname") },
	}
	response := BuildWebserviceResponseFromRequest(r, appConfig, extractor)

	if response.Success == false {
		if response.Message == "Domain not set" {
			w.Write([]byte("notfqdn\n"))
		} else {
			w.Write([]byte("badauth\n"))
		}
		return
	}

	for _, domain := range response.Domains {
		result := UpdateRecord(domain, response.Address, response.AddrType)

		if result != "" {
			response.Success = false
			response.Message = result

			w.Write([]byte("dnserr\n"))
			return
		}
	}

	response.Success = true
	response.Message = fmt.Sprintf("Updated %s record for %s to IP address %s", response.AddrType, response.Domain, response.Address)

	w.Write([]byte(fmt.Sprintf("good %s\n", response.Address)))
}

func Update(w http.ResponseWriter, r *http.Request) {
	extractor := RequestDataExtractor{
		Address: func(r *http.Request) string { return r.URL.Query().Get("addr") },
		Secret:  func(r *http.Request) string { return r.URL.Query().Get("secret") },
		Domain:  func(r *http.Request) string { return r.URL.Query().Get("domain") },
	}
	response := BuildWebserviceResponseFromRequest(r, appConfig, extractor)

	if response.Success == false {
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, domain := range response.Domains {
		result := UpdateRecord(domain, response.Address, response.AddrType)

		if result != "" {
			response.Success = false
			response.Message = result

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response.Success = true
	response.Message = fmt.Sprintf("Updated %s record for %s to IP address %s", response.AddrType, response.Domain, response.Address)

	json.NewEncoder(w).Encode(response)
}

func UpdateRecord(domain string, ipaddr string, addrType string) string {
	log.Println(fmt.Sprintf("%s record update request: %s -> %s", addrType, domain, ipaddr))

	f, err := ioutil.TempFile(os.TempDir(), "dyndns")
	if err != nil {
		return err.Error()
	}

	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf("server %s\n", appConfig.Server))
	w.WriteString(fmt.Sprintf("zone %s\n", appConfig.Zone))
	w.WriteString(fmt.Sprintf("update delete %s.%s %s\n", domain, appConfig.Domain, addrType))
	w.WriteString(fmt.Sprintf("update add %s.%s %v %s %s\n", domain, appConfig.Domain, appConfig.RecordTTL, addrType, ipaddr))
	w.WriteString("send\n")

	w.Flush()
	f.Close()

	cmd := exec.Command(appConfig.NsupdateBinary, f.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err.Error() + ": " + stderr.String()
	}

	return out.String()
}

func Resolve(w http.ResponseWriter, r *http.Request) {
	extractor := RequestDataExtractor{
		Address: func(r *http.Request) string { return r.URL.Query().Get("addr") },
		Secret:  func(r *http.Request) string { return r.URL.Query().Get("secret") },
		Domain:  func(r *http.Request) string { return r.URL.Query().Get("domain") },
	}
	response := BuildWebserviceResponseFromRequest(r, appConfig, extractor)


	cmd := exec.Command("nslookup", fmt.Sprintf(response.Domain + "." + appConfig.Domain), "localhost")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	f, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	defer os.Remove(f.Name())
	wf := bufio.NewWriter(f)

	wf.WriteString(fmt.Sprintf("%s", out.String()))
	wf.Flush()
	f.Close()


	cmd2 := exec.Command("grep", "Address", f.Name())
	var out2 bytes.Buffer
	var stderr2 bytes.Buffer
	cmd2.Stdout = &out2
	cmd2.Stderr = &stderr2
	cmd2.Run()


	f2, err2 := ioutil.TempFile(os.TempDir(), "temp")
	if err2 != nil {
		response.Success = false
		response.Message = err2.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	defer os.Remove(f2.Name())
	wf2 := bufio.NewWriter(f2)

	wf2.WriteString(fmt.Sprintf("%s", out2.String()))
	wf2.Flush()
	f2.Close()


	cmd3 := exec.Command("tail", "-1", f2.Name())
	var out3 bytes.Buffer
	var stderr3 bytes.Buffer
	cmd3.Stdout = &out3
	cmd3.Stderr = &stderr3
	cmd3.Run()


	f3, err3 := ioutil.TempFile(os.TempDir(), "temp")
	if err3 != nil {
		response.Success = false
		response.Message = err3.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	defer os.Remove(f3.Name())
	wf3 := bufio.NewWriter(f3)

	wf3.WriteString(fmt.Sprintf("%s", out3.String()))
	wf3.Flush()
	f3.Close()


	cmd4 := exec.Command("cut", "-f", "2", "-d", " ", f3.Name())
	var out4 bytes.Buffer
	var stderr4 bytes.Buffer
	cmd4.Stdout = &out4
	cmd4.Stderr = &stderr4
	cmd4.Run()

	var result string
	
	fmt.Sscanf(out4.String(), "%s\n", &result)

	response.Success = true
	response.Message = result
	json.NewEncoder(w).Encode(response)
}