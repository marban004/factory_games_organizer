package handler

import (
	"fmt"
	"io"
	"net/http"
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
	request, err := http.NewRequestWithContext(r.Context(), r.Method, fmt.Sprintf("https://%s/%s", microserviceAddressArray[h.NextMicroservice], requestEndpoint), r.Body)
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
