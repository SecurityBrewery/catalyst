package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func dashboardResponseID(doc *model.DashboardResponse) []driver.DocumentID {
	if doc == nil {
		return nil
	}
	return templateID(doc.ID)
}

func dashboardID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.DashboardCollectionName, id))}
}

func (s *Service) ListDashboards(ctx context.Context) ([]*model.DashboardResponse, error) {
	return s.database.DashboardList(ctx)
}

func (s *Service) CreateDashboard(ctx context.Context, dashboard *model.Dashboard) (doc *model.DashboardResponse, err error) {
	defer s.publishRequest(ctx, err, "CreateDashboard", dashboardResponseID(doc))
	return s.database.DashboardCreate(ctx, dashboard)
}

func (s *Service) GetDashboard(ctx context.Context, id string) (*model.DashboardResponse, error) {
	return s.database.DashboardGet(ctx, id)
}

func (s *Service) UpdateDashboard(ctx context.Context, id string, form *model.Dashboard) (doc *model.DashboardResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateDashboard", dashboardResponseID(doc))
	return s.database.DashboardUpdate(ctx, id, form)
}

func (s *Service) DeleteDashboard(ctx context.Context, id string) (err error) {
	defer s.publishRequest(ctx, err, "DeleteDashboard", dashboardID(id))
	return s.database.DashboardDelete(ctx, id)
}

func (s *Service) DashboardData(ctx context.Context, aggregation string, filter *string) (map[string]interface{}, error) {
	return s.database.WidgetData(ctx, aggregation, filter)
}
