package users

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// DeleteUserURI is pattern for the endpoint to delete a user
	DeleteUserURI = "/api/users/%s"
)

// DeleteUserList deletes a list of users from the Platform
func DeleteUserList(users []string, config control.Configuration, debug bool) error {

	if debug {
		log.Printf("Deleting %v users...", len(users))
	}

	client, _, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + DeleteUserURI

	for _, v := range users {
		uri := fmt.Sprintf(url, v)
		if debug {
			log.Println("DELETE", uri)
		}
		req, err := http.NewRequest("DELETE", uri, nil)
		if err != nil {
			return err
		}
		req = control.AddCsrfHeader(req, client)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if debug {
			log.Println(resp.Status)
		}

		if resp.StatusCode == 200 {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("User %s deleted (%s).\n", bodyBytes, v)
		} else {
			return errors.New("Unable to delete user " + v + " (" + resp.Status + ")")
		}

	}

	return nil
}
