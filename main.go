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

	pt, _ := createPartyCommand.Exec(dto.CreatePartyDTO{
		Initials: "PT",
		Name:     "Partido dos Trabalhadores",
	})

	psl, _ := createPartyCommand.Exec(dto.CreatePartyDTO{
		Initials: "PSL",
		Name:     "Partido Social Liberal",
	})

	lula, _ := createCandidateCommand.Exec(dto.CreateCandidateDTO{
		Name:          "Lula",
		PartyInitials: pt.Initials,
		ImageURL:      "https://s2.glbimg.com/4whWaUxUmO9OZKbeUXLvPK9nnZ0=/1200x/smart/filters:cover():strip_icc()/i.s3.glbimg.com/v1/AUTH_59edd422c0c84a879bd37670ae4f538a/internal_photos/bs/2022/N/A/x4O7j9R5i4N0ecWBU7vw/lula-31.jpg",
	})

	bolsonaro, _ := createCandidateCommand.Exec(dto.CreateCandidateDTO{
		Name:          "Bolsonaro",
		PartyInitials: psl.Initials,
		ImageURL:      "https://uploads.metropoles.com/wp-content/uploads/2022/11/01174231/Apo%CC%81s-reunia%CC%83o-com-ministros-e-aliados-jair-Bolsonaro-fara%CC%81-primeiro-pronunciamento-aos-brasileiros-9.jpeg",
	})

	candidatureLula, _ := createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		Code:          "13",
		CandidateName: lula.Name,
		Position:      "presidência",
		Year:          2022,
	})

	_, _ = createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		Code:          "22",
		CandidateName: bolsonaro.Name,
		Position:      "presidência",
		Year:          2022,
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
		CandidatureCode: candidatureLula.Code,
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

	app.GET("/", func(c echo.Context) error {
		candidatures, err := listCandidaturesQuery.Exec("year = ?", time.Now().Year())

		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "index", candidatures)
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
