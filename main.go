package main

import (
	"fmt"
	"log"

	"github.com/ElioenaiFerrari/youdecide/app/command"
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/ElioenaiFerrari/youdecide/infra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("youdecide.db"))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&entity.Party{})
	db.AutoMigrate(&entity.Candidate{})
	db.AutoMigrate(&entity.Candidature{})

	validator := infra.NewValidator()
	findCandidateQuery := query.NewFindCandidateQuery(db)
	createCandidatureCommand := command.NewCreateCandidatureCommand(db, validator, findCandidateQuery)

	candidature, err := createCandidatureCommand.Exec(dto.CreateCandidatureDTO{
		CandidateCode: "13",
		Position:      "Presidency",
		Year:          2022,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(candidature)

}
