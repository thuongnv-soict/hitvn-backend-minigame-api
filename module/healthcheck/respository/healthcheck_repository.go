package respository

import (
	"context"
	"database/sql"
	"g-tech.com/infrastructure/repository"
	"time"
)

type IHealthCheckRepository interface {
	GetSQLStatus(ctx context.Context) error
}

type HealthCheckRepository struct {
	MySql 			repository.MySqlRepository
	Timeout    		time.Duration
}

/**
 * Return a new HealthCheckRepository
 */
func NewHealthCheckRepository(dbContext *sql.DB, timeout time.Duration) IHealthCheckRepository {
	repository := HealthCheckRepository{}
	repository.MySql.SetDbContext(dbContext)
	repository.Timeout = timeout
	return &repository
}

/**
 * Returns Couchbase status
 */
func (repository *HealthCheckRepository) GetSQLStatus(ctx context.Context) error {
	// 	Setting up timeout
	ctx, cancel := context.WithTimeout(ctx, repository.Timeout)
	defer cancel()

	err := repository.MySql.DbContext.Ping()
	if err != nil {
		return err
	}

	return nil
}



