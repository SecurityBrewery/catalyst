package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	// "github.com/xeipuuv/gojsonschema"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

type HTTPError struct {
	Status   int
	Internal error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError(%d): %s", e.Status, e.Internal)
}

type Service interface {
	ListAutomations(context.Context) ([]*model.AutomationResponse, error)
	CreateAutomation(context.Context, *model.AutomationForm) (*model.AutomationResponse, error)
	GetAutomation(context.Context, string) (*model.AutomationResponse, error)
	UpdateAutomation(context.Context, string, *model.AutomationForm) (*model.AutomationResponse, error)
	DeleteAutomation(context.Context, string) error
	CurrentUser(context.Context) (*model.UserResponse, error)
	CurrentUserData(context.Context) (*model.UserDataResponse, error)
	UpdateCurrentUserData(context.Context, *model.UserData) (*model.UserDataResponse, error)
	ListJobs(context.Context) ([]*model.JobResponse, error)
	RunJob(context.Context, *model.JobForm) error
	GetJob(context.Context, string) (*model.JobResponse, error)
	UpdateJob(context.Context, string, *model.Job) (*model.JobResponse, error)
	GetLogs(context.Context, string) ([]*model.LogEntry, error)
	ListPlaybooks(context.Context) ([]*model.PlaybookTemplateResponse, error)
	CreatePlaybook(context.Context, *model.PlaybookTemplateForm) (*model.PlaybookTemplateResponse, error)
	GetPlaybook(context.Context, string) (*model.PlaybookTemplateResponse, error)
	UpdatePlaybook(context.Context, string, *model.PlaybookTemplateForm) (*model.PlaybookTemplateResponse, error)
	DeletePlaybook(context.Context, string) error
	GetSettings(context.Context) (*model.Settings, error)
	GetStatistics(context.Context) (*model.Statistics, error)
	ListTasks(context.Context) ([]*model.TaskWithContext, error)
	ListTemplates(context.Context) ([]*model.TicketTemplateResponse, error)
	CreateTemplate(context.Context, *model.TicketTemplateForm) (*model.TicketTemplateResponse, error)
	GetTemplate(context.Context, string) (*model.TicketTemplateResponse, error)
	UpdateTemplate(context.Context, string, *model.TicketTemplateForm) (*model.TicketTemplateResponse, error)
	DeleteTemplate(context.Context, string) error
	ListTickets(context.Context, *string, *int, *int, []string, []bool, *string) (*model.TicketList, error)
	CreateTicket(context.Context, *model.TicketForm) (*model.TicketResponse, error)
	CreateTicketBatch(context.Context, []*model.TicketForm) error
	GetTicket(context.Context, int64) (*model.TicketWithTickets, error)
	UpdateTicket(context.Context, int64, *model.Ticket) (*model.TicketWithTickets, error)
	DeleteTicket(context.Context, int64) error
	AddArtifact(context.Context, int64, *model.Artifact) (*model.TicketWithTickets, error)
	GetArtifact(context.Context, int64, string) (*model.Artifact, error)
	SetArtifact(context.Context, int64, string, *model.Artifact) (*model.TicketWithTickets, error)
	RemoveArtifact(context.Context, int64, string) (*model.TicketWithTickets, error)
	EnrichArtifact(context.Context, int64, string, *model.EnrichmentForm) (*model.TicketWithTickets, error)
	RunArtifact(context.Context, int64, string, string) error
	AddComment(context.Context, int64, *model.CommentForm) (*model.TicketWithTickets, error)
	RemoveComment(context.Context, int64, int) (*model.TicketWithTickets, error)
	LinkFiles(context.Context, int64, []*model.File) (*model.TicketWithTickets, error)
	AddTicketPlaybook(context.Context, int64, *model.PlaybookTemplateForm) (*model.TicketWithTickets, error)
	RemoveTicketPlaybook(context.Context, int64, string) (*model.TicketWithTickets, error)
	SetTask(context.Context, int64, string, string, *model.Task) (*model.TicketWithTickets, error)
	CompleteTask(context.Context, int64, string, string, map[string]interface{}) (*model.TicketWithTickets, error)
	RunTask(context.Context, int64, string, string) error
	SetReferences(context.Context, int64, []*model.Reference) (*model.TicketWithTickets, error)
	SetSchema(context.Context, int64, string) (*model.TicketWithTickets, error)
	LinkTicket(context.Context, int64, int64) (*model.TicketWithTickets, error)
	UnlinkTicket(context.Context, int64, int64) (*model.TicketWithTickets, error)
	ListTicketTypes(context.Context) ([]*model.TicketTypeResponse, error)
	CreateTicketType(context.Context, *model.TicketTypeForm) (*model.TicketTypeResponse, error)
	GetTicketType(context.Context, string) (*model.TicketTypeResponse, error)
	UpdateTicketType(context.Context, string, *model.TicketTypeForm) (*model.TicketTypeResponse, error)
	DeleteTicketType(context.Context, string) error
	ListUserData(context.Context) ([]*model.UserDataResponse, error)
	GetUserData(context.Context, string) (*model.UserDataResponse, error)
	UpdateUserData(context.Context, string, *model.UserData) (*model.UserDataResponse, error)
	ListUsers(context.Context) ([]*model.UserResponse, error)
	CreateUser(context.Context, *model.UserForm) (*model.NewUserResponse, error)
	GetUser(context.Context, string) (*model.UserResponse, error)
	UpdateUser(context.Context, string, *model.UserForm) (*model.UserResponse, error)
	DeleteUser(context.Context, string) error
}

func NewServer(service Service, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares...)

	s := &server{service}

	r.Get("/automations", s.listAutomationsHandler)
	r.Put("/automations", s.createAutomationHandler)
	r.Get("/automations/{id}", s.getAutomationHandler)
	r.Put("/automations/{id}", s.updateAutomationHandler)
	r.Delete("/automations/{id}", s.deleteAutomationHandler)
	r.Get("/currentuser", s.currentUserHandler)
	r.Get("/currentuserdata", s.currentUserDataHandler)
	r.Put("/currentuserdata", s.updateCurrentUserDataHandler)
	r.Get("/jobs", s.listJobsHandler)
	r.Put("/jobs", s.runJobHandler)
	r.Get("/jobs/{id}", s.getJobHandler)
	r.Put("/jobs/{id}", s.updateJobHandler)
	r.Get("/logs/{reference}", s.getLogsHandler)
	r.Get("/playbooks", s.listPlaybooksHandler)
	r.Put("/playbooks", s.createPlaybookHandler)
	r.Get("/playbooks/{id}", s.getPlaybookHandler)
	r.Put("/playbooks/{id}", s.updatePlaybookHandler)
	r.Delete("/playbooks/{id}", s.deletePlaybookHandler)
	r.Get("/settings", s.getSettingsHandler)
	r.Get("/statistics", s.getStatisticsHandler)
	r.Get("/tasks", s.listTasksHandler)
	r.Get("/templates", s.listTemplatesHandler)
	r.Put("/templates", s.createTemplateHandler)
	r.Get("/templates/{id}", s.getTemplateHandler)
	r.Put("/templates/{id}", s.updateTemplateHandler)
	r.Delete("/templates/{id}", s.deleteTemplateHandler)
	r.Get("/tickets", s.listTicketsHandler)
	r.Put("/tickets", s.createTicketHandler)
	r.Put("/tickets/batch", s.createTicketBatchHandler)
	r.Get("/tickets/{id}", s.getTicketHandler)
	r.Put("/tickets/{id}", s.updateTicketHandler)
	r.Delete("/tickets/{id}", s.deleteTicketHandler)
	r.Put("/tickets/{id}/artifacts", s.addArtifactHandler)
	r.Get("/tickets/{id}/artifacts/{name}", s.getArtifactHandler)
	r.Put("/tickets/{id}/artifacts/{name}", s.setArtifactHandler)
	r.Delete("/tickets/{id}/artifacts/{name}", s.removeArtifactHandler)
	r.Put("/tickets/{id}/artifacts/{name}/enrich", s.enrichArtifactHandler)
	r.Put("/tickets/{id}/artifacts/{name}/run/{automation}", s.runArtifactHandler)
	r.Put("/tickets/{id}/comments", s.addCommentHandler)
	r.Delete("/tickets/{id}/comments/{commentID}", s.removeCommentHandler)
	r.Put("/tickets/{id}/files", s.linkFilesHandler)
	r.Put("/tickets/{id}/playbooks", s.addTicketPlaybookHandler)
	r.Delete("/tickets/{id}/playbooks/{playbookID}", s.removeTicketPlaybookHandler)
	r.Put("/tickets/{id}/playbooks/{playbookID}/task/{taskID}", s.setTaskHandler)
	r.Put("/tickets/{id}/playbooks/{playbookID}/task/{taskID}/complete", s.completeTaskHandler)
	r.Put("/tickets/{id}/playbooks/{playbookID}/task/{taskID}/run", s.runTaskHandler)
	r.Put("/tickets/{id}/references", s.setReferencesHandler)
	r.Put("/tickets/{id}/schema", s.setSchemaHandler)
	r.Patch("/tickets/{id}/tickets", s.linkTicketHandler)
	r.Delete("/tickets/{id}/tickets", s.unlinkTicketHandler)
	r.Get("/tickettypes", s.listTicketTypesHandler)
	r.Put("/tickettypes", s.createTicketTypeHandler)
	r.Get("/tickettypes/{id}", s.getTicketTypeHandler)
	r.Put("/tickettypes/{id}", s.updateTicketTypeHandler)
	r.Delete("/tickettypes/{id}", s.deleteTicketTypeHandler)
	r.Get("/userdata", s.listUserDataHandler)
	r.Get("/userdata/{id}", s.getUserDataHandler)
	r.Put("/userdata/{id}", s.updateUserDataHandler)
	r.Get("/users", s.listUsersHandler)
	r.Put("/users", s.createUserHandler)
	r.Get("/users/{id}", s.getUserHandler)
	r.Put("/users/{id}", s.updateUserHandler)
	r.Delete("/users/{id}", s.deleteUserHandler)
	return r
}

