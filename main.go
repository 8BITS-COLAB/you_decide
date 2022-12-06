package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ElioenaiFerrari/youdecide/app/command"
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func take(items map[string]cache.Item, size int) []dto.CreateVoteDTO {
	var votes []dto.CreateVoteDTO

	for _, v := range items {
		vote := v.Object.(dto.CreateVoteDTO)
		votes = append(votes, vote)
		if len(votes) == size {
			return votes
		}
	}

	return nil
}

func worker(ch chan int, rwm *sync.RWMutex, votesPool *cache.Cache, createVoteCommand *command.CreateVoteCommand, size int) {
	for range time.Tick(time.Second) {
		poolSize := votesPool.ItemCount()
		ch <- poolSize

		if poolSize >= size {
			items := votesPool.Items()
			votes := take(items, size)

			rwm.Lock()
			createVoteCommand.Exec(votes)
			rwm.Unlock()

			for k := range items {
				votesPool.Delete(k)
			}

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

	votesPool := cache.New(15*time.Minute, 30*time.Minute)

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
	findPartyQuery := query.NewFindPartyQuery(db)
	findVoterQuery := query.NewFindVoterQuery(db)
	findLastVoteQuery := query.NewFindLastVoteQuery(db)
	createVoterCommand := command.NewCreateVoterCommand(db, validator)
	createCandidateCommand := command.NewCreateCandidateCommand(db, validator, findPartyQuery)
	createPartyCommand := command.NewCreatePartyCommand(db, validator)
	createCandidatureCommand := command.NewCreateCandidatureCommand(db, validator, findCandidateQuery)
	createBlockCommand := command.NewCreateBlockCommand(db, findLastBlockQuery)
	createVoteCommand := command.NewCreateVoteCommand(db, createBlockCommand, validator, findCandidatureQuery, findVoterQuery, findLastVoteQuery)
	createVotePoolCommand := command.NewCreateVotePoolCommand(validator, votesPool)

	ch := make(chan int)

	go worker(ch, new(sync.RWMutex), votesPool, createVoteCommand, 2)

	party, err := createPartyCommand.Exec(dto.CreatePartyDTO{
		Initials: "PT",
		Name:     "Partido dos Trabalhadores",
	})

	if err != nil {
		log.Fatal(err)
	}

	candidate, err := createCandidateCommand.Exec(dto.CreateCandidateDTO{
		Code:          "13",
		Name:          "LULA",
		PartyInitials: party.Initials,
	})

	if err != nil {
		log.Fatal(err)
	}

	candidature, err := createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		Code:          "13",
		CandidateName: candidate.Name,
		Position:      "Presidency",
		Year:          2022,
	})

	if err != nil {
		log.Fatal(err)
	}

	voter1, mnemonic1, err := createVoterCommand.Exec(dto.CreateVoterDTO{
		Email:    "elioenaferrari@gmail.com",
		Phone:    "27992150059",
		Name:     "Eli",
		Password: "123456",
	})

	if err != nil {
		log.Fatal(err)
	}

	voter2, mnemonic2, err := createVoterCommand.Exec(dto.CreateVoterDTO{
		Email:    "elioenaferrari2@gmail.com",
		Phone:    "27992150058",
		Name:     "Eli",
		Password: "123456",
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := createVotePoolCommand.Exec(dto.CreateVoteDTO{
		Email:           voter1.Email,
		CandidatureCode: candidature.Code,
		Mnemonic:        *mnemonic1,
		Password:        "123456",
	}); err != nil {
		log.Fatal(err)
	}

	if err := createVotePoolCommand.Exec(dto.CreateVoteDTO{
		Email:           voter2.Email,
		CandidatureCode: candidature.Code,
		Mnemonic:        *mnemonic2,
		Password:        "123456",
	}); err != nil {
		log.Fatal(err)
	}

	for msg := range ch {
		fmt.Printf("pool size in %d\n", msg)
	}
}
