package prototypes

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId         uint
	UserLogin      string
	UserPasswdHash string
}

func GeneratePasswordHash(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 12)
	if err != nil {
		return "", fmt.Errorf("could not generate hash of provided string: %w", err)
	}
	return string(bytes), nil
}

func CheckPassword(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

func CreateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Minute * 30).Unix(),
			"iat":    time.Now().Unix(),
		})
	fileContents, err := os.ReadFile("crud_microservice_secret.pem")
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(fileContents)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// logins must be unique, put in User repository, 0 means user not verified as user id in database starts from 1
func VerifyUser(ctx context.Context, db *sql.DB, login string, passwd string) (uint, error) {
	user := User{}
	query := fmt.Sprintf(`SELECT * FROM users where login = "%s"`, strings.ToLower(login))
	err := db.QueryRowContext(ctx, query).Scan(&user.UserId, &user.UserLogin, &user.UserPasswdHash)
	if err != sql.ErrNoRows && err != nil {
		return 0, fmt.Errorf("could not verify user: %w", err)
	}
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if CheckPassword(passwd, user.UserPasswdHash) {
		return user.UserId, nil
	}
	return 0, nil
}

// put in User repository, only creates does not verify validity of login and password
func CreateUser(ctx context.Context, db *sql.DB, user User) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO users VALUES (null, "%s", "%s")`, strings.ToLower(user.UserLogin), user.UserPasswdHash)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not insert new user: %w", err)
	}
	return result, nil
}

// put in User repository, does not verify validity of login and password, replace UserId in user with id from jwt
func UpdateUser(ctx context.Context, db *sql.DB, user User) (sql.Result, error) {
	query := `UPDATE users SET `
	if len(user.UserLogin) > 0 {
		query += fmt.Sprintf(`login = "%s" `, strings.ToLower(user.UserLogin))
	}
	if len(user.UserPasswdHash) > 0 {
		query += fmt.Sprintf(`passwdhash = "%s" `, user.UserPasswdHash)
	}
	query += fmt.Sprintf(`WHERE id = %d`, user.UserId)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("user has not been updated: %w", err)
	}
	return result, nil
}

func DeleteUser(ctx context.Context, db *sql.DB, userId uint) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM users WHERE id = %d", userId)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func VerifyUserLogin(login string) (bool, error) {
	valid, err := VerifyAgainstRegexMatches(login, `^.{3,64}$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexNotMatches(login, `^.*((?i)(select)|(?i)(update)|(?i)(insert)|(?i)(delete)|(?i)(drop)).*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexNotMatches(login, `^.*["';].*$`)
	if !valid || err != nil {
		return valid, err
	}
	return valid, err
}

func VerifyUserPassword(passwd string) (bool, error) {
	valid, err := VerifyAgainstRegexMatches(passwd, `^.{8,72}$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexMatches(passwd, `^.*[a-z].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexMatches(passwd, `^.*[A-Z].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexMatches(passwd, `^.*[0-9].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexMatches(passwd, `^.*[!@#$&*].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexNotMatches(passwd, `^.*["';].*$`)
	if !valid || err != nil {
		return valid, err
	}
	valid, err = VerifyAgainstRegexNotMatches(passwd, `^.*((?i)(select)|(?i)(update)|(?i)(insert)|(?i)(delete)|(?i)(drop)).*$`)
	if !valid || err != nil {
		return valid, err
	}
	return valid, err
}

func VerifyAgainstRegexMatches(text string, pattern string) (bool, error) {
	matched, err := regexp.Match(pattern, []byte(text))
	if !matched {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not parse provided pattern: %w", err)
	}
	return true, nil
}

func VerifyAgainstRegexNotMatches(text string, pattern string) (bool, error) {
	matched, err := regexp.Match(pattern, []byte(text))
	if matched {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not parse provided pattern: %w", err)
	}
	return true, nil
}

func VerifyJWT(jwtString string) (bool, int) {
	token, err := jwt.Parse(jwtString, func(*jwt.Token) (interface{}, error) {
		return "replace with secret key later", nil
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
