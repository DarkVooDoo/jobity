package store

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/lib/pq"
)

const (
    FONT_NAME="SUSE"
    LEFT_SECTION=.35
    RIGHT_SECTION=.65
    LEFT_MARGIN=8
)

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

type Curriculum struct{
    Id string
    Name string
    Age string
    Description string
    Adresse string
    Gender string
    Email string
    Skill []string `json:"skill"`
    Interest []string `json:"interest"`
    Work []Work `json:"work"`
    School []School `json:"school"`
    UserId string
    JobId string
}

func (c *Curriculum) Get(userId string)error{
    var school, work, description sql.NullString
    conn, err := GetDBConn()
    if err != nil{
        log.Println("err db conn")
        return errors.New("error conn db")
    }
    userId = DecryptCurriculumId(userId)
    curriculumRow := conn.QueryRowContext(context.Background(), `SELECT CONCAT(u.firstname,' ',u.lastname), CONCAT(u.city,', ', u.postal), CONCAT(EXTRACT(YEAR FROM AGE(NOW(), u.birthdate)), ' Ans'), u.interest, u.skill, u.school, u.work, u.email, u.description FROM Users AS u WHERE u.id=$1`, userId)
    if err := curriculumRow.Scan(&c.Name, &c.Adresse, &c.Age, pq.Array(&c.Interest), pq.Array(&c.Skill), &school, &work, &c.Email, &description); err != nil{
        log.Println(err)
        return errors.New("error query")
    }
    c.Description = description.String
    if err := json.Unmarshal([]byte(school.String), &c.School); err != nil{
        log.Println(err)
    }
    if err := json.Unmarshal([]byte(work.String), &c.Work); err != nil{
        log.Println(err)
    }
    return nil
}

func (c *Curriculum) Save(userId string)error{
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error db conn")
    }
    defer conn.Close()
    work, _ := json.Marshal(c.Work)
    school, _ := json.Marshal(c.School)
    _, err = conn.ExecContext(context.Background(), `UPDATE Users SET interest=$1, skill=$2, school=$3, work=$4 WHERE id=$5`, pq.Array(c.Interest), pq.Array(c.Skill), string(school), string(work), userId)
    if err != nil{
        log.Println(err)    
        return errors.New("error in the query")
    }
    return nil
}

func CreateCurriculumPDF(curriculum Curriculum) *gofpdf.Fpdf{
    pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, "")
    pdf.AddPage()
    pdf.SetFontLocation("static")
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
    printCurriculumInfo(pdf, width*.25, "Email", curriculum.Email)
    printCurriculumInfo(pdf, width*.25, "Telephone", "0635392048")
    printCurriculumInfo(pdf, width*.25, "Adresse", curriculum.Adresse)
    printCurriculumInfo(pdf, width*.25, "Age", curriculum.Age)

    printCurriculumAbout(pdf, width, "Compétences", curriculum.Skill)
    printCurriculumAbout(pdf, width, "Centres d'intérêt", curriculum.Interest)

    pdf.SetXY(rightSectionStartPos, top)
    pdf.SetFont("SUSE", "B", 22)
    pdf.Cell(rightElementWidth ,6, curriculum.Name)
    pdf.Ln(10)
    pdf.SetX(rightSectionStartPos)
    pdf.SetFont("SUSE", "", 12)
    pdf.MultiCell(rightElementWidth, 6, curriculum.Description, "", "", false)
    pdf.Ln(10)
    pdf.SetX(rightSectionStartPos)
    pdf.SetFont("SUSE", "B", 22)
    pdf.Cell(rightElementWidth, 6, "Expériences Professionnelles")
    pdf.Ln(8)
    drawRightSectionHorizontalLine(pdf, rightSectionStartPos, width-LEFT_MARGIN)
    pdf.Ln(8)
    for _, v := range curriculum.Work{
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
    for _, val := range curriculum.School{
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
    return pdf
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

func GetJobCurriculum(jobId string)(candidates []Curriculum, interview []Curriculum){
    var id, name, userId, status string
    var adresse, workString, schoolString, gender sql.NullString 
    var age sql.NullInt16
    var skill, interest []string
    var school []School
    var work []Work
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return 
    }
    defer conn.Close()
    curriculumRows, err := conn.QueryContext(context.Background(), `SELECT c.id, u.id, CONCAT(u.firstname, ' ', u.lastname), CONCAT(u.City, ', ', u.Postal), DATE_PART('year', AGE(NOW(), u.birthdate)), u.gender, u.skill, u.interest, u.school,u.work, c.status FROM JobApplication AS c LEFT JOIN Users AS u ON u.id=c.user_id WHERE c.job_id=$1 AND c.status < 'Reject'`, jobId)
    if err != nil{
        log.Println(err)
        return 
    }
    for curriculumRows.Next(){
        curriculumRows.Scan(&id, &userId, &name, &adresse, &age, &gender, pq.Array(&skill), pq.Array(&interest), &schoolString, &workString, &status)
        json.Unmarshal([]byte(string(schoolString.String)), &school)
        json.Unmarshal([]byte(string(workString.String)), &work)
        userId = EncryptCurriculumId(userId)
        if status == "Interview"{
            interview = append(interview, Curriculum{Name: name, Adresse: adresse.String, Skill: skill, Interest: interest, Id: id, Work: work, School: school, Age: fmt.Sprintf("%v Ans", age.Int16), Gender: gender.String, UserId: userId, JobId: jobId})
        }else{
            candidates = append(candidates, Curriculum{Name: name, Adresse: adresse.String, Skill: skill, Interest: interest, Id: id, Work: work, School: school, Age: fmt.Sprintf("%v Ans", age.Int16), Gender: gender.String, UserId: userId, JobId: jobId})
        }
    }
    return 
}

func GetInterviewType()[]string{
    var interviewType []string
    var name string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return interviewType
    }
    defer conn.Close()
    row, err := conn.QueryContext(context.Background(), `SELECT unnest(enum_range(NULL::interview_type))`)
    if err != nil{
        log.Println(err)
        return interviewType
    }
    for row.Next(){
        if err := row.Scan(&name); err != nil{
            log.Println(err)
            return interviewType
        }
        interviewType = append(interviewType, name)
    }
    return interviewType
}

func EncryptCurriculumId(id string)string{
    block, err := aes.NewCipher([]byte(os.Getenv("CURRICULUM_INDEX_KEY")))
    if err != nil{
        log.Printf("error creating cypher: %v", err)
    }
    gmc, err := cipher.NewGCM(block)
    if err != nil{
        log.Printf("error gcm: %v", err)
    }
    nonce := make([]byte, gmc.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil{
        log.Printf("error read full %v", err)
    }
    cipherText := gmc.Seal(nonce, nonce, []byte(id), nil)
    enc := hex.EncodeToString(cipherText)
    return enc
}

func DecryptCurriculumId(encString string)string{
    block, err := aes.NewCipher([]byte(os.Getenv("CURRICULUM_INDEX_KEY")))
    if err != nil{
        log.Printf("error creating cypher: %v", err)
    }
    gmc, err := cipher.NewGCM(block)
    if err != nil{
        log.Printf("error gcm: %v", err)

    }
    decodedCipherText, err := hex.DecodeString(encString)
    if err != nil{
        log.Printf("error decoding %v", err)
    }
    decryptData, err := gmc.Open(nil, decodedCipherText[:gmc.NonceSize()], decodedCipherText[gmc.NonceSize():], nil)
    if err != nil{
        log.Printf("error opening %v", err)
    }

    return string(decryptData)
}
