package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
    userRow := conn.QueryRowContext(context.Background(), `SELECT id, CONCAT(LEFT(firstname, 1), LEFT(lastname, 1)), picture, password, salt FROM Users WHERE email=$1`, email)
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

func SigninProUser(email string, password string)(ConnectedUser, error){
    var user ConnectedUser
    var cryptPassword, salt string
    var picture sql.NullString
    conn, err := GetDBConn()
    if err !=  nil{
        log.Println(err)
        return user, errors.New("error conn to the db")
    
    }
    userRow := conn.QueryRowContext(context.Background(), `SELECT id, name, picture, password, salt FROM Entreprise WHERE email=$1`, email)
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
    tk, err := jwt.ParseWithClaims(token[:len(token)-1], &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
        if token[len(token)-1:] == "0"{
            return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
        }else{
            return []byte(os.Getenv("PRO_ACCESS_TOKEN_KEY")), nil
        }
    })
    if err != nil {
        return ConnectedUser{}, "", errors.New("error parsing token")
    } else if claims, ok := tk.Claims.(*UserClaims); ok {
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
        specialCharacter, _ := regexp.MatchString(`[#?!@$%^&*-]`, string(r))
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
