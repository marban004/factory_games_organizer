package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/model"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/repository/user"
	"golang.org/x/crypto/bcrypt"
)

type JSONData struct {
	UserLogin    string
	UserPassword string
}

type Users struct {
	UserRepo *user.MySQLRepo
	Secret   []byte
}

func (h *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	//no parameters aree required for this request
	inputData := JSONData{}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	valid, err := h.verifyUserLogin(inputData.UserLogin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server could not resolve regex pattern, contact server administrator"))
		return
	}
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Provided login is invalid. Login needs to be minimum 3 characters long, maximum 64 characters long and cannot contain ""","'" or ";" characters`))
		return
	}
	valid, err = h.verifyUserPassword(inputData.UserPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server could not resolve regex pattern, contact server administrator"))
		return
	}
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Provided password is invalid. Password needs to be minimum 8 characters long, maximum 72 characters long, needs to contain a lowercase letter, an uppercase letter, a digit, a special character and cannot contain " ", """, "'" or ";" characters`))
		return
	}
	hash, err := h.generatePasswordHash(inputData.UserPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server could not generate password hash, contact server administrator"))
		return
	}
	result, err := h.UserRepo.CreateUser(r.Context(), model.UserInfo{UserId: 0, UserLogin: inputData.UserLogin, UserPasswdHash: hash})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not create requested user, reason: %w", err).Error()))
		return
	}
	noRows, err := result.RowsAffected()
	if err != nil {
		w.Write([]byte("database driver does not support returning numbers of rows affected"))
	}
	var response struct {
		UsersCreated uint
	}
	response.UsersCreated = uint(noRows)
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("user has been created, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

func (h *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//parameters for request are:
	//jwt = token with dispatcher server secret key, id of user who received the token and issue date of the token, not optional
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, userId := h.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	inputData := JSONData{}
	hash := ""
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	if len(inputData.UserLogin) > 0 {
		valid, err = h.verifyUserLogin(inputData.UserLogin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server could not resolve regex pattern, contact server administrator"))
			return
		}
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`Provided login is invalid. Login needs to be minimum 3 characters long, maximum 64 characters long and cannot contain ""","'" or ";" characters`))
			return
		}
	}
	if len(inputData.UserPassword) > 0 {
		valid, err = h.verifyUserPassword(inputData.UserPassword)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server could not resolve regex pattern, contact server administrator"))
			return
		}
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`Provided password is invalid. Password needs to be minimum 8 characters long, maximum 72 characters long, needs to contain a lowercase letter, an uppercase letter, a digit, a special character and cannot contain " ", """, "'" or ";" characters`))
			return
		}
		hash, err = h.generatePasswordHash(inputData.UserPassword)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server could not generate password hash, contact server administrator"))
			return
		}
	}
	result, err := h.UserRepo.UpdateUser(r.Context(), model.UserInfo{UserId: uint(userId), UserLogin: inputData.UserLogin, UserPasswdHash: hash})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not update requested user, reason: %w", err).Error()))
		return
	}
	noRows, err := result.RowsAffected()
	if err != nil {
		w.Write([]byte("database driver does not support returning numbers of rows affected"))
	}
	var response struct {
		UsersUpdated uint
	}
	response.UsersUpdated = uint(noRows)
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("user has been updated, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// todo: implement verification of jwt
func (h *Users) verifyJWT(jwtString string) (bool, int) {
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

// func (h *Users) convertArrToInt(input []string) []int {
// 	result := []int{}
// 	for _, value := range input {
// 		intValue, err := strconv.Atoi(value)
// 		if err != nil {
// 			continue
// 		}
// 		result = append(result, intValue)
// 	}
// 	return result
// }

func (h *Users) generatePasswordHash(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 12)
	if err != nil {
		return "", fmt.Errorf("could not generate hash of provided string: %w", err)
	}
	return string(bytes), nil
}

func (h *Users) checkPassword(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

// returns true and nil if provided login is valid, returns false and nil if provided login is invalid, returns false and an error if an error occured while matching login to regex
func (h *Users) verifyUserLogin(login string) (bool, error) {
	valid, err := h.verifyAgainstRegexMatches(login, `^.{3,64}$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexNotMatches(login, `^.*((?i)(select)|(?i)(update)|(?i)(insert)|(?i)(delete)|(?i)(drop)).*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexNotMatches(login, `^.*[\ "';].*$`)
	if !valid || err != nil {
		return valid, err
	}
	return valid, err
}

// returns true and nil if provided login is valid, returns false and nil if provided login is invalid, returns false and an error if an error occured while matching password to regex
func (h *Users) verifyUserPassword(passwd string) (bool, error) {
	valid, err := h.verifyAgainstRegexMatches(passwd, `^.{8,72}$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexMatches(passwd, `^.*[a-z].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexMatches(passwd, `^.*[A-Z].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexMatches(passwd, `^.*[0-9].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexMatches(passwd, `^.*[!@#$&*].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexNotMatches(passwd, `^.*[\ "';].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = h.verifyAgainstRegexNotMatches(passwd, `^.*((?i)(select)|(?i)(update)|(?i)(insert)|(?i)(delete)|(?i)(drop)).*$`)
	if !valid || err != nil {
		return valid, err
	}
	return valid, err
}

// returns false if string does not match regex pattern or an error occured while parsing pattern in which case also returns said error, otherwise returns true and nil, if pattern does not match and no error occurred returns false and nil
func (h *Users) verifyAgainstRegexMatches(text string, pattern string) (bool, error) {
	matched, err := regexp.Match(pattern, []byte(text))
	if !matched {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not parse provided pattern: %w", err)
	}
	return true, nil
}

// returns false if string matches regex pattern or an error occured while parsing pattern in which case also returns said error, otherwise returns true and nil, if pattern matchea and no error occurred returns false and nil
func (h *Users) verifyAgainstRegexNotMatches(text string, pattern string) (bool, error) {
	matched, err := regexp.Match(pattern, []byte(text))
	if matched {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not parse provided pattern: %w", err)
	}
	return true, nil
}
