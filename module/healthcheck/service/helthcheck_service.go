package service

import (
	"context"
	"g-tech.com/module/healthcheck/respository"
)

type IHealthCheckService interface {
	GetSQLStatus(ctx context.Context) error
}

type HealthCheckService struct {
	Repository respository.IHealthCheckRepository
}

/**
 * Returns a new HealthCheckService
 */
func NewHealthCheckService(respository respository.IHealthCheckRepository) IHealthCheckService {
	return &HealthCheckService{
		Repository: respository,
	}
}

/**
 * Return Couchbase status
 */
func (service *HealthCheckService) GetSQLStatus(ctx context.Context) error {

	return service.Repository.GetSQLStatus(ctx)
}
