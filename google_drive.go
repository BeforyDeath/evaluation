package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var GoogleDrive googleDrive

type googleDrive struct {
	Service      *drive.Service
	modifiedTime string
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func (gd *googleDrive) getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	log.Info("Google Drive get client")

	cacheFile, err := gd.tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := gd.tokenFromFile(cacheFile)
	if err != nil {
		tok = gd.getTokenFromWeb(config)
		gd.saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func (gd *googleDrive) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	log.Info("Google Drive get Token From Web")

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Infof("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func (gd *googleDrive) tokenCacheFile() (string, error) {
	log.Info("Google Drive get token Cache File")

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("drive-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func (gd *googleDrive) tokenFromFile(file string) (*oauth2.Token, error) {
	log.Info("Google Drive get token From File")

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func (gd *googleDrive) saveToken(file string, token *oauth2.Token) {
	log.Infof("Saving credential file to: %s", file)

	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (gd *googleDrive) Init() error {
	log.Info("Google Drive initialization")

	ctx := context.Background()

	b, err := ioutil.ReadFile("config/client_secret.json")
	if err != nil {
		return err
	}

	// file path ~/.credentials/drive-go-quickstart.json
	config, err := google.ConfigFromJSON(b, drive.DriveScope, drive.DriveFileScope, drive.DriveReadonlyScope)
	if err != nil {
		return err
	}
	client := gd.getClient(ctx, config)

	srv, err := drive.New(client)
	if err != nil {
		return err
	}

	gd.Service = srv

	return nil
}
