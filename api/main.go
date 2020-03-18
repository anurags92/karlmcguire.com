package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	DB = "data.gob"
)

type Views struct {
	sync.Mutex
	access chan struct{}
	count  map[string]uint64
}

func NewViews() *Views {
	v := &Views{access: make(chan struct{}, 100)}
	if err := v.Load(); err != nil {
		fmt.Println("error loading: " + DB)
		v.count = make(map[string]uint64)
	}
	go v.Saver()
	return v
}

func (v *Views) Saver() {
	for i := uint64(0); ; i++ {
		<-v.access
		if i == 100 {
			v.Save()
			i = 0
		}
	}
}

func (v *Views) Add(path string) uint64 {
	v.access <- struct{}{}
	v.Lock()
	defer v.Unlock()
	n := v.count[path]
	v.count[path]++
	return n
}

func (v *Views) Save() error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	v.Lock()
	err := enc.Encode(v.count)
	if err != nil {
		return err
	}
	v.Unlock()
	return ioutil.WriteFile(DB, buf.Bytes(), 0644)
}

func (v *Views) Load() error {
	data, err := ioutil.ReadFile(DB)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&v.count)
}

func main() {
	views := NewViews()
	http.HandleFunc("/views", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "GET")
		//w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		r.ParseForm()
		fmt.Fprintf(w, "%d", views.Add(r.FormValue("path")))
	})
	fmt.Println("listening on :80")
	panic(http.ListenAndServe(":80", nil))
}
