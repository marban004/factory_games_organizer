package prototypes

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
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

func CreateToken(userId int) (string, error) {
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

// logins must be unique
func VerifyUser(ctx context.Context, db *sql.DB, login string, passwd string) (uint, error) {
	user := User{}
	query := fmt.Sprintf(`SELECT * FROM users where login = "%s"`, login)
	err := db.QueryRowContext(ctx, query).Scan(&user)
	if err != nil {
		return 0, fmt.Errorf("could not verify user: %w", err)
	}
	if CheckPassword(passwd, user.UserPasswdHash) {
		return user.UserId, nil
	}
	return 0, fmt.Errorf("invalid credentials")
}

func CreateUser(ctx context.Context, db *sql.DB, login string, passwd string) (uint, error) {
	user := User{}
	query := fmt.Sprintf(`SELECT * FROM users where login = "%s"`, login)
	err := db.QueryRowContext(ctx, query).Scan(&user)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("this login is already in use, provide a different login")
	}
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	matched, err := regexp.Match(`^.{8,72}$"`, []byte(passwd))
	if !matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*[a-z].*$"`, []byte(passwd))
	if !matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*[0-9].*$"`, []byte(passwd))
	if !matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*[!@#$&*].*"`, []byte(passwd))
	if !matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*[A-Z].*$"`, []byte(passwd))
	if !matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*["';].*$"`, []byte(passwd))
	if matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	matched, err = regexp.Match(`^.*((?i)(select)|(?i)(update)|(?i)(insert)|(?i)(delete)|(?i)(from)).*$"`, []byte(passwd))
	if matched {
		return 0, fmt.Errorf("password is invalid, provide a different password")
	}
	if err != nil {
		return 0, fmt.Errorf("could not process provided password: %w", err)
	}
	hash, err := GeneratePasswordHash(passwd)
	if err != nil {
		return 0, fmt.Errorf("could not generate password hash: %w", err)
	}
	query = fmt.Sprintf(`INSERT INTO users VALUES (null, %s, %s)`, login, hash)
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("could not insert new user: %w", err)
	}
	userId, err := VerifyUser(ctx, db, login, passwd)
	if err != nil {
		return 0, fmt.Errorf("could not verify created user: %w", err)
	}
	return userId, nil
}
