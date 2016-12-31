package main

import (
  "os"
  "fmt"
  "time"
  "errors"
  "github.com/dgrijalva/jwt-go"
  "github.com/nu7hatch/gouuid"
)

type AuthClaims struct {
  User string `json:"user"`
  jwt.StandardClaims
}

var secretKey interface{} = []byte("AuthKey?")

func getSecretKey(token *jwt.Token) (interface{}, error) {
  return secretKey, nil
}

func CreateAuthToken(user string) (string, error) {
  standardClaims := jwt.StandardClaims{
    ExpiresAt: time.Now().Add(time.Hour).Unix(),
    Issuer:    "test",
  }

  claims := AuthClaims{
    user,
    standardClaims,
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  return token.SignedString(secretKey)
}

func HasValidAuthToken(v string) bool {
  token, err := jwt.ParseWithClaims(v, &AuthClaims{}, getSecretKey)

  if err != nil {
    fmt.Println(err)
    return false
  }

  if _, ok := token.Claims.(*AuthClaims); ok && token.Valid {
    return true
  } else {
    return false
  }
}

func GenerateHash() (*uuid.UUID, error) {
  return uuid.NewV4()
}

func IsValidHash(h string) bool {
  return h == "12345"
}

func SendAuthEmail(recipient string) error {
  pass := os.Getenv("EMAIL_PASS")

  if pass != "" {
    return actuallySendAuthEmail(recipient)
  }

  return errors.New("Failed to send email")
}

func actuallySendAuthEmail(rec string) error {
  sender := NewSender("k@blwsk.com", os.Getenv("EMAIL_PASS"))

  recipients := []string{rec}

  u, err := GenerateHash()

  if err != nil {
    return err
  }

  hash := u.String()
  subject := hash
  body := "try this: http://localhost:8080/auth/" + hash

  return sender.SendMail(recipients, subject, body + hash)
}