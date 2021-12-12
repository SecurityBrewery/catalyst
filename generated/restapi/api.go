package restapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/automations"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/jobs"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/logs"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/playbooks"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/settings"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/statistics"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tasks"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/templates"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickettypes"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/userdata"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/users"
	"github.com/SecurityBrewery/catalyst/role"
)

// Service is the interface that must be implemented in order to provide
// business logic for the Server service.
type Service interface {
	AddArtifact(ctx context.Context, params *tickets.AddArtifactParams) *api.Response
	AddComment(ctx context.Context, params *tickets.AddCommentParams) *api.Response
	AddTicketPlaybook(ctx context.Context, params *tickets.AddTicketPlaybookParams) *api.Response
	CompleteTask(ctx context.Context, params *tickets.CompleteTaskParams) *api.Response
	CreateAutomation(ctx context.Context, params *automations.CreateAutomationParams) *api.Response
	CreatePlaybook(ctx context.Context, params *playbooks.CreatePlaybookParams) *api.Response
	CreateTemplate(ctx context.Context, params *templates.CreateTemplateParams) *api.Response
	CreateTicket(ctx context.Context, params *tickets.CreateTicketParams) *api.Response
	CreateTicketBatch(ctx context.Context, params *tickets.CreateTicketBatchParams) *api.Response
	CreateTicketType(ctx context.Context, params *tickettypes.CreateTicketTypeParams) *api.Response
	CreateUser(ctx context.Context, params *users.CreateUserParams) *api.Response
	CurrentUser(ctx context.Context) *api.Response
	CurrentUserData(ctx context.Context) *api.Response
	DeleteAutomation(ctx context.Context, params *automations.DeleteAutomationParams) *api.Response
	DeletePlaybook(ctx context.Context, params *playbooks.DeletePlaybookParams) *api.Response
	DeleteTemplate(ctx context.Context, params *templates.DeleteTemplateParams) *api.Response
	DeleteTicket(ctx context.Context, params *tickets.DeleteTicketParams) *api.Response
	DeleteTicketType(ctx context.Context, params *tickettypes.DeleteTicketTypeParams) *api.Response
	DeleteUser(ctx context.Context, params *users.DeleteUserParams) *api.Response
	EnrichArtifact(ctx context.Context, params *tickets.EnrichArtifactParams) *api.Response
	GetArtifact(ctx context.Context, params *tickets.GetArtifactParams) *api.Response
	GetAutomation(ctx context.Context, params *automations.GetAutomationParams) *api.Response
	GetJob(ctx context.Context, params *jobs.GetJobParams) *api.Response
	GetLogs(ctx context.Context, params *logs.GetLogsParams) *api.Response
	GetPlaybook(ctx context.Context, params *playbooks.GetPlaybookParams) *api.Response
	GetSettings(ctx context.Context) *api.Response
	GetStatistics(ctx context.Context) *api.Response
	GetTemplate(ctx context.Context, params *templates.GetTemplateParams) *api.Response
	GetTicket(ctx context.Context, params *tickets.GetTicketParams) *api.Response
	GetTicketType(ctx context.Context, params *tickettypes.GetTicketTypeParams) *api.Response
	GetUser(ctx context.Context, params *users.GetUserParams) *api.Response
	GetUserData(ctx context.Context, params *userdata.GetUserDataParams) *api.Response
	LinkFiles(ctx context.Context, params *tickets.LinkFilesParams) *api.Response
	LinkTicket(ctx context.Context, params *tickets.LinkTicketParams) *api.Response
	ListAutomations(ctx context.Context) *api.Response
	ListJobs(ctx context.Context) *api.Response
	ListPlaybooks(ctx context.Context) *api.Response
	ListTasks(ctx context.Context) *api.Response
	ListTemplates(ctx context.Context) *api.Response
	ListTicketTypes(ctx context.Context) *api.Response
	ListTickets(ctx context.Context, params *tickets.ListTicketsParams) *api.Response
	ListUserData(ctx context.Context) *api.Response
	ListUsers(ctx context.Context) *api.Response
	RemoveArtifact(ctx context.Context, params *tickets.RemoveArtifactParams) *api.Response
	RemoveComment(ctx context.Context, params *tickets.RemoveCommentParams) *api.Response
	RemoveTicketPlaybook(ctx context.Context, params *tickets.RemoveTicketPlaybookParams) *api.Response
	RunArtifact(ctx context.Context, params *tickets.RunArtifactParams) *api.Response
	RunJob(ctx context.Context, params *jobs.RunJobParams) *api.Response
	RunTask(ctx context.Context, params *tickets.RunTaskParams) *api.Response
	SetArtifact(ctx context.Context, params *tickets.SetArtifactParams) *api.Response
	SetReferences(ctx context.Context, params *tickets.SetReferencesParams) *api.Response
	SetSchema(ctx context.Context, params *tickets.SetSchemaParams) *api.Response
	SetTask(ctx context.Context, params *tickets.SetTaskParams) *api.Response
	UnlinkTicket(ctx context.Context, params *tickets.UnlinkTicketParams) *api.Response
	UpdateAutomation(ctx context.Context, params *automations.UpdateAutomationParams) *api.Response
	UpdateCurrentUserData(ctx context.Context, params *userdata.UpdateCurrentUserDataParams) *api.Response
	UpdateJob(ctx context.Context, params *jobs.UpdateJobParams) *api.Response
	UpdatePlaybook(ctx context.Context, params *playbooks.UpdatePlaybookParams) *api.Response
	UpdateTemplate(ctx context.Context, params *templates.UpdateTemplateParams) *api.Response
	UpdateTicket(ctx context.Context, params *tickets.UpdateTicketParams) *api.Response
	UpdateTicketType(ctx context.Context, params *tickettypes.UpdateTicketTypeParams) *api.Response
	UpdateUser(ctx context.Context, params *users.UpdateUserParams) *api.Response
	UpdateUserData(ctx context.Context, params *userdata.UpdateUserDataParams) *api.Response
}

// Config defines the config options for the API server.
type Config struct {
	Address      string
	InsecureHTTP bool
	TLSCertFile  string
	TLSKeyFile   string
}

// Server defines the Server service.
type Server struct {
	*gin.Engine
	config  *Config
	server  *http.Server
	service Service

	ApiGroup *gin.RouterGroup

	RoleAuth func([]role.Role) gin.HandlerFunc
}

// New initializes a new Server service.
func New(svc Service, config *Config) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Server{
		Engine:  engine,
		service: svc,
		config:  config,
		server: &http.Server{
			Addr:         config.Address,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},

		ApiGroup: engine.Group("/api"),

		RoleAuth: func(i []role.Role) gin.HandlerFunc { return func(c *gin.Context) { c.Next() } },
	}
}

// ConfigureRoutes configures the routes for the Server service.
// Configuring of routes includes setting up Auth if it is enabled.
func (s *Server) ConfigureRoutes() {
	s.ApiGroup.POST("/tickets/:id/artifacts", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.AddArtifactEndpoint(s.service.AddArtifact))
	s.ApiGroup.POST("/tickets/:id/comments", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.AddCommentEndpoint(s.service.AddComment))
	s.ApiGroup.POST("/tickets/:id/playbooks", s.RoleAuth([]role.Role{}), tickets.AddTicketPlaybookEndpoint(s.service.AddTicketPlaybook))
	s.ApiGroup.PUT("/tickets/:id/playbooks/:playbookID/task/:taskID/complete", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.CompleteTaskEndpoint(s.service.CompleteTask))
	s.ApiGroup.POST("/automations", s.RoleAuth([]role.Role{role.AutomationWrite}), automations.CreateAutomationEndpoint(s.service.CreateAutomation))
	s.ApiGroup.POST("/playbooks", s.RoleAuth([]role.Role{role.PlaybookWrite}), playbooks.CreatePlaybookEndpoint(s.service.CreatePlaybook))
	s.ApiGroup.POST("/templates", s.RoleAuth([]role.Role{role.TemplateWrite}), templates.CreateTemplateEndpoint(s.service.CreateTemplate))
	s.ApiGroup.POST("/tickets", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.CreateTicketEndpoint(s.service.CreateTicket))
	s.ApiGroup.POST("/tickets/batch", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.CreateTicketBatchEndpoint(s.service.CreateTicketBatch))
	s.ApiGroup.POST("/tickettypes", s.RoleAuth([]role.Role{role.TickettypeWrite}), tickettypes.CreateTicketTypeEndpoint(s.service.CreateTicketType))
	s.ApiGroup.POST("/users", s.RoleAuth([]role.Role{role.UserWrite}), users.CreateUserEndpoint(s.service.CreateUser))
	s.ApiGroup.GET("/currentuser", s.RoleAuth([]role.Role{role.CurrentuserRead}), users.CurrentUserEndpoint(s.service.CurrentUser))
	s.ApiGroup.GET("/currentuserdata", s.RoleAuth([]role.Role{role.CurrentuserdataRead}), userdata.CurrentUserDataEndpoint(s.service.CurrentUserData))
	s.ApiGroup.DELETE("/automations/:id", s.RoleAuth([]role.Role{role.AutomationWrite}), automations.DeleteAutomationEndpoint(s.service.DeleteAutomation))
	s.ApiGroup.DELETE("/playbooks/:id", s.RoleAuth([]role.Role{role.PlaybookWrite}), playbooks.DeletePlaybookEndpoint(s.service.DeletePlaybook))
	s.ApiGroup.DELETE("/templates/:id", s.RoleAuth([]role.Role{role.TemplateWrite}), templates.DeleteTemplateEndpoint(s.service.DeleteTemplate))
	s.ApiGroup.DELETE("/tickets/:id", s.RoleAuth([]role.Role{role.TicketDelete}), tickets.DeleteTicketEndpoint(s.service.DeleteTicket))
	s.ApiGroup.DELETE("/tickettypes/:id", s.RoleAuth([]role.Role{role.TickettypeWrite}), tickettypes.DeleteTicketTypeEndpoint(s.service.DeleteTicketType))
	s.ApiGroup.DELETE("/users/:id", s.RoleAuth([]role.Role{role.UserWrite}), users.DeleteUserEndpoint(s.service.DeleteUser))
	s.ApiGroup.POST("/tickets/:id/artifacts/:name/enrich", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.EnrichArtifactEndpoint(s.service.EnrichArtifact))
	s.ApiGroup.GET("/tickets/:id/artifacts/:name", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.GetArtifactEndpoint(s.service.GetArtifact))
	s.ApiGroup.GET("/automations/:id", s.RoleAuth([]role.Role{role.AutomationRead}), automations.GetAutomationEndpoint(s.service.GetAutomation))
	s.ApiGroup.GET("/jobs/:id", s.RoleAuth([]role.Role{role.JobRead}), jobs.GetJobEndpoint(s.service.GetJob))
	s.ApiGroup.GET("/logs/:reference", s.RoleAuth([]role.Role{role.LogRead}), logs.GetLogsEndpoint(s.service.GetLogs))
	s.ApiGroup.GET("/playbooks/:id", s.RoleAuth([]role.Role{role.PlaybookRead}), playbooks.GetPlaybookEndpoint(s.service.GetPlaybook))
	s.ApiGroup.GET("/settings", s.RoleAuth([]role.Role{role.SettingsRead}), settings.GetSettingsEndpoint(s.service.GetSettings))
	s.ApiGroup.GET("/statistics", s.RoleAuth([]role.Role{role.TicketRead}), statistics.GetStatisticsEndpoint(s.service.GetStatistics))
	s.ApiGroup.GET("/templates/:id", s.RoleAuth([]role.Role{role.TemplateRead}), templates.GetTemplateEndpoint(s.service.GetTemplate))
	s.ApiGroup.GET("/tickets/:id", s.RoleAuth([]role.Role{role.TicketRead}), tickets.GetTicketEndpoint(s.service.GetTicket))
	s.ApiGroup.GET("/tickettypes/:id", s.RoleAuth([]role.Role{role.TickettypeRead}), tickettypes.GetTicketTypeEndpoint(s.service.GetTicketType))
	s.ApiGroup.GET("/users/:id", s.RoleAuth([]role.Role{role.UserRead}), users.GetUserEndpoint(s.service.GetUser))
	s.ApiGroup.GET("/userdata/:id", s.RoleAuth([]role.Role{role.UserdataRead}), userdata.GetUserDataEndpoint(s.service.GetUserData))
	s.ApiGroup.PUT("/tickets/:id/files", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.LinkFilesEndpoint(s.service.LinkFiles))
	s.ApiGroup.PATCH("/tickets/:id/tickets", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.LinkTicketEndpoint(s.service.LinkTicket))
	s.ApiGroup.GET("/automations", s.RoleAuth([]role.Role{role.AutomationRead}), automations.ListAutomationsEndpoint(s.service.ListAutomations))
	s.ApiGroup.GET("/jobs", s.RoleAuth([]role.Role{role.JobRead}), jobs.ListJobsEndpoint(s.service.ListJobs))
	s.ApiGroup.GET("/playbooks", s.RoleAuth([]role.Role{role.PlaybookRead}), playbooks.ListPlaybooksEndpoint(s.service.ListPlaybooks))
	s.ApiGroup.GET("/tasks", s.RoleAuth([]role.Role{role.TicketRead}), tasks.ListTasksEndpoint(s.service.ListTasks))
	s.ApiGroup.GET("/templates", s.RoleAuth([]role.Role{role.TemplateRead}), templates.ListTemplatesEndpoint(s.service.ListTemplates))
	s.ApiGroup.GET("/tickettypes", s.RoleAuth([]role.Role{role.TickettypeRead}), tickettypes.ListTicketTypesEndpoint(s.service.ListTicketTypes))
	s.ApiGroup.GET("/tickets", s.RoleAuth([]role.Role{role.TicketRead}), tickets.ListTicketsEndpoint(s.service.ListTickets))
	s.ApiGroup.GET("/userdata", s.RoleAuth([]role.Role{role.UserdataRead}), userdata.ListUserDataEndpoint(s.service.ListUserData))
	s.ApiGroup.GET("/users", s.RoleAuth([]role.Role{role.UserRead}), users.ListUsersEndpoint(s.service.ListUsers))
	s.ApiGroup.DELETE("/tickets/:id/artifacts/:name", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.RemoveArtifactEndpoint(s.service.RemoveArtifact))
	s.ApiGroup.DELETE("/tickets/:id/comments/:commentID", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.RemoveCommentEndpoint(s.service.RemoveComment))
	s.ApiGroup.DELETE("/tickets/:id/playbooks/:playbookID", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.RemoveTicketPlaybookEndpoint(s.service.RemoveTicketPlaybook))
	s.ApiGroup.POST("/tickets/:id/artifacts/:name/run/:automation", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.RunArtifactEndpoint(s.service.RunArtifact))
	s.ApiGroup.POST("/jobs", s.RoleAuth([]role.Role{role.JobWrite}), jobs.RunJobEndpoint(s.service.RunJob))
	s.ApiGroup.POST("/tickets/:id/playbooks/:playbookID/task/:taskID/run", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.RunTaskEndpoint(s.service.RunTask))
	s.ApiGroup.PUT("/tickets/:id/artifacts/:name", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.SetArtifactEndpoint(s.service.SetArtifact))
	s.ApiGroup.PUT("/tickets/:id/references", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.SetReferencesEndpoint(s.service.SetReferences))
	s.ApiGroup.PUT("/tickets/:id/schema", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.SetSchemaEndpoint(s.service.SetSchema))
	s.ApiGroup.PUT("/tickets/:id/playbooks/:playbookID/task/:taskID", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.SetTaskEndpoint(s.service.SetTask))
	s.ApiGroup.DELETE("/tickets/:id/tickets", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.UnlinkTicketEndpoint(s.service.UnlinkTicket))
	s.ApiGroup.PUT("/automations/:id", s.RoleAuth([]role.Role{role.AutomationWrite}), automations.UpdateAutomationEndpoint(s.service.UpdateAutomation))
	s.ApiGroup.PUT("/currentuserdata", s.RoleAuth([]role.Role{role.CurrentuserdataWrite}), userdata.UpdateCurrentUserDataEndpoint(s.service.UpdateCurrentUserData))
	s.ApiGroup.PUT("/jobs/:id", s.RoleAuth([]role.Role{role.JobWrite}), jobs.UpdateJobEndpoint(s.service.UpdateJob))
	s.ApiGroup.PUT("/playbooks/:id", s.RoleAuth([]role.Role{role.PlaybookWrite}), playbooks.UpdatePlaybookEndpoint(s.service.UpdatePlaybook))
	s.ApiGroup.PUT("/templates/:id", s.RoleAuth([]role.Role{role.TemplateWrite}), templates.UpdateTemplateEndpoint(s.service.UpdateTemplate))
	s.ApiGroup.PUT("/tickets/:id", s.RoleAuth([]role.Role{role.TicketWrite}), tickets.UpdateTicketEndpoint(s.service.UpdateTicket))
	s.ApiGroup.PUT("/tickettypes/:id", s.RoleAuth([]role.Role{role.TickettypeWrite}), tickettypes.UpdateTicketTypeEndpoint(s.service.UpdateTicketType))
	s.ApiGroup.PUT("/users/:id", s.RoleAuth([]role.Role{role.UserWrite}), users.UpdateUserEndpoint(s.service.UpdateUser))
	s.ApiGroup.PUT("/userdata/:id", s.RoleAuth([]role.Role{role.UserdataWrite}), userdata.UpdateUserDataEndpoint(s.service.UpdateUserData))
}

// run the Server. It will listen on either HTTP or HTTPS depending on the
// config passed to NewServer.
func (s *Server) run() error {
	log.Printf("Serving on address %s\n", s.server.Addr)
	if s.config.InsecureHTTP {
		return s.server.ListenAndServe()
	}
	return s.server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
}

// Shutdown will gracefully shutdown the Server.
func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

// RunWithSigHandler runs the Server with SIGTERM handling automatically
// enabled. The server will listen for a SIGTERM signal and gracefully shutdown
// the web server.
// It's possible to optionally pass any number shutdown functions which will
// execute one by one after the webserver has been shutdown successfully.
func (s *Server) RunWithSigHandler(shutdown ...func() error) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		s.Shutdown()
	}()

	err := s.run()
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	for _, fn := range shutdown {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}
