package main

import (
	"log"

	"github.com/ElioenaiFerrari/youdecide/app/command"
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	findPartyQuery := query.NewFindPartyQuery(db)
	findVoterQuery := query.NewFindVoterQuery(db)
	findLastVoteQuery := query.NewFindLastVoteQuery(db)
	createVoterCommand := command.NewCreateVoterCommand(db, validator)
	createCandidateCommand := command.NewCreateCandidateCommand(db, validator, findPartyQuery)
	createPartyCommand := command.NewCreatePartyCommand(db, validator)
	createCandidatureCommand := command.NewCreateCandidatureCommand(db, validator, findCandidateQuery)
	createBlockCommand := command.NewCreateBlockCommand(db, findLastBlockQuery)
	createVoteCommand := command.NewCreateVoteCommand(db, createBlockCommand, validator, findCandidatureQuery, findVoterQuery, findLastVoteQuery)

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

	_, mnemonic, err := createVoterCommand.Exec(dto.CreateVoterDTO{
		Email:    "elioenaferrari@gmail.com",
		Phone:    "27992150059",
		Name:     "Eli",
		Password: "123456",
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := createVoteCommand.Exec(dto.CreateVoteDTO{
		CandidatureCode: candidature.Code,
		Mnemonic:        *mnemonic,
		Password:        "123456",
	}); err != nil {
		log.Fatal(err)
	}

	if err := createVoteCommand.Exec(dto.CreateVoteDTO{
		CandidatureCode: candidature.Code,
		Mnemonic:        *mnemonic,
		Password:        "123456",
	}); err != nil {
		log.Fatal(err)
	}
}