type server struct {
	service Service
}

func (s *server) listAutomationsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListAutomations(r.Context())
	response(w, result, err)
}

func (s *server) createAutomationHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.AutomationFormSchema.Validate(jl)

	var automationP *model.AutomationForm
	if err := parseBody(r, &automationP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreateAutomation(r.Context(), automationP)
	response(w, result, err)
}

func (s *server) getAutomationHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetAutomation(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateAutomationHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.AutomationFormSchema.Validate(jl)

	var automationP *model.AutomationForm
	if err := parseBody(r, &automationP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateAutomation(r.Context(), idP, automationP)
	response(w, result, err)
}

func (s *server) deleteAutomationHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	response(w, nil, s.service.DeleteAutomation(r.Context(), idP))
}

func (s *server) currentUserHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.CurrentUser(r.Context())
	response(w, result, err)
}

func (s *server) currentUserDataHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.CurrentUserData(r.Context())
	response(w, result, err)
}

func (s *server) updateCurrentUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.UserDataSchema.Validate(jl)

	var userdataP *model.UserData
	if err := parseBody(r, &userdataP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateCurrentUserData(r.Context(), userdataP)
	response(w, result, err)
}

func (s *server) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListJobs(r.Context())
	response(w, result, err)
}

func (s *server) runJobHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.JobFormSchema.Validate(jl)

	var jobP *model.JobForm
	if err := parseBody(r, &jobP); err != nil {
		JSONError(w, err)
		return
	}

	response(w, nil, s.service.RunJob(r.Context(), jobP))
}

func (s *server) getJobHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetJob(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateJobHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.JobSchema.Validate(jl)

	var jobP *model.Job
	if err := parseBody(r, &jobP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateJob(r.Context(), idP, jobP)
	response(w, result, err)
}

func (s *server) getLogsHandler(w http.ResponseWriter, r *http.Request) {
	referenceP := chi.URLParam(r, "reference")

	result, err := s.service.GetLogs(r.Context(), referenceP)
	response(w, result, err)
}

func (s *server) listPlaybooksHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListPlaybooks(r.Context())
	response(w, result, err)
}

