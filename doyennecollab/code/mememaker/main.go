package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/nathj07/talks/doyennecollab/code/mememaker/data"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-multierror"
)

var (
	username = flag.String("username", "", "Your username for the api.imgflip.com service")
	password = flag.String("password", "", "Your password for the api.imgflip.com service")
	action   = flag.String("action", "GET", "The action to perform against api.imgflip.com")
	memeID   = flag.Int("meme", 0, "Meme template ID to ue in creating a new meme. Not needed for GET requests")
	text0 = flag.String("text0", "", "top Line Of TExt fOR yUoR memE")
	text1 = flag.String("text1", "", "bototom Line Of TExt fOR yUoR memE")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
      %s is a tool to call the api.imgflip.com service for the purpose of finding and creating memes
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if strings.EqualFold(*action, http.MethodPost) {
		// check the input is valid
		var err *multierror.Error
		if *username == "" {
			err = multierror.Append(err, fmt.Errorf("You must supply a username"))
		}
		if *password == "" {
			err = multierror.Append(err, fmt.Errorf("You must supply a password"))
		}
		if *memeID == 0 {
			err = multierror.Append(err, fmt.Errorf("You must supply a meme template ID"))
		}
		if e := err.ErrorOrNil(); e != nil {
			log.Fatal(err)
		}
		// TODO: Build up post request - this will need more cli args too - see the request struct
		// plus an output path to write the file to
	}

	if strings.EqualFold(*action, http.MethodGet) {
		makeGetRequest()
		return
	}
	makePostRequest()
}

// makeGet request is here as a simple example of how to make a
// basic HTTP GET request with Go.
// For production purposes we would need to define our own HTTP client and not rely on the default
func makeGetRequest() {
	resp, err := http.Get("http://api.imgflip.com/get_memes")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code returned from GET request: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	memes := &data.MemeGetResponse{}
	if err := json.Unmarshal(body, memes); err != nil {
		log.Fatal(err)
	}
	spew.Dump(memes)
}

// makePostRequest POSTS the details to the API
// This needs to be completed.
// You will need to:
// - add more args to the cli, and validate them
// - make POST request
// - unmarshal the response, the data structure will depend on the status code
// - fetch the created meme and write it to disk
// (in the future we may work on displaying the meme, you use
// os.Exec with the open command if you feel like it)
func makePostRequest() {
	values := url.Values{}

	// ignoring errors here for now; need to address that
	values.Set("template_id", strconv.Itoa(*memeID))
	values.Set("username", *username)
	values.Set("password", *password)
	values.Set("text0", *text0)
	values.Set("text1", *text1)

	resp, err := http.PostForm("https://api.imgflip.com/caption_image", values)
	if err != nil {
		log.Fatalf("unable to post form %v", err)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read body %v", err)
	}

	memeResp := &data.MemePostResponse{}
	err = json.Unmarshal(d, memeResp)
	if err != nil {
		log.Fatalf("unable to read json %v", err)

	}
	spew.Dump(memeResp)
}

// Once the makeRequest is done feel free to improve upon this code, look at using an http client,
// other than the default one used in the Get example.
// Whatever features you want to add, or improvements you want to make please go ahead.
