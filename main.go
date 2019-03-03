package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"hh-actors/models"
)

const (
	MaxTurn = 10
)

type Charactor struct {
	Name       string
	ImagePath  string
	VoiceActor string
}

var db *gorm.DB
var charcterImages []models.Image
var successImages []models.Image
var failureImages []models.Image

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = connectDb()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	images := models.FindImage(db)
	for _, image := range images {
		switch image.Type {
		case "character":
			charcterImages = append(charcterImages, image)
		case "success":
			successImages = append(successImages, image)
		case "failure":
			failureImages = append(failureImages, image)
		}
	}

	r := createRouter()

	setRoute(r)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	r.Run(":" + port)
}

func createRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "assets/css")
	r.Static("/image", "assets/image")
	return r
}

func setRoute(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		index(c)
	})

	router.POST("/choice", func(c *gin.Context) {
		finished := choice(c)
		if finished {
			c.Request.URL.Path = "/result"
		} else {
			c.Request.URL.Path = "/"
		}
		c.Request.Method = "GET"
		router.HandleContext(c)
	})

	router.GET("/result", func(c *gin.Context) {
		result(c)
	})
}

func index(c *gin.Context) {
	sessionId := getIntCookie(c, "sessionId")
	session := models.FindSession(db, sessionId)
	setCookie(c, "sessionId", fmt.Sprint(session.ID))

	indication, charactors := generateChactors(c, session)

	paths := strings.Split(session.ResultSummary, ",")
	iconPaths := make([]string, len(paths)-1)
	for i := range iconPaths {
		iconPaths[i] = "image/icon/" + paths[i]
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"indication": indication,
		"charactors": charactors,
		"iconPaths":  iconPaths,
	})
}

func generateChactors(c *gin.Context, session *models.UserSession) (string, []Charactor) {
	var images []models.Image
	for _, image := range charcterImages {
		if !strings.Contains(session.AnswerCharacterNames, image.Name) {
			images = append(images, image)
		}
	}

	choiced := randomChoice(images, 3)
	charactors := make([]Charactor, len(choiced))
	for i, v := range choiced {
		charactors[i] = Charactor{v.Name, "image/character/" + v.Path, v.VoiceActor}
	}

	if session.Turn == 0 || session.Turn >= MaxTurn {
		session.Turn = 1
		session.Hit = 0
		setCookie(c, "hit", "0")
	} else {
		session.Turn += 1
	}

	answer := charactors[rand.Intn(3)]
	session.Answer = answer.Name
	session.AnswerCharacterNames = session.AnswerCharacterNames + answer.Name + ","
	models.UpdateSession(db, session)

	var indication string
	if session.Turn >= 2 {
		indication = fmt.Sprintf("%d人目 %s (これまでの正解 %d人)", session.Turn, answer.VoiceActor, session.Hit)
	} else {
		indication = fmt.Sprintf("%d人目 %s", session.Turn, answer.VoiceActor)
	}

	return indication, charactors
}

func choice(c *gin.Context) bool {
	coiced := c.PostForm("choice")
	sessionId := getIntCookie(c, "sessionId")
	session := models.FindSession(db, sessionId)

	if session.Turn > MaxTurn || session.Turn == 0 {
		return false
	}

	if session.Answer == coiced {
		session.Hit += 1
		choiced := randomChoice(successImages, 1)
		session.ResultSummary = session.ResultSummary + choiced[0].Path + ","
	} else {
		choiced := randomChoice(failureImages, 1)
		session.ResultSummary = session.ResultSummary + choiced[0].Path + ","
	}

	models.UpdateSession(db, session)
	if session.Turn == MaxTurn {
		models.InsertResultLog(db, c.ClientIP(), session.Hit)
		return true
	} else {
		return false
	}
}

func randomChoice(images []models.Image, choice int) []models.Image {
	n := len(images)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		images[i], images[j] = images[j], images[i]
	}
	return images[:choice]
}

func result(c *gin.Context) {
	sessionId := getIntCookie(c, "sessionId")
	session := models.FindSession(db, sessionId)
	setCookie(c, "sessionId", fmt.Sprint(session.ID))

	var indication string
	if session.Hit == MaxTurn {
		indication = fmt.Sprintf("正解 %d人 PERFECT!!", MaxTurn)
	} else {
		indication = fmt.Sprintf("正解 %d人", session.Hit)
	}

	paths := strings.Split(session.ResultSummary, ",")
	iconPaths := make([]string, len(paths)-1)
	for i := range iconPaths {
		iconPaths[i] = "image/icon/" + paths[i]
	}

	hitAverage := fmt.Sprintf("直近の正解率： %d%%", models.HitAverage(db, c.ClientIP()))

	session.Turn = 0
	session.ResultSummary = ""
	session.AnswerCharacterNames = ""
	models.UpdateSession(db, session)

	c.HTML(http.StatusOK, "result.tmpl", gin.H{
		"indication": indication,
		"iconPaths":  iconPaths,
		"hitAverage": hitAverage,
	})
}

func setCookie(c *gin.Context, key string, value string) {
	c.SetCookie(key, value, 300, "/", "", false, true)
}

func getCookie(c *gin.Context, key string) string {
	value, err := c.Cookie(key)
	if err != nil {
		return ""
	}
	return value
}

func getIntCookie(c *gin.Context, key string) int {
	value := getCookie(c, key)
	if value == "" {
		return 0
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func connectDb() (*gorm.DB, error) {
	configs := map[string]string{}
	configs["host"] = os.Getenv("DB_HOST")
	configs["port"] = os.Getenv("DB_PORT")
	configs["user"] = os.Getenv("DB_USER")
	configs["password"] = os.Getenv("DB_PASSWORD")
	configs["dbname"] = os.Getenv("DB_SCHEMA")
	configs["sslmode"] = os.Getenv("DB_SSL")

	buf := []string{}
	for k, v := range configs {
		buf = append(buf, k+"="+v)
	}
	params := strings.Join(buf, " ")

	db, err := gorm.Open("postgres", params)
	if err != nil {
		log.Println("DB connect error!!")
		log.Println(err)
		return nil, err
	}

	log.Println("DB connect success!")

	return db, nil
}
