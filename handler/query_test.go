package handler

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/jung-kurt/gofpdf"
)

type Recom struct{
    Label string
    Postal string
    Contract string
    Fulltime bool
}

type Work struct{
    Position string `json:"position"`
    Entreprise string `json:"entreprise"`
    StartDate string `json:"start_date"` 
    EndDate string `json:"end_date"` 
    Description []string `json:"description"`
}

type School struct{
    Name string `json:"name"`
    Establishment string `json:"establishment"`
    EndDate string `json:"end_date"` 
    StartDate string `json:"start_date"` 
}

const (
    FONT_NAME="SUSE"
    LEFT_SECTION=.35
    RIGHT_SECTION=.65
    LEFT_MARGIN=8
)

func TestCreateCurriculumPDF(t *testing.T){
    interest := []string{"Golang", "Javascript"}
    skill := []string{"Musique", "True baller", "Moto"}
    myWorks := []Work{
        Work{Position: "Preparateur de commande", StartDate: "22-03-2022", Entreprise: "Amazon", EndDate: "30-06-2024", Description: []string{"Test", "Long description dsq k kdsq jsq, sq,djqsjd,q", "no way its home"}},
        Work{Position: "Preparateur de commande", StartDate: "Jun-2022", Entreprise: "Amazon", EndDate: "Avr-2024", Description: []string{"Test", "Long description dsq k kdsq jsq, sq,djqsjd,q", "no way its home"}},
        Work{Position: "Preparateur de commande", StartDate: "Nov-2022", Entreprise: "Amazon", EndDate: "Dec-2024", Description: []string{"Test", "Long description dsq k kdsq jsq, sq,djqsjd,q cuanto te vas no morire de amor ella se burla de ti de mi", "no way its home"}},

    }
    mySchool := []School{
        School{Establishment: "Lycée simon bolivar", StartDate: "Dec-2024", EndDate: "Jan-2022", Name: "Bac"},
        School{Establishment: "Lycée simon bolivar", StartDate: "Dec-2024", EndDate: "Jan-2022", Name: "Bac"},
    }
    pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, "")
    pdf.AddPage()
    pdf.AddUTF8Font("SUSE", "", "SUSE-Regular.ttf")
    pdf.AddUTF8Font("SUSE", "B", "SUSE-Bold.ttf")
    width, height := pdf.GetPageSize()
    rightSectionStartPos := width*LEFT_SECTION+LEFT_MARGIN
    rightElementWidth := width-rightSectionStartPos - LEFT_MARGIN
    left, top,_ ,_ := pdf.GetMargins()
	pdf.SetFillColor(183, 208, 229)
    pdf.Rect(0, 0, width*LEFT_SECTION, height, "F")
    pdf.RegisterImageOptions("test.jpeg", gofpdf.ImageOptions{})
    pdf.ClipCircle(left+((width*.25)/2), top+((width*.25)/2), (width*.25)/2, false)
    pdf.ImageOptions("test.jpeg", left, top, width*0.25, width*0.25, false, gofpdf.ImageOptions{}, 0, "")
    pdf.ClipEnd()
    pdf.SetY(width*.25+top)
    pdf.Ln(3)
    pdf.SetFont(FONT_NAME, "B", 18)
    pdf.Cell(width*.25, 6, "Coordonnées")
    pdf.Ln(9)
    pdf.MoveTo(pdf.GetXY())
    pdf.LineTo(width*LEFT_SECTION, pdf.GetY())
    pdf.ClosePath()
    pdf.SetLineWidth(1)
    pdf.SetFillColor(4,4,4)
    pdf.DrawPath("DF")
    pdf.Ln(10)
    printCurriculumInfo(pdf, width*.25, "Email", "test@test.com")
    printCurriculumInfo(pdf, width*.25, "Telephone", "0635392048")
    printCurriculumInfo(pdf, width*.25, "Adresse", "Paris, 75002")
    printCurriculumInfo(pdf, width*.25, "Age", "24 Ans")

    printCurriculumAbout(pdf, width, "Compétences", skill)
    printCurriculumAbout(pdf, width, "Centres d'intérêt", interest)

    pdf.SetXY(rightSectionStartPos, top)
    pdf.SetFont("SUSE", "B", 22)
    pdf.Cell(rightElementWidth ,6, "Ines Narayaninnaiken")
    pdf.Ln(10)
    pdf.SetX(rightSectionStartPos)
    pdf.SetFont("SUSE", "", 12)
    pdf.MultiCell(rightElementWidth, 6, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce ornare tempus nibh. Quisque congue urna nec efficitur hendrerit. Duis id auctor mauris.", "", "", false)
    pdf.Ln(10)
    pdf.SetX(rightSectionStartPos)
    pdf.SetFont("SUSE", "B", 22)
    pdf.Cell(rightElementWidth, 6, "Expériences Professionnelles")
    pdf.Ln(8)
    drawRightSectionHorizontalLine(pdf, rightSectionStartPos, width-LEFT_MARGIN)
    pdf.Ln(8)
    for _, v := range myWorks{
        pdf.SetX(rightSectionStartPos)
        pdf.SetFont("SUSE", "B", 13)
        pdf.MultiCell(rightElementWidth, 4, v.Position, "", "", false)
        pdf.SetFont("SUSE", "", 10)
        pdf.Ln(6)
        pdf.SetX(rightSectionStartPos)
        pdf.SetTextColor(123, 123, 123)
        pdf.Cellf(rightElementWidth, 4, "%v - %v", v.StartDate, v.EndDate)
        pdf.Ln(6)
        pdf.SetX(rightSectionStartPos)
        pdf.SetTextColor(0,0,0)
        pdf.SetFontSize(13)
        pdf.Cell(rightElementWidth, 6, v.Entreprise)
        pdf.Ln(8)
        for i, description := range v.Description{
            pdf.SetX(rightSectionStartPos)
            pdf.Circle(pdf.GetX()+3, pdf.GetY()+3, 1, "F")
            pdf.SetLineWidth(0.2)
            pdf.SetX(pdf.GetX()+6)
            pdf.SetFontSize(12)
            pdf.MultiCell(rightElementWidth-15, 6, description, "", "", false)
            if i == len(v.Description)-1{
                pdf.Ln(8)
            }
        }
    }

    pdf.SetX(rightSectionStartPos)
    pdf.SetFont("SUSE", "B", 22)
    pdf.Cell(rightElementWidth, 6, "Formation")
    pdf.Ln(8)
    drawRightSectionHorizontalLine(pdf, rightSectionStartPos, width-LEFT_MARGIN)
    pdf.Ln(8)
    for _, val := range mySchool{
        pdf.SetX(rightSectionStartPos)
        pdf.SetFont("SUSE", "B", 13)
        pdf.MultiCell(rightElementWidth, 4, val.Establishment, "", "", false)
        pdf.SetFont("SUSE", "", 10)
        pdf.Ln(6)
        pdf.SetX(rightSectionStartPos)
        pdf.SetTextColor(123, 123, 123)
        pdf.Cellf(rightElementWidth, 4, "%v - %v", val.StartDate, val.EndDate)
        pdf.Ln(6)
        pdf.SetX(rightSectionStartPos)
        pdf.SetTextColor(0,0,0)
        pdf.SetFontSize(13)
        pdf.Cell(rightElementWidth, 4, val.Name)
        pdf.Ln(8)
    }
    pdf.OutputFileAndClose("test.pdf")
}

