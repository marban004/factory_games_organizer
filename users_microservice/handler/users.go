package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	custommiddleware "github.com/marban004/factory_games_organizer/custom_middleware"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/model"
	"github.com/marban004/factory_games_organizer/microservice_logic_users/repository/user"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/crypto/bcrypt"
)

type JSONData struct {
	UserLogin    string
	UserPassword string
}

type Users struct {
	UserRepo    *user.MySQLRepo
	Secret      []byte
	StatTracker *custommiddleware.DefaultApiStatTracker
}

type CreateUserResponse struct {
	UsersCreated uint
}

type UpdateUserResponse struct {
	UsersUpdated uint
}

type LoginResponse struct {
	Jwt string
}

type DeleteUserResponse struct {
	UsersDeleted uint
}

type HealthResponse struct {
	MicroserviceStatus string
	DatabaseStatus     string
}

type StatsResponse struct {
	ApiUsageStats    *orderedmap.OrderedMap[string, map[string]int]
	TrackingPeriodMs int64
	NoPeriods        uint64
}

// Health return the status of microservice and associated database
//
//	@Description	Return the status of microservice and it's database. Default working state is signified by status "up".
//	@Tags			Users
//	@Success		200	{object}	handler.HealthResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/health [get]
func (h *Users) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		MicroserviceStatus: "up",
	}
	err := h.UserRepo.DB.PingContext(r.Context())
	if err != nil {
		response.DatabaseStatus = "connection disrupted"
	} else {
		response.DatabaseStatus = "up"
	}
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// Stats return the usage stats of microservice
//
//	@Description	Return the usage stats of microservice.
//	@Tags			Users
//	@Success		200	{object}	handler.StatsResponse
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/stats [get]
func (h *Users) Stats(w http.ResponseWriter, r *http.Request) {
	endpointResponse := StatsResponse{ApiUsageStats: h.StatTracker.GetStats(), TrackingPeriodMs: h.StatTracker.Period, NoPeriods: h.StatTracker.MaxLen}
	byteJSONRepresentation, err := json.Marshal(endpointResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// CreateUser create new user
//
//	@Description	Insert data of new user into database. Logins of every user must be unique. Passwords must be at least 8 characters long and maximum 72 characters long. Passwords must contain a lowercase and uppercase letter, a digit and a special character that is not a space, quote, double quote or semicolon. Logins must be at least 3 characters long and maximum 64 characters long. Logins cannot contain a space, quote, double quote or semicolon. Logins ignore letter case when logging in.
//	@Param			createUser	body	handler.JSONData	true	"New user data to be inserted into database"
//	@Tags			Users
//
//	@Accept			json
//
//	@Success		200	{object}	handler.CreateUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [post]
func (h *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	//no parameters are required for this request
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
	response := CreateUserResponse{}
	response.UsersCreated = uint(noRows)
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("user has been created, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(byteJSONRepresentation)
}

// UpdateUser update user's data
//
//	@Description	Update user's data in database. The user whose data is updated is the user who presented the authentication token. Same login and password rules apply as when creating a new user account.
//	@Param			updateUser	body	handler.JSONData	true	"New data of the user to be saved into database"
//	@Tags			Users Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.UpdateUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [put]
//
//	@Security		apiTokenAuth
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
		w.WriteHeader(http.StatusUnauthorized)
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
	response := UpdateUserResponse{}
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

// LoginUser login users
//
//	@Description	Authenticate users against the database. If verification is successfull a jwt(authentication token) is returned, that can be used to prove the user's identity to other microservices in Factory Games Organizer api.
//	@Param			login	body	handler.JSONData	true	"Login data for the user."
//	@Tags			Users
//
//	@Accept			json
//
//	@Success		200	{object}	handler.LoginResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format or invalid login data has been sent"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/login [post]
func (h *Users) LoginUser(w http.ResponseWriter, r *http.Request) {
	// no parameters are required for this request
	inputData := JSONData{}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("could not parse received body, reason: %w", err).Error()))
		return
	}
	user, err := h.UserRepo.SelectUserByLogin(r.Context(), inputData.UserLogin)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such user"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not retrieve user data: %w", err).Error()))
		return
	}
	if !h.checkPassword(inputData.UserPassword, user.UserPasswdHash) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid credentials"))
		return
	}
	token, err := h.createToken(int(user.UserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not generate authentication token: %w", err).Error()))
	}
	response := LoginResponse{}
	response.Jwt = token
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("user has been updated, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

// DeleteUser delete user's data
//
//	@Description	Delete user's data in database. The user whose data is deleted is the user who presented the authentication token.
//	@Tags			Users Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.DeleteUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/ [delete]
//
//	@Security		apiTokenAuth
func (h *Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	result, err := h.UserRepo.DeleteUser(r.Context(), uint(userId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("could not delete requested user, reason: %w", err).Error()))
		return
	}
	noRows, err := result.RowsAffected()
	if err != nil {
		w.Write([]byte("database driver does not support returning numbers of rows affected"))
	}
	response := DeleteUserResponse{}
	response.UsersDeleted = uint(noRows)
	byteJSONRepresentation, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("user has been deleted, but could not generate json representation of response, reason: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteJSONRepresentation)
}

func (h *Users) createToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Minute * 30).Unix(),
			"iat":    time.Now().Unix(),
		})
	tokenString, err := token.SignedString(h.Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
