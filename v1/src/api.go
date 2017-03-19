/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 11/22/2014

 Module: api
 Description: RESTful API that parts of the web ui uses.

*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type API interface {
	POST(string, []byte) error
	GET(string) ([]byte, error)
}

// V1 represents a RESTful JSON API.
type V1 struct {
	controller *Controller
}

// Unsafe represents what I am calling 'unsafe' API calls. Which are basically
// straight up console commands sent to the controller.
type Unsafe struct {
	controller *Controller
}

func APIHandler(w http.ResponseWriter, req *http.Request, controller *Controller) {
	// TODO: set up the correct error codes.
	r := strings.SplitN(strings.Trim(req.URL.Path, "/"), "/", 3)
	var err error
	if len(r) >= 3 {
		var api API
		var payload []byte
		switch r[1] {
		case "v1":
			api = V1{controller: controller}
		case "unsafe":
			api = Unsafe{controller: controller}
		default:
			err = errors.New("Bad Request")
			return
		}
		switch req.Method {
		case "POST":
			if payload, err = ioutil.ReadAll(req.Body); err == nil {
				err = api.POST(r[2], payload)
			}

		case "GET":
			if payload, err = api.GET(r[2]); err == nil {
				w.Write(payload)
			}
		default:
			err = errors.New("Bad Request")
		}
	} else {
		err = errors.New("Bad Request")
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 400)
	}
}

func (v V1) POST(resource string, payload []byte) error {
	fmt.Println(string(payload))
	switch resource {
	case "tourney":
		filename := filepath.Join(Settings.TourneyDirectory, strconv.FormatInt(time.Now().UTC().UnixNano(), 10)+".tourney")
		fmt.Println("Saving", filename)
		if err := ioutil.WriteFile(filename, payload, 0644); err == nil {
			v.controller.Enque("stop")
			v.controller.Enque("load " + filename)
			v.controller.Enque("host")
		} else {
			return errors.New("500 " + err.Error())
		}

	default:
		return errors.New("404")
	}
	return nil
}

func (v V1) GET(resource string) ([]byte, error) {
	return []byte{}, nil
}

func (v Unsafe) POST(resource string, payload []byte) error {

	return nil
}

func (v Unsafe) GET(resource string) ([]byte, error) {
	r := strings.Split(strings.Trim(resource, "/"), "/")
	cmd, arg := "", ""
	if len(r) >= 2 {
		arg = " " + r[1]
	}
	if len(r) >= 1 {
		cmd = r[0]
		fmt.Println("[API-RAW] Recieved: " + cmd + arg)
		v.controller.Enque(cmd + arg)
	}
	return []byte{}, nil
}
