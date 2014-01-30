package main

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/keyring"
	"github.com/skycoin/skycoin/src/util"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Profile struct {
	Name    string
	Hobbies []string
}

type test_struct struct {
	WalletName string
}

type walletData struct {
	Seed      string
	Addresses []string
	History   []string
}

var (
	logger = logging.MustGetLogger("skycoin.gui")
)

var WalletFile = walletData{}

func main() {
	static_path, _ := filepath.Abs("../../static/app/")
	logger.Debug("Serving %s", static_path)

	//readWriteFile()

	
	http.Handle("/", http.FileServer(http.Dir(static_path)))

	http.HandleFunc("/api/loadWallet", loadWallet)

	http.HandleFunc("/api/saveWallet", saveWallet)

	http.HandleFunc("/api/newAddress", newAddress)

	fmt.Printf("Server Running on: 127.0.0.1:3003 \n")
	err := http.ListenAndServe(":3003", nil)
	if err != nil {
		log.Panic(err)
	}

}

func loadWallet(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t test_struct

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Panic(err)
	}

	value := t.WalletName

	logger.Debug("Serving %s", value)

	fmt.Printf("walletName = %s", value+".wallet")

	//LoadedJSON := util.LoadJSON(value+".wallet", WalletFile)

	err = util.LoadJSON(value+".wallet", &WalletFile)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("json= %v type= %T \n", WalletFile, WalletFile)

	logger.Debug("LoadedJSON = %s", WalletFile)

	js, err := json.Marshal(WalletFile)

	_ = err

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)

}

func saveWallet(w http.ResponseWriter, r *http.Request) {
	//data, err := json.Marshal(thing)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t walletData

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Panic(err)
	}

	year, month, day := time.Now().Date()

	DateStr := fmt.Sprintf("%04d_%02d_%02d", year, month, day)

	util.SaveJSON(DateStr+".wallet", t)

	//fmt.Printf("address= %s \n", DateStr)

	js, err := json.Marshal(DateStr)
	_ = err
	//js, err := json.Marshal(addr.Address.String())

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}

func newAddress(w http.ResponseWriter, r *http.Request) {

	//js, err := json.Marshal(profile)
	addr := keyring.NewAddress()

	//walletFile.Addresses = append(walletFile.Addresses, addr)
	fmt.Printf("address= %s \n", addr.Address.String())

	js, err := json.Marshal(addr.Address.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LoadJSON(filename string, thing interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, thing)
}

func SaveJSON(filename string, thing interface{}) error {
	data, err := json.Marshal(thing)
	if err != nil {
		return err
	}
	tmpname := filename + ".tmp"
	err = ioutil.WriteFile(tmpname, data, os.FileMode(0644))
	if err != nil {
		return err
	}
	return os.Rename(tmpname, filename)
}