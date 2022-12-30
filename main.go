package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/ElioenaiFerrari/youdecide/app/command"
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/jaswdr/faker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var votesPool []dto.CreateVoteDTO

func worker(ch chan int, rwm *sync.RWMutex, createVoteCommand *command.CreateVoteCommand, size int) {
	for range time.Tick(time.Second) {
		poolSize := len(votesPool)
		ch <- poolSize

		if poolSize >= size {
			rwm.Lock()
			votes := votesPool[:size-1]

			if err := createVoteCommand.Exec(votes); err != nil {
				log.Fatal(err)
			} else {
				votesPool = votesPool[size:]
			}

			rwm.Unlock()
		}
	}

}

func main() {
	db, err := gorm.Open(sqlite.Open("youdecide.db"), &gorm.Config{
		FullSaveAssociations: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&entity.Party{})
	db.AutoMigrate(&entity.Candidate{})
	db.AutoMigrate(&entity.Candidature{})
	db.AutoMigrate(&entity.Voter{})
	db.AutoMigrate(&entity.Vote{})
	db.AutoMigrate(&entity.Block{})

	validator := validator.New()
	findLastBlockQuery := query.NewFindLastBlockQuery(db)
	findCandidatureQuery := query.NewFindCandidatureQuery(db)
	findCandidateQuery := query.NewFindCandidateQuery(db)
	listCandidaturesQuery := query.NewListCandidaturesQuery(db)
	findPartyQuery := query.NewFindPartyQuery(db)
	findVoterQuery := query.NewFindVoterQuery(db)
	findLastVoteQuery := query.NewFindLastVoteQuery(db)
	createVoterCommand := command.NewCreateVoterCommand(db, validator)
	createCandidateCommand := command.NewCreateCandidateCommand(db, validator, findPartyQuery)
	createPartyCommand := command.NewCreatePartyCommand(db, validator)
	createCandidatureCommand := command.NewCreateCandidatureCommand(db, validator, findCandidateQuery)
	createBlockCommand := command.NewCreateBlockCommand(db, findLastBlockQuery)
	createGenesisBlockCommand := command.NewCreateGenesisBlockCommand(db)
	createVoteCommand := command.NewCreateVoteCommand(db, createBlockCommand, validator, findCandidatureQuery, findVoterQuery, findLastVoteQuery)
	createVotePoolCommand := command.NewCreateVotePoolCommand(validator, &votesPool)

	ch := make(chan int)

	go func() {
		for msg := range ch {
			fmt.Printf("pool size in %d\n", msg)
		}
	}()
	go worker(ch, new(sync.RWMutex), createVoteCommand, 2)

	fake := faker.New()

	party1, _ := createPartyCommand.Exec(dto.CreatePartyDTO{
		Initials: fake.Lorem().Word(),
		Name:     fake.Lorem().Sentence(3),
	})

	party2, _ := createPartyCommand.Exec(dto.CreatePartyDTO{
		Initials: fake.Lorem().Word(),
		Name:     fake.Lorem().Sentence(3),
	})

	candidate1, _ := createCandidateCommand.Exec(dto.CreateCandidateDTO{
		Name:          fake.Person().FirstName(),
		PartyInitials: party1.Initials,
		ImageURL:      "https://images.pexels.com/photos/1222271/pexels-photo-1222271.jpeg?auto=compress&cs=tinysrgb&w=800",
	})

	candidate2, _ := createCandidateCommand.Exec(dto.CreateCandidateDTO{
		Name:          fake.Person().FirstName(),
		PartyInitials: party2.Initials,
		ImageURL:      "https://images.pexels.com/photos/1239291/pexels-photo-1239291.jpeg?auto=compress&cs=tinysrgb&w=800",
	})

	candidature1, _ := createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		Code:          fake.UUID().V4(),
		CandidateName: candidate1.Name,
		Position:      "presidência",
		Year:          time.Now().Year(),
	})

	_, _ = createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		Code:          fake.UUID().V4(),
		CandidateName: candidate2.Name,
		Position:      "presidência",
		Year:          time.Now().Year(),
	})

	voter1, mnemonic1, err := createVoterCommand.Exec(dto.CreateVoterDTO{
		Email:    "elioenaferrari@gmail.com",
		Phone:    "27992150059",
		Name:     "Eli",
		Password: "123456",
	})

	if err != nil {
		log.Fatal(err)
	}

	_, mnemonic2, err := createVoterCommand.Exec(dto.CreateVoterDTO{
		Email:    "elioenaferrari2@gmail.com",
		Phone:    "27992150058",
		Name:     "Eli",
		Password: "123456",
	})

	log.Println("MNEMONIC", *mnemonic2)

	if err != nil {
		log.Fatal(err)
	}

	block, err := createGenesisBlockCommand.Exec()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GENESIS BLOCK CREATED %d", block.Index)

	if err := createVotePoolCommand.Exec(dto.CreateVoteDTO{
		Email:           voter1.Email,
		CandidatureCode: candidature1.Code,
		Mnemonic:        *mnemonic1,
		Password:        "123456",
	}); err != nil {
		log.Fatal(err)
	}

	app := echo.New()
	v1 := app.Group("/api/v1")
	websocket := melody.New()

	app.Use(middleware.Logger())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	app.Renderer = t

	app.GET("/candidatures", func(c echo.Context) error {
		candidatures, err := listCandidaturesQuery.Exec("year = ?", time.Now().Year())

		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "index", candidatures)
	})

	app.GET("/candidatures/:code", func(c echo.Context) error {
		candidature, err := findCandidatureQuery.Exec("code = ? AND year = ?", c.Param("code"), time.Now().Year())

		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "candidature", candidature)
	})

	websocket.HandleMessage(func(s *melody.Session, b []byte) {
		var createVoteDTO dto.CreateVoteDTO

		if err := sonic.Unmarshal(b, &createVoteDTO); err != nil {
			s.Write([]byte(err.Error()))
			return
		}

		if err := createVotePoolCommand.Exec(createVoteDTO); err != nil {
			s.Write([]byte(err.Error()))
			s.Write([]byte(err.Error()))
			return
		}

		s.Write([]byte("VOTE REGISTERED!"))
	})

	v1.GET("/ws", func(c echo.Context) error {
		return websocket.HandleRequest(c.Response().Writer, c.Request())
	})

	app.Logger.Fatal(app.Start(":4000"))
}