func drawRightSectionHorizontalLine(pdf *gofpdf.Fpdf, startX float64, lineWidth float64){
    pdf.MoveTo(startX, pdf.GetY())
    pdf.LineTo(lineWidth, pdf.GetY())
    pdf.ClosePath()
    pdf.SetLineWidth(1)
    pdf.DrawPath("DF")
}

func printCurriculumInfo(pdf *gofpdf.Fpdf, width float64, label string, text string){
    pdf.SetFont(FONT_NAME, "B", 12)
    pdf.Cellf(width*.25, 4, "%v", label)
    pdf.Ln(6)
    pdf.SetFont(FONT_NAME, "", 12)
    pdf.Cellf(width*.25, 4, "%v", text)
    pdf.Ln(10)
}

func printCurriculumAbout(pdf *gofpdf.Fpdf, width float64, header string, data []string){
    pdf.SetFont(FONT_NAME, "B", 18)
    pdf.Ln(8)
    pdf.Cell(width*.25, 5, header)
    pdf.Ln(9)
    pdf.MoveTo(pdf.GetXY())
    pdf.LineTo(width*LEFT_SECTION, pdf.GetY())
    pdf.ClosePath()
    pdf.SetLineWidth(1)
    pdf.SetFillColor(4,4,4)
    pdf.DrawPath("DF")
    pdf.Ln(8)
    pdf.SetFont("SUSE", "", 12)
    for _, v := range data{
        pdf.Cell(width*.25, 4, v)
        pdf.Ln(6)
    }
    pdf.Ln(10)
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
