package lib

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"reflect"
	"slices"
)

var Contract map[string]string = map[string]string{"SAI": "Saisonnier", "MIS": "Intérim", "Saisonnier": "Saisonnier", "CCE": "CCE", "FRA": "FRA", "LIB": "LIB", "REP": "REP", "TTI": "TTI", "DDI": "DDI", "DIN": "DIN", "Alternance": "Alternance", "CDD": "CDD", "CDI": "CDI"}

type QueueEnqueueterface[K any] interface{
    Enqueue(K)
    Dequeue()K
    Front()K
    Rear()K
    IsFull()bool
    Queue()[]K
    Size()int
}

type Queue[K any] struct{
    List []K
    Length int
}

type RecomendationToken struct{
    Label string `json:"label"`
    Postal string `json:"postal"`
    Contract string `json:"contract"`
    Fulltime bool `json:"fulltime"`
}

func (q *Queue[K]) Enqueue(value K){
    if q.Length > len(q.List){
        q.List = append(q.List, value)
    }
}

func (q *Queue[K])Dequeue()K{
    out := q.List[0]
    q.List = slices.Delete(q.List, 0, 1)
    return out
}

func (q *Queue[K])Front()K{
    return q.List[0]
}

func (q *Queue[K])Rear()K{
    return q.List[len(q.List)-1]
}

func (q *Queue[K])IsFull()bool{
    if q.Length == len(q.List) {
        return true
    }else{
        return false
    }
}

func (q *Queue[K]) Queue()[]K{
    return q.List
}

func (q *Queue[K])Size()int{
    return len(q.List)
}

func MostSearch(storedSearch []RecomendationToken)RecomendationToken{
    field := -1
    lengthStore := len(storedSearch)
    var resultStruct RecomendationToken
    for row := 0; row < 4; row++{
        field++
        var matrix [5][15]string
        for t := 0; t < len(storedSearch); t++{
            currentField := reflect.ValueOf(storedSearch[t]).Field(row)
            if field == 3{
                //Handle bool 
            }else{
                for column := 0; column < currentField.Len(); column++{
                    matrix[t][column] = currentField.String()[column:column+1]
                } 
            }
        }
        result := ""
        for c := 0; c < 15; c++{
            mapCount := map[string]int{}
            answer := ""
            count := 0
            for r := 0; r < lengthStore; r++{
                mapCount[matrix[r][c]] = mapCount[matrix[r][c]] + 1
                if count < mapCount[matrix[r][c]]{
                    answer = matrix[r][c]
                    count = mapCount[matrix[r][c]]
                }
            }
            result += answer
            
        }
        if field == 0{
            resultStruct.Label = result
        }else if field == 1{
            resultStruct.Postal = result
        }else if field == 2{
            resultStruct.Contract = result
        }
    }
    return resultStruct
}

func ReadBody(body io.ReadCloser) map[string]string {
	var result map[string]string = map[string]string{}
	buf, _ := io.ReadAll(body)
	payload := bytes.Split(buf, []byte("&"))
    if len(payload[0]) > 0{
        for _, value := range payload {
            keyValue := bytes.Split(value, []byte("="))
            key, _ := url.QueryUnescape(string(keyValue[0]))
            value, _ := url.QueryUnescape(string(keyValue[1]))
            result[key] = value
        }    
    }
	return result
}

func SendMail(to []string, subject string, body string) error{
    msg := []byte("Subject: "+subject+ "\n"+
    "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
    body)
    auth := smtp.PlainAuth("", "moisestazaro@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
    if err := smtp.SendMail("smtp.gmail.com:587", auth, "test", to, msg); 
    err != nil{
        log.Println(err)
        return errors.New("error sending mail")
    }
    return nil
}
