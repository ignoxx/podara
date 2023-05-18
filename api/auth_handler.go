package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	// create user
	bodyJson := struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal(body, &bodyJson)

	if err != nil {
		return err
	}

	if bodyJson.Email == "" || bodyJson.Name == "" || bodyJson.Password == "" {
		return errors.New("invalid data")
	}

	// save user
	user, err := s.store.CreateUser(bodyJson.Email, bodyJson.Name, bodyJson.Password)

	if err != nil {
		return err
	}

	signedToken, err := createJwtToken(user)

	if err != nil {
		return err
	}

	w.Header().Add("Authorization", signedToken)
	return WriteJSON(w, http.StatusOK, user)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// read request body as json
	body, err := ioutil.ReadAll(r.Body)

	// check for errors
	if err != nil {
		return err
	}

	defer r.Body.Close()

	bodyJson := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal(body, &bodyJson)

	if err != nil {
		return err
	}

	if bodyJson.Email == "" || bodyJson.Password == "" {
		return errors.New("invalid data")
	}

	user, err := s.store.GetUserByEmail(bodyJson.Email)
	if err != nil {
		return err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyJson.Password))
	if err != nil {
		return err
	}

	signedToken, err := createJwtToken(user)

	if err != nil {
		return err
	}

	w.Header().Add("Authorization", signedToken)

	return WriteJSON(w, http.StatusOK, user)
}
