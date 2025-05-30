package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	machine1 := "/452FG92"
	machine2 := "/1ASD987"
	machine3 := "/1TREW89"
	machine4 := "/452FG92/get_mac_addr"

	mux := http.NewServeMux()
	mux.HandleFunc(machine1, MachineSuccessHandler)
	mux.HandleFunc(machine2, MachineFailedHandler)
	mux.HandleFunc(machine3, MachineTimeoutRequestHandler)
	mux.HandleFunc(machine4, GetMachineMacAddrHandler)

	http.ListenAndServe(":8000", mux)
}

func MachineSuccessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle machine 1")
	w.WriteHeader(http.StatusOK)
}

func MachineFailedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle machine 2")
	w.WriteHeader(http.StatusNotFound)
}

func MachineTimeoutRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle machine 3")
	time.Sleep(time.Second * 10)
	w.WriteHeader(http.StatusNotFound)
}

func GetMachineMacAddrHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle machine 1")
	payload := []byte(`{"router_bssid": "12.12.12.12"}`)
	w.Write(payload)
	w.WriteHeader(http.StatusOK)
}