func (s *server) createPlaybookHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.PlaybookTemplateFormSchema.Validate(jl)

	var playbookP *model.PlaybookTemplateForm
	if err := parseBody(r, &playbookP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreatePlaybook(r.Context(), playbookP)
	response(w, result, err)
}

func (s *server) getPlaybookHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetPlaybook(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updatePlaybookHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.PlaybookTemplateFormSchema.Validate(jl)

	var playbookP *model.PlaybookTemplateForm
	if err := parseBody(r, &playbookP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdatePlaybook(r.Context(), idP, playbookP)
	response(w, result, err)
}

func (s *server) deletePlaybookHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	response(w, nil, s.service.DeletePlaybook(r.Context(), idP))
}

func (s *server) getSettingsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.GetSettings(r.Context())
	response(w, result, err)
}

func (s *server) getStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.GetStatistics(r.Context())
	response(w, result, err)
}

func (s *server) listTasksHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListTasks(r.Context())
	response(w, result, err)
}

func (s *server) listTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListTemplates(r.Context())
	response(w, result, err)
}

func (s *server) createTemplateHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketTemplateFormSchema.Validate(jl)

	var templateP *model.TicketTemplateForm
	if err := parseBody(r, &templateP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreateTemplate(r.Context(), templateP)
	response(w, result, err)
}

func (s *server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetTemplate(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateTemplateHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketTemplateFormSchema.Validate(jl)

	var templateP *model.TicketTemplateForm
	if err := parseBody(r, &templateP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateTemplate(r.Context(), idP, templateP)
	response(w, result, err)
}

func (s *server) deleteTemplateHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	response(w, nil, s.service.DeleteTemplate(r.Context(), idP))
}

func (s *server) listTicketsHandler(w http.ResponseWriter, r *http.Request) {
	typeP := r.URL.Query().Get("type")

	offsetP, err := parseQueryOptionalInt(r, "offset")
	if err != nil {
		JSONError(w, err)
		return
	}

	countP, err := parseQueryOptionalInt(r, "count")
	if err != nil {
		JSONError(w, err)
		return
	}

	sortP, err := parseQueryOptionalStringArray(r, "sort")
	if err != nil {
		JSONError(w, err)
		return
	}

	descP, err := parseQueryOptionalBoolArray(r, "desc")
	if err != nil {
		JSONError(w, err)
		return
	}

	queryP := r.URL.Query().Get("query")

	result, err := s.service.ListTickets(r.Context(), &typeP, offsetP, countP, sortP, descP, &queryP)
	response(w, result, err)
}

func (s *server) createTicketHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketFormSchema.Validate(jl)

	var ticketP *model.TicketForm
	if err := parseBody(r, &ticketP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreateTicket(r.Context(), ticketP)
	response(w, result, err)
}

func (s *server) createTicketBatchHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// []*model.TicketFormSchema.Validate(jl)

	var ticketP []*model.TicketForm
	if err := parseBody(r, &ticketP); err != nil {
		JSONError(w, err)
		return
	}

	response(w, nil, s.service.CreateTicketBatch(r.Context(), ticketP))
}

func (s *server) getTicketHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.GetTicket(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateTicketHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketSchema.Validate(jl)

	var ticketP *model.Ticket
	if err := parseBody(r, &ticketP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateTicket(r.Context(), idP, ticketP)
	response(w, result, err)
}

func (s *server) deleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	response(w, nil, s.service.DeleteTicket(r.Context(), idP))
}

func (s *server) addArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.ArtifactSchema.Validate(jl)

	var artifactP *model.Artifact
	if err := parseBody(r, &artifactP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.AddArtifact(r.Context(), idP, artifactP)
	response(w, result, err)
}

func (s *server) getArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	nameP := chi.URLParam(r, "name")

	result, err := s.service.GetArtifact(r.Context(), idP, nameP)
	response(w, result, err)
}

func (s *server) setArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	nameP := chi.URLParam(r, "name")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.ArtifactSchema.Validate(jl)

	var artifactP *model.Artifact
	if err := parseBody(r, &artifactP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.SetArtifact(r.Context(), idP, nameP, artifactP)
	response(w, result, err)
}

