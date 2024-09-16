package api_client_go

import (
	"bufio"
	"encoding/json"
	"github.com/nuvla/api-client-go/common"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

type NuvlaCookies struct {
	jar        http.CookieJar
	lastCookie []*http.Cookie
	endpoint   *url.URL
	cookieFile string
}

// NewNuvlaCookies creates a new instance of the NuvlaCookies struct.
// It takes two parameters: cookieFile and endpoint.
//
// Parameters:
//   - cookieFile (string): This is the path to the file where jar will be saved or loaded from.
//   - endpoint (string): This is the URL endpoint for which the jar are relevant.
//
// The function does the following:
// 1. Creates a new NuvlaCookies instance and sets the cookieFile field.
// 2. Parses the endpoint string into a url.URL object and sets the endpoint field of the NuvlaCookies instance.
// 3. Checks if the cookieFile exists and is not empty. If it is, it creates a new cookiejar.Jar, attempts to load jar from the cookieFile into the cookiejar.Jar, and sets the jar field of the NuvlaCookies instance.
//
// Returns:
//   - A pointer to the newly created NuvlaCookies instance.
//
// Example:
//
//	jar := client.NewNuvlaCookies("/path/to/jar.txt", "http://example.com")
//	In this example, a new NuvlaCookies instance is created. The jar relevant to the "http://example.com" endpoint will be saved to or loaded from the "/path/to/jar.txt" file.
func NewNuvlaCookies(cookieFile string, endpoint string) *NuvlaCookies {
	// Create new NuvlaCookies
	j, _ := cookiejar.New(nil)
	c := &NuvlaCookies{
		cookieFile: cookieFile,
		jar:        j,
	}
	if c.cookieFile == "" {
		c.cookieFile = types.DefaultCookieFile
	}

	// Parse endpoint
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Error("Error parsing endpoint URL")
		return nil
	}
	c.endpoint = u

	// Try to load jar from file
	if common.FileExistsAndNotEmpty(cookieFile) {
		err := c.load()
		if err != nil {
			log.Info("Error loading jar from file")
		}

	}

	return c
}

func (c *NuvlaCookies) load() error {
	log.Debug("Loading cookies from file ...")
	// Open the file
	file, err := os.Open(c.cookieFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		// Get the current line
		line := scanner.Text()

		// Unmarshal the line into a cookie
		var cookie http.Cookie
		err := json.Unmarshal([]byte(line), &cookie)
		if err != nil {
			return err
		}

		// Add the cookie to the jar
		c.jar.SetCookies(c.endpoint, []*http.Cookie{&cookie})
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	log.Debug("Loading cookies from file... success")
	return nil
}

func (c *NuvlaCookies) Save() error {
	// Create folder structure if required
	// TODO: This check probably needs to be moved to a more appropriate place
	err := common.BuildDirectoryStructureIfNotExists(types.DefaultConfigPath)
	if err != nil {
		return err
	}

	// Create target file
	cookieFile, err := os.Create(c.cookieFile)
	if err != nil {
		log.Errorf("Error creating file: %s", err)
		return err
	}
	defer cookieFile.Close()

	// Write jar to file
	log.Infof("Cookies saved to file %s ... ", c.cookieFile)
	for _, cookie := range c.jar.Cookies(c.endpoint) {
		cookieJson, err := json.Marshal(cookie)
		if err != nil {
			log.Errorf("Error marshalling cookie: %s", err)
			return err
		}
		_, err = cookieFile.Write(cookieJson)
		if err != nil {
			log.Errorf("Error writing cookie to file: %s", err)
			return err
		}
		_, err = cookieFile.WriteString("\n")
		if err != nil {
			log.Errorf("Error writing new line to file: %s", err)
			return err
		}
	}

	log.Infof("Cookies saved to file %s ... Success", c.cookieFile)
	return nil
}

// SaveIfNeeded jar if needed
func (c *NuvlaCookies) SaveIfNeeded(newCookie http.CookieJar) error {
	// Get cookies
	newCookies := newCookie.Cookies(c.endpoint)

	// Compare jar
	if !compareCookies(c.lastCookie, newCookies) {
		log.Debugf("Cookies are different, saving new jar: %s", c.cookieFile)
		// If jar are different, save new jar
		c.jar.SetCookies(c.endpoint, newCookies)
		c.lastCookie = newCookies
		return c.Save()
	}
	log.Debugf("Cookies are the same, not saving jar: %s", c.cookieFile)
	return nil

}

func compareCookies(cookies1, cookies2 []*http.Cookie) bool {
	if len(cookies1) != len(cookies2) {
		return false
	}

	for i, cookie := range cookies1 {
		if cookie.String() != cookies2[i].String() {
			return false
		}
	}

	return true
}
