package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/kovetskiy/godocs"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/colorgful"
	"github.com/reconquest/ser-go"
)

var (
	version = "[manual build]"
	usage   = "rutracker " + version + `


Usage:
  rutracker [options] -Q <query>
  rutracker -h | --help
  rutracker --version

Options:
  -Q --query <query>  Query torrent tracker.
  -h --help           Show this screen.
  --version           Show version.
  -c --config <path>  Use specified configuration file.
                       [default: $HOME/.config/rutracker/rutracker.conf]
  --debug             Print debug information.
`
)

var (
	logger    = lorg.NewLog()
	debugMode = false
)

type Tracker struct {
	client  *http.Client
	cookies *cookiejar.Jar
	baseURL string
}

func NewTracker(baseURL string) (*Tracker, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	tracker := &Tracker{
		client: &http.Client{
			Jar: jar,
		},
		cookies: jar,
		baseURL: baseURL,
	}

	return tracker, nil
}

func main() {
	args := godocs.MustParse(os.ExpandEnv(usage), version, godocs.UsePager)

	config, err := LoadConfig(args["--config"].(string))
	if err != nil {
		fatalh(err, "unable to load config: %s", args["--config"].(string))
	}

	logger.SetFormat(
		colorgful.MustApplyDefaultTheme(
			"${time} ${level:[%s]:right:short} ${prefix}%s",
			colorgful.Dark,
		),
	)

	debugMode = args["--debug"].(bool)
	if debugMode {
		logger.SetLevel(lorg.LevelDebug)
	}

	tracker, err := NewTracker(config.BaseURL)
	if err != nil {
		fatalh(
			err, "unable to initialize",
		)
	}

	client, err := authorize(
		config.BaseURL, config.Username, config.Password,
	)
	if err != nil {
		fatalh(
			err, "unable to authorize %s at %s",
			config.Username, config.BaseURL,
		)
	}

	fmt.Printf("XXXXXX main.go:64 client: %#v\n", client)

	switch {
	case args["--query"] != nil:
		err = handleQuery(args, config)
	}

	if err != nil {
		fatalln(err)
	}
}

func (tracker *Tracker) Authorize(
	baseURL, username, password string,
) error {
	payload := url.Values{}
	payload.Set("login_username", username)
	payload.Set("login_password", password)

	target := strings.TrimRight(baseURL, "/") + "/forum/login.php"

	request, err := http.NewRequest("POST", target, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return err
	}

	response, err := client.PostForm(target, payload)
	if err != nil {
		return nil, ser.Errorf(
			err, "unable to send POST request to %s", target,
		)
	}

	debugf("response status code: %s", response.Status)

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	debugf("response body: %s", body)

	if response.StatusCode == http.StatusMovedPermanently {
		return nil, nil
	}

	return nil, errors.New("invalid username or password")
}
