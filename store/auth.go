package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type GoogleUserClaim struct{
    jwt.RegisteredClaims
    Expire int `json:"exp"`
    Email string `json:"email"`
    Verified bool `json:"email_verified"`
    Picture string `json:"picture"`
    Firstname string `json:"given_name"`
    Lastname string `json:"family_name"`
    Name string `json:"name"`
    Locale string `json:"locale"`
}

type GoogleCert struct{
    Cert1 string `json:"a50f6e70ef4b548a5fd9142eecd1fb8f54dce9ee"`
    Cert2 string `json:"28a421cafbe3dd889271df900f4bbf16db5c24d4"`
}

type UserClaims struct{
    ConnectedUser
    jwt.RegisteredClaims
}

type ConnectedUser struct{
    Id string
    ShortName string
    Picture string
}

func SigninUser(email string, password string)(ConnectedUser, error){
    var user ConnectedUser
    var cryptPassword, salt string
    var picture sql.NullString
    conn, err := GetDBConn()
    if err !=  nil{
        log.Println(err)
        return user, errors.New("error conn to the db")
    }
    defer conn.Close()
    userRow := conn.QueryRowContext(context.Background(), `SELECT id, CONCAT(LEFT(firstname, 1), LEFT(lastname, 1)), picture, password, salt FROM Users WHERE email=$1`, email)
    if err := userRow.Scan(user.Id, &user.ShortName, &picture, &cryptPassword, &salt); err != nil{
        log.Println(err)
        return user, errors.New("error selecting users")
    }
    //cypher := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", password,salt,os.Getenv("PASSWORD_SECRET_KEY"))))
    //userPassword := fmt.Sprintf("%x", cypher)
    //if cryptPassword != userPassword{
    //    return user, errors.New("wrong password")
    //}
    user.Picture = picture.String
    return user, nil
}

func SignGoogleUser(user *GoogleUserClaim, token string)error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error db conn")
        return errors.New("error db conn")
    }
    defer conn.Close()
    if _, err := conn.ExecContext(context.Background(), `INSERT INTO Users(firstname, lastname, email, picture, verified, auth_service, auth_expire, token) VALUES($1,$2,$3,$4,$5,$6,$7,$8) 
    ON CONFLICT(email) DO UPDATE SET token=EXCLUDED.token`, user.Firstname, user.Lastname, user.Email, user.Picture, user.Verified, "Google", user.Expire, token); err != nil{
        log.Printf("error in the query: %v", err)
        return errors.New("error in the query")
    }
    return nil 
}

func SigninProUser(email string, password string)(ConnectedUser, error){
    var user ConnectedUser
    var cryptPassword, salt string
    var picture sql.NullString
    conn, err := GetDBConn()
    if err !=  nil{
        log.Println(err)
        return user, errors.New("error conn to the db")
    
    }
    defer conn.Close()
    userRow := conn.QueryRowContext(context.Background(), `SELECT id, LEFT(name, 1), picture, password, salt FROM Entreprise WHERE email=$1`, email)
    if err := userRow.Scan(&user.Id, &user.ShortName, &picture, &cryptPassword, &salt); err != nil{
        log.Println(err)
        return user, errors.New("error selecting users")
    }
    //cypher := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", password,salt,os.Getenv("PASSWORD_SECRET_KEY"))))
    //userPassword := fmt.Sprintf("%x", cypher)
    //if cryptPassword != userPassword{
    //    return user, errors.New("wrong password")
    //}
    user.Picture = picture.String
    return user, nil
}

//Access token duration 6 hour
func CreateToken(data ConnectedUser) (string, error){
    key := []byte(os.Getenv("ACCESS_TOKEN_KEY"))
    claim := UserClaims{data, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(6*time.Hour)), Issuer: "1"}}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tk, err := token.SignedString(key)
    if err != nil{
        return "", errors.New("error signing token")
    }
    return tk+"0", nil
}

//Access token duration 6 hour
func CreateTokenPro(data ConnectedUser) (string, error){
    key := []byte(os.Getenv("PRO_ACCESS_TOKEN_KEY"))
    claim := UserClaims{data, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(6*time.Hour)), Issuer: "1"}}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tk, err := token.SignedString(key)
    if err != nil{
        return "", errors.New("error signing token")
    }
    return tk+"1", nil
}

//Refresh Token duration 5 days
func CreateRefreshToken(data ConnectedUser) (string, error){
    key := []byte(os.Getenv("REFRESH_TOKEN_KEY"))
    claim := UserClaims{data, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour*5)), Issuer: "1"}}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tk, err := token.SignedString(key)
    if err != nil{
        return "", errors.New("error signing token")
    }
    return tk, nil
}

