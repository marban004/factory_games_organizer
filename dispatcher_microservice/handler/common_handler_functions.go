//     This is Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.
//     Copyright (C) 2025  Marek Banaś

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see https://www.gnu.org/licenses/.

package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CommonHandlerFunctions struct {
	Secret           []byte
	NextMicroservice uint
	Client           *http.Client
}

func (h *CommonHandlerFunctions) useNextMicroservice(len uint) {
	if len == 1 {
		return
	}
	if h.NextMicroservice+1 == len {
		h.NextMicroservice = 0
		return
	}
	h.NextMicroservice += 1
}

func (h *CommonHandlerFunctions) redirectRequest(w http.ResponseWriter, r *http.Request, requestEndpoint string, microserviceAddressArray []string) {
	redirectURI := fmt.Sprintf("https://%s/%s", microserviceAddressArray[h.NextMicroservice], requestEndpoint)
	_, params, paramsPresent := strings.Cut(r.RequestURI, "?")
	if paramsPresent {
		redirectURI += "?" + params
	}
	request, err := http.NewRequestWithContext(r.Context(), r.Method, redirectURI, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create request to microservice"))
		return
	}
	response, err := h.Client.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not communicate with microservice"))
		return
	}
	w.WriteHeader(response.StatusCode)
	temp := make([]byte, 1)
	for {
		_, err = response.Body.Read(temp)
		if err == io.EOF {
			break
		}
		w.Write(temp)
	}
	h.useNextMicroservice(uint(len(microserviceAddressArray)))
}

// todo: implement verification of jwt
func (h *CommonHandlerFunctions) verifyJWT(jwtString string) (bool, int) {
	token, err := jwt.Parse(jwtString, func(*jwt.Token) (interface{}, error) {
		return h.Secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return false, 0
	}
	if !token.Valid {
		return false, 0
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	expTime := int64(claims["exp"].(float64))
	issueTime := int64(claims["iat"].(float64))
	if time.Now().Unix() > expTime {
		return false, 0
	}
	if time.Now().Unix() <= issueTime {
		return false, 0
	}
	return true, int(userId)
}
