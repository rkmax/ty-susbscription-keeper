package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	launchWebServer = true
	redirectUrl     = "http://localhost:8090"
	tcpHost         = "localhost:8090"
)

func getClient(scope string) *http.Client {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_web_secrets.json")
	if err != nil {
		log.Fatalf("Unable to read client secrets %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secrets %v", err)
	}

	if launchWebServer {
		config.RedirectURL = redirectUrl
	} else {
		config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	}

	cacheFile, err := tokenCacheFile()

	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	token, err := readTokenFromFile(cacheFile)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

		if launchWebServer {
			fmt.Println("trying to get token from web")
			token, err = getTokenFromWeb(config, authURL)
		} else {
			fmt.Println("trying to get token from prompt")
			token, err = getTokenFromPrompt(config, authURL)
		}

		if err == nil {
			saveTokenToFile(cacheFile, token)
		}
	}

	return config.Client(ctx, token)
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")

	err = os.MkdirAll(tokenCacheDir, 0700)
	if err != nil {
		return "", err
	}

	return filepath.Join(tokenCacheDir, url.QueryEscape("ty-susbscription-keeper.json")), err
}

func readTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()

	return t, err
}

func saveTokenToFile(file string, token *oauth2.Token) {
	fmt.Println("trying to save token")
	fmt.Printf("Saving credential fila to: %s\n", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache credential %v", err)
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("Unable to cache credential %v", err)
	}

}

func getTokenFromWeb(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	codeCh, err := startWebServer()
	if err != nil {
		fmt.Printf("Unable to start web server")
	}

	err = openURL(authURL)
	if err != nil {
		log.Fatalf("Unable to open authorization URL in web server: %v", err)
	} else {
		fmt.Println("Your web browser has been opened to an authorization URL", "This program will resume once authorization has been provided.", authURL)
	}

	code := <-codeCh

	return exchangeToken(config, code)
}

func getTokenFromPrompt(config *oauth2.Config, authURL string) (*oauth2.Token, error) {
	var code string

	fmt.Printf("Go to the following link in your browser. After completing"+"the authentization flow, enter the authorization code on the command line: \n%v\n", authURL)

	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	return exchangeToken(config, code)
}

func startWebServer() (codeCh chan string, err error) {
	listener, err := net.Listen("tcp", tcpHost)
	if err != nil {
		return nil, err
	}

	codeCh = make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		code := request.FormValue("code")
		codeCh <- code
		defer listener.Close()
		writer.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(writer, "Code received: \r\nYou can now safely close this browser window.")
	}))

	return codeCh, nil
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("cannot open URL %s on this platform", url)
	}

	return err
}

func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}

	return token, nil
}