func VerifyToken(token string)(ConnectedUser, string, error){
    tKey, err := jwt.ParseWithClaims(token[:len(token)-1], &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
        if token[len(token)-1:] == "0"{
            return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
        }else{
            return []byte(os.Getenv("PRO_ACCESS_TOKEN_KEY")), nil
        } 
    })
    if err != nil {
        googleUser, _, err := VerifyGoogleToken(token)
        if err != nil{
            return ConnectedUser{}, "", errors.New("error parsing token")
        }
        return googleUser, token, nil
    } else if claims, ok := tKey.Claims.(*UserClaims); ok {
        connedtedUser := ConnectedUser{Id: claims.Id, ShortName: claims.ShortName, Picture: claims.Picture}
        var accessToken string
        if strings.HasSuffix(token, "0"){
            accessToken, _ = CreateToken(connedtedUser)
        }else{
            accessToken, _ = CreateTokenPro(connedtedUser)
        }
        return connedtedUser, accessToken, nil
    } else {
        return ConnectedUser{}, "", errors.New("error token")
    }
}

func IsValidPassword(password string) bool{
    gotNumber := false
    gotLowerCase := false
    gotUpperCase := false
    gotSpecialCharacter := false
    isValid := strings.ContainsFunc(password, func(r rune) bool {
        isUpperCase, _ := regexp.MatchString(`[A-Z]`, string(r))
        isLowerCase, _ := regexp.MatchString(`[a-z]`, string(r))
        hasNumber, _ := regexp.MatchString(`[0-9]`, string(r))
        specialCharacter, _ := regexp.MatchString(`[#?!@$%^&*]`, string(r))
        if isUpperCase{
            gotUpperCase = true
        }else if isLowerCase{
            gotLowerCase = true
        }else if hasNumber{
            gotNumber = true
        }else if specialCharacter{
            gotSpecialCharacter =  true
        }
        
        if gotLowerCase && gotNumber && gotUpperCase && gotSpecialCharacter{
            return true
        }else{
            return false
        }
    })
    if len(password) > 6 && isValid{
        return true
    }else{
        return false
    }
}

func VerifyGoogleToken(googleToken string)(ConnectedUser, *GoogleUserClaim,  error){
    var googleCert GoogleCert
    var user ConnectedUser
    var token *jwt.Token
    requestCert, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/certs", nil)
    if err != nil{
        log.Printf("Error in the request google key: %v",err)
        return user, nil, errors.New("error in the request to google")
    }
    responseCert, err := http.DefaultClient.Do(requestCert)

    decCert := json.NewDecoder(responseCert.Body)
    decCert.Decode(&googleCert)
    if err != nil{
        log.Printf("error doing the request: %v", err)
        return user, nil, errors.New("error doing the request")
    }
    rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(googleCert.Cert1))
    rsaPublicKey1, err := jwt.ParseRSAPublicKeyFromPEM([]byte(googleCert.Cert2))
    if err != nil{
        log.Printf("error parsing public key: %v", err)
        return user, nil, errors.New("error parsing the public key")
    }
    token, err = jwt.ParseWithClaims(googleToken, &GoogleUserClaim{}, func(t *jwt.Token) (interface{}, error) {
        return rsaPublicKey, nil 
    }, jwt.WithValidMethods([]string{"RS256"}))
    if err != nil{
        token, err = jwt.ParseWithClaims(googleToken, &GoogleUserClaim{}, func(t *jwt.Token) (interface{}, error) {
            return rsaPublicKey1, nil 
        }, jwt.WithValidMethods([]string{"RS256"}))
        if err != nil{
            log.Printf("error verifying token? %v", err)
            return user, nil, errors.New("error verifying token")
        }
    }
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error in the db: %v", err)
        return user, nil, errors.New("error conn to the db")
    }
    googleUser := token.Claims.(*GoogleUserClaim)
    userRow := conn.QueryRowContext(context.Background(), `SELECT id FROM Users WHERE email=$1`, googleUser.Email)
    if err := userRow.Scan(&user.Id); err != nil{
        log.Printf("error scanining id in verify token: %v", err)
        return user, nil, errors.New("error scanining id in verify token")
    }
    user.ShortName = fmt.Sprintf("%v%v", googleUser.Firstname[:1], googleUser.Lastname[:1])
    user.Picture = googleUser.Picture

    return user, googleUser, nil
}





















