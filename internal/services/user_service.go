package services

import (
	"cmd/internal/db"
	"cmd/internal/leetcode_api"
	"cmd/internal/storage/dbs"
	v1 "cmd/internal/templates/v1"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	repository *db.Repository
}

func NewUserService(repository *db.Repository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) Upsert(ctx context.Context, userSlug string) int64 {
	id, err := s.repository.Queries().UserGetBySocialProviderId(ctx, userSlug)

	if err == sql.ErrNoRows {
		fmt.Println("no rows in result set________________")
		//creating new user
		now := time.Now().UTC()

		matchedUser := leetcode_api.MatchedUserMapToUserProfile(userSlug)
		userData := v1.NewLcUserDataFromReq(*matchedUser)
		fmt.Println("userData", userData)

		id, err := s.repository.Queries().UserNewAndParse(ctx, dbs.UserNewAndParseParams{
			Username:             userData.Username,
			SocialProviderUserID: userData.UserSlug})
		fmt.Println("----------created id: ", id)
		if err != nil {
			fmt.Println("39str ", err)

			return 0
		}
		err = s.repository.Queries().InsertStatsInfo(ctx, dbs.InsertStatsInfoParams{
			UserID: id,
			Rank:   int64(userData.Rank),
			EasySubmits: sql.NullInt64{
				Int64: userData.EasyCount,
				Valid: true},
			MediumSubmits: sql.NullInt64{
				Int64: userData.MediumCount,
				Valid: true},
			HardSubmits: sql.NullInt64{
				Int64: userData.HardCount,
				Valid: true},
			TotalSubmits: sql.NullInt64{
				Int64: userData.TotalCount,
				Valid: true},
			CreatedAt: now,
			UpdatedAt: now})
		if err != nil {
			log.Printf("Database err %s\n", err)
		}
		return id
	} else if err != nil {
		fmt.Println("56str ", err)
		log.Printf("Database err %s\n", err)
	}
	return id
}

func (s *UserService) GetByStatsById(ctx context.Context, userId int64) (*v1.LcUserData, error) {
	userStatsByIDRow, err := s.repository.Queries().UserGetStatsByID(ctx, userId)
	if err == sql.ErrNoRows {
		fmt.Println("No user stats row: ", err)
		return nil, err
	} else if err != nil {
		fmt.Println("Database err ", err)
		return nil, err
	}

	userProfileData := &v1.LcUserData{
		Username:    userStatsByIDRow.Username,
		UserSlug:    userStatsByIDRow.Username,
		Rank:        float64(userStatsByIDRow.Rank),
		EasyCount:   userStatsByIDRow.EasySubmits.Int64,
		MediumCount: userStatsByIDRow.MediumSubmits.Int64,
		HardCount:   userStatsByIDRow.HardSubmits.Int64,
		TotalCount:  userStatsByIDRow.TotalSubmits.Int64,
	}

	return userProfileData, nil
}
