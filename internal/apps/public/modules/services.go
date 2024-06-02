package modules

import "github.com/tamboto2000/99-backend-exercise/internal/apps/public/services/aggregator"

type Services struct {
	AggregatorService aggregator.AggregatorService
}

func NewServices(exSvcs ExternalServices) Services {
	// aggregator service
	aggrSvc := aggregator.NewAggregatorService(exSvcs.UserService, exSvcs.ListingService)

	return Services{
		AggregatorService: aggrSvc,
	}
}
