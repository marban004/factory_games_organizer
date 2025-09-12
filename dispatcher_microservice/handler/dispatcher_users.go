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
	"net/http"
)

type DispatcherUsers struct {
	CommonHandlerFunctions      CommonHandlerFunctions
	UsersMicroservicesAddresses []string
}

// LoginUser login users
//
//	@Description	Authenticate users against the database. If verification is successfull a jwt(authentication token) is returned, that can be used to prove the user's identity to other microservices in Factory Games Organizer api.
//	@Param			login	body	handler.JSONDataUsers	true	"Login data for the user."
//	@Tags			Users
//
//	@Accept			json
//
//	@Success		200	{object}	handler.LoginResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format or invalid login data has been sent"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/users/login [post]
func (h *DispatcherUsers) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.CommonHandlerFunctions.redirectRequest(w, r, "login", h.UsersMicroservicesAddresses)
}

// CreateUser create new user
//
//	@Description	Insert data of new user into database. Logins of every user must be unique. Passwords must be at least 8 characters long and maximum 72 characters long. Passwords must contain a lowercase and uppercase letter, a digit and a special character that is not a space, quote, double quote or semicolon. Logins must be at least 3 characters long and maximum 64 characters long. Logins cannot contain a space, quote, double quote or semicolon. Logins ignore letter case when logging in.
//	@Param			createUser	body	handler.JSONDataUsers	true	"New user data to be inserted into database"
//	@Tags			Users
//
//	@Accept			json
//
//	@Success		200	{object}	handler.CreateUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/users [post]
func (h *DispatcherUsers) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.UsersMicroservicesAddresses)
}

// UpdateUser update user's data
//
//	@Description	Update user's data in database. The user whose data is updated is the user who presented the authentication token. Same login and password rules apply as when creating a new user account.
//	@Param			updateUser	body	handler.JSONDataUsers	true	"New data of the user to be saved into database"
//	@Tags			Users Authorization required
//
//	@Accept			json
//
//	@Success		200	{object}	handler.UpdateUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/users [put]
//
//	@Security		apiTokenAuth
func (h *DispatcherUsers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.UsersMicroservicesAddresses)
}

// DeleteUser delete user's data
//
//	@Description	Delete user's data in database. The user whose data is deleted is the user who presented the authentication token.
//	@Tags			Users Authorization required
//
//	@Success		200	{object}	handler.DeleteUserResponse
//	@Failure		400	{string}	string	"Bad request. One of required parameters is missing or is not of valid format"
//	@Failure		401	{string}	string	"Authentication error"
//	@Failure		500	{string}	string	"Unexpected serverside error"
//	@Router			/users [delete]
//
//	@Security		apiTokenAuth
func (h *DispatcherUsers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if len(jwt) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("jwt parameter cannot be empty"))
		return
	}
	valid, _ := h.CommonHandlerFunctions.verifyJWT(jwt)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("provided jwt is invalid"))
		return
	}
	h.CommonHandlerFunctions.redirectRequest(w, r, "", h.UsersMicroservicesAddresses)
}