func (s *server) removeArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	nameP := chi.URLParam(r, "name")

	result, err := s.service.RemoveArtifact(r.Context(), idP, nameP)
	response(w, result, err)
}

func (s *server) enrichArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	nameP := chi.URLParam(r, "name")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.EnrichmentFormSchema.Validate(jl)

	var dataP *model.EnrichmentForm
	if err := parseBody(r, &dataP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.EnrichArtifact(r.Context(), idP, nameP, dataP)
	response(w, result, err)
}

func (s *server) runArtifactHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	nameP := chi.URLParam(r, "name")

	automationP := chi.URLParam(r, "automation")

	response(w, nil, s.service.RunArtifact(r.Context(), idP, nameP, automationP))
}

func (s *server) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.CommentFormSchema.Validate(jl)

	var commentP *model.CommentForm
	if err := parseBody(r, &commentP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.AddComment(r.Context(), idP, commentP)
	response(w, result, err)
}

func (s *server) removeCommentHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	commentIDP, err := parseURLInt(r, "commentID")
	if err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.RemoveComment(r.Context(), idP, commentIDP)
	response(w, result, err)
}

func (s *server) linkFilesHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// []*model.FileSchema.Validate(jl)

	var filesP []*model.File
	if err := parseBody(r, &filesP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.LinkFiles(r.Context(), idP, filesP)
	response(w, result, err)
}

func (s *server) addTicketPlaybookHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.PlaybookTemplateFormSchema.Validate(jl)

	var playbookP *model.PlaybookTemplateForm
	if err := parseBody(r, &playbookP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.AddTicketPlaybook(r.Context(), idP, playbookP)
	response(w, result, err)
}

func (s *server) removeTicketPlaybookHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	playbookIDP := chi.URLParam(r, "playbookID")

	result, err := s.service.RemoveTicketPlaybook(r.Context(), idP, playbookIDP)
	response(w, result, err)
}

func (s *server) setTaskHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	playbookIDP := chi.URLParam(r, "playbookID")

	taskIDP := chi.URLParam(r, "taskID")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TaskSchema.Validate(jl)

	var taskP *model.Task
	if err := parseBody(r, &taskP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.SetTask(r.Context(), idP, playbookIDP, taskIDP, taskP)
	response(w, result, err)
}

func (s *server) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	playbookIDP := chi.URLParam(r, "playbookID")

	taskIDP := chi.URLParam(r, "taskID")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// map[string]interface{}Schema.Validate(jl)

	var dataP map[string]interface{}
	if err := parseBody(r, &dataP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CompleteTask(r.Context(), idP, playbookIDP, taskIDP, dataP)
	response(w, result, err)
}

func (s *server) runTaskHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	playbookIDP := chi.URLParam(r, "playbookID")

	taskIDP := chi.URLParam(r, "taskID")

	response(w, nil, s.service.RunTask(r.Context(), idP, playbookIDP, taskIDP))
}

func (s *server) setReferencesHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// []*model.ReferenceSchema.Validate(jl)

	var referencesP []*model.Reference
	if err := parseBody(r, &referencesP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.SetReferences(r.Context(), idP, referencesP)
	response(w, result, err)
}

func (s *server) setSchemaHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// stringSchema.Validate(jl)

	var schemaP string
	if err := parseBody(r, &schemaP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.SetSchema(r.Context(), idP, schemaP)
	response(w, result, err)
}

func (s *server) linkTicketHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// int64Schema.Validate(jl)

	var linkedIDP int64
	if err := parseBody(r, &linkedIDP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.LinkTicket(r.Context(), idP, linkedIDP)
	response(w, result, err)
}

func (s *server) unlinkTicketHandler(w http.ResponseWriter, r *http.Request) {
	idP, err := parseURLInt64(r, "id")
	if err != nil {
		JSONError(w, err)
		return
	}

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// int64Schema.Validate(jl)

	var linkedIDP int64
	if err := parseBody(r, &linkedIDP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UnlinkTicket(r.Context(), idP, linkedIDP)
	response(w, result, err)
}

func (s *server) listTicketTypesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListTicketTypes(r.Context())
	response(w, result, err)
}

