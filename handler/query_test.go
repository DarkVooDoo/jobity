package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type Recom struct{
    Label string
    Postal string
    Contract string
    Fulltime bool
}

func TestContext(t *testing.T){
    ctx := context.Background()

    myContext := context.WithValue(ctx, "db", "Hello World")
    
    select{
    case <-time.After(10 * time.Second):
        log.Printf("10 Seconds")
    case <-myContext.Done():
        log.Printf("done")

    }
}

func TestMostCommunQuery(t *testing.T){
    data := []Recom{{"Vendeur", "94", "CDI", true}}

    GetTest(data, 0)

    //Vendeur(se) temps partiel (H/F)
    //Halloween NIGLOLAND : Vendeur en boutique  (H/F)
    //Vendeur / Vendeuse en horlogerie                            (H/F)
    //Vendeur comptoir en centre auto (H/F)
    //Vendeur / Vendeuse en charcuterie Saint cloud (92) (H/F)
    //Agent préparation commandes F/H

}

func GetTest(arr []Recom, pos int){
    field := -1
    var resultStruct Recom
    var lengthRecommendation = len(arr)
    for row := 0; row < 4; row++{
        field++
        var matrix [4][15]string
        //matrix := make([][]string, lengthRecommendation, 15)
        for t := 0; t < len(arr); t++{
            currentField := reflect.ValueOf(arr[t]).Field(row)
            if field == 3{
                log.Printf("BOol: %v", strconv.FormatBool(currentField.Bool()))
            }else{
                for column := 0; column < currentField.Len(); column++{
                    //value += currentField.String()[column:column+1]
                    matrix[t][column] = currentField.String()[column:column+1]
                } 
            }
        }
        result := ""
        for c := 0; c < 15; c++{
            mapCount := map[string]int{}
            answer := ""
            count := 0
            for r := 0; r < lengthRecommendation; r++{
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
    log.Println(resultStruct)

}

func TestCalculateDistance(t *testing.T){
    result := 1275.6
    var r float64 = 6371
    lat1 := 48.8396952
    lon1 := 2.2399123
    lat2 := 48.918167
    lon2 := 2.292385
    
    q1 := lat1 * math.Pi / 180
    q2 := lat2 * math.Pi / 180
    y1 := (lat2 - lat1) * math.Pi / 180
    y2 := (lon2 - lon1) * math.Pi / 180
    a := math.Sin(y1/2) * math.Sin(y1/2) + math.Cos(q1) * math.Cos(q2) * math.Sin(y2/2) * math.Sin(y2/2)  
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    mettres := r * c 
    if fmt.Sprintf("%.1f",mettres) == fmt.Sprintf("%.1f",result){
        log.Printf("Correct: %.1f", mettres)
    }else{
        t.Fatalf("Incorrect: result is: %v \tneed to be %v", mettres, result)
    }
}