func (s *server) createTicketTypeHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketTypeFormSchema.Validate(jl)

	var tickettypeP *model.TicketTypeForm
	if err := parseBody(r, &tickettypeP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreateTicketType(r.Context(), tickettypeP)
	response(w, result, err)
}

func (s *server) getTicketTypeHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetTicketType(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateTicketTypeHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.TicketTypeFormSchema.Validate(jl)

	var tickettypeP *model.TicketTypeForm
	if err := parseBody(r, &tickettypeP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateTicketType(r.Context(), idP, tickettypeP)
	response(w, result, err)
}

func (s *server) deleteTicketTypeHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	response(w, nil, s.service.DeleteTicketType(r.Context(), idP))
}

func (s *server) listUserDataHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListUserData(r.Context())
	response(w, result, err)
}

func (s *server) getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetUserData(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateUserDataHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.UserDataSchema.Validate(jl)

	var userdataP *model.UserData
	if err := parseBody(r, &userdataP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateUserData(r.Context(), idP, userdataP)
	response(w, result, err)
}

func (s *server) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.ListUsers(r.Context())
	response(w, result, err)
}

func (s *server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.UserFormSchema.Validate(jl)

	var userP *model.UserForm
	if err := parseBody(r, &userP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.CreateUser(r.Context(), userP)
	response(w, result, err)
}

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	result, err := s.service.GetUser(r.Context(), idP)
	response(w, result, err)
}

func (s *server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	// jl, _ := gojsonschema.NewReaderLoader(r.Body)
	// *model.UserFormSchema.Validate(jl)

	var userP *model.UserForm
	if err := parseBody(r, &userP); err != nil {
		JSONError(w, err)
		return
	}

	result, err := s.service.UpdateUser(r.Context(), idP, userP)
	response(w, result, err)
}

func (s *server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idP := chi.URLParam(r, "id")

	response(w, nil, s.service.DeleteUser(r.Context(), idP))
}

func parseURLInt64(r *http.Request, s string) (int64, error) {
	i, err := strconv.ParseInt(chi.URLParam(r, s), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseURLInt(r *http.Request, s string) (int, error) {
	i, err := strconv.Atoi(chi.URLParam(r, s))
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseQueryInt(r *http.Request, s string) (int, error) {
	i, err := strconv.Atoi(r.URL.Query().Get(s))
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseQueryBool(r *http.Request, s string) (bool, error) {
	b, err := strconv.ParseBool(r.URL.Query().Get(s))
	if err != nil {
		return false, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return b, nil
}

func parseQueryStringArray(r *http.Request, key string) ([]string, error) {
	stringArray, ok := r.URL.Query()[key]
	if !ok {
		return nil, nil
	}
	return stringArray, nil
}

func parseQueryBoolArray(r *http.Request, key string) ([]bool, error) {
	stringArray, ok := r.URL.Query()[key]
	if !ok {
		return nil, nil
	}
	var boolArray []bool
	for _, s := range stringArray {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
		}
		boolArray = append(boolArray, b)
	}

	return boolArray, nil
}

func parseQueryOptionalInt(r *http.Request, key string) (*int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return nil, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return &i, nil
}

func parseQueryOptionalStringArray(r *http.Request, key string) ([]string, error) {
	return parseQueryStringArray(r, key)
}

func parseQueryOptionalBoolArray(r *http.Request, key string) ([]bool, error) {
	return parseQueryBoolArray(r, key)
}

func parseBody(r *http.Request, i interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(i)
	if err != nil {
		return fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return nil
}

func JSONError(w http.ResponseWriter, err error) {
	JSONErrorStatus(w, http.StatusInternalServerError, err)
}

func JSONErrorStatus(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(b)
}

func response(w http.ResponseWriter, v interface{}, err error) {
	if err != nil {
		var httpError *HTTPError
		if errors.As(err, &httpError) {
			JSONErrorStatus(w, httpError.Status, httpError.Internal)
		}
		log.Println(err)
		JSONError(w, err)
		return
	}

	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(v)
	w.Write(b)
}
