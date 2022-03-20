package model

import (
	"time"

	"github.com/xeipuuv/gojsonschema"
)

var (
	schemaLoader                   = gojsonschema.NewSchemaLoader()
	ArtifactSchema                 = new(gojsonschema.Schema)
	ArtifactOriginSchema           = new(gojsonschema.Schema)
	AutomationSchema               = new(gojsonschema.Schema)
	AutomationFormSchema           = new(gojsonschema.Schema)
	AutomationResponseSchema       = new(gojsonschema.Schema)
	CommentSchema                  = new(gojsonschema.Schema)
	CommentFormSchema              = new(gojsonschema.Schema)
	ContextSchema                  = new(gojsonschema.Schema)
	DashboardSchema                = new(gojsonschema.Schema)
	DashboardResponseSchema        = new(gojsonschema.Schema)
	EnrichmentSchema               = new(gojsonschema.Schema)
	EnrichmentFormSchema           = new(gojsonschema.Schema)
	FileSchema                     = new(gojsonschema.Schema)
	JobSchema                      = new(gojsonschema.Schema)
	JobFormSchema                  = new(gojsonschema.Schema)
	JobResponseSchema              = new(gojsonschema.Schema)
	JobUpdateSchema                = new(gojsonschema.Schema)
	LogEntrySchema                 = new(gojsonschema.Schema)
	MessageSchema                  = new(gojsonschema.Schema)
	NewUserResponseSchema          = new(gojsonschema.Schema)
	OriginSchema                   = new(gojsonschema.Schema)
	PlaybookSchema                 = new(gojsonschema.Schema)
	PlaybookResponseSchema         = new(gojsonschema.Schema)
	PlaybookTemplateSchema         = new(gojsonschema.Schema)
	PlaybookTemplateFormSchema     = new(gojsonschema.Schema)
	PlaybookTemplateResponseSchema = new(gojsonschema.Schema)
	ReferenceSchema                = new(gojsonschema.Schema)
	ReferenceArraySchema           = new(gojsonschema.Schema)
	SettingsSchema                 = new(gojsonschema.Schema)
	SettingsResponseSchema         = new(gojsonschema.Schema)
	StatisticsSchema               = new(gojsonschema.Schema)
	TaskSchema                     = new(gojsonschema.Schema)
	TaskOriginSchema               = new(gojsonschema.Schema)
	TaskResponseSchema             = new(gojsonschema.Schema)
	TaskWithContextSchema          = new(gojsonschema.Schema)
	TicketSchema                   = new(gojsonschema.Schema)
	TicketFormSchema               = new(gojsonschema.Schema)
	TicketFormArraySchema          = new(gojsonschema.Schema)
	TicketListSchema               = new(gojsonschema.Schema)
	TicketResponseSchema           = new(gojsonschema.Schema)
	TicketSimpleResponseSchema     = new(gojsonschema.Schema)
	TicketTemplateSchema           = new(gojsonschema.Schema)
	TicketTemplateFormSchema       = new(gojsonschema.Schema)
	TicketTemplateResponseSchema   = new(gojsonschema.Schema)
	TicketTypeSchema               = new(gojsonschema.Schema)
	TicketTypeFormSchema           = new(gojsonschema.Schema)
	TicketTypeResponseSchema       = new(gojsonschema.Schema)
	TicketWithTicketsSchema        = new(gojsonschema.Schema)
	TypeSchema                     = new(gojsonschema.Schema)
	UserSchema                     = new(gojsonschema.Schema)
	UserDataSchema                 = new(gojsonschema.Schema)
	UserDataResponseSchema         = new(gojsonschema.Schema)
	UserFormSchema                 = new(gojsonschema.Schema)
	UserResponseSchema             = new(gojsonschema.Schema)
	WidgetSchema                   = new(gojsonschema.Schema)
)

func init() {
	err := schemaLoader.AddSchemas(
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"enrichments":{"type":"object","additionalProperties":{"$ref":"#/definitions/Enrichment"}},"kind":{"type":"string"},"name":{"type":"string"},"status":{"type":"string"},"type":{"type":"string"}},"required":["name"],"$id":"#/definitions/Artifact"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifact":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["ticket_id","artifact"],"$id":"#/definitions/ArtifactOrigin"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"image":{"type":"string"},"schema":{"type":"string"},"script":{"type":"string"},"type":{"items":{"type":"string","enum":["artifact","playbook","global"]},"type":"array"}},"required":["image","script","type"],"$id":"#/definitions/Automation"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"image":{"type":"string"},"schema":{"type":"string"},"script":{"type":"string"},"type":{"items":{"type":"string","enum":["artifact","playbook","global"]},"type":"array"}},"required":["id","image","script","type"],"$id":"#/definitions/AutomationForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"image":{"type":"string"},"schema":{"type":"string"},"script":{"type":"string"},"type":{"items":{"type":"string","enum":["artifact","playbook","global"]},"type":"array"}},"required":["id","image","script","type"],"$id":"#/definitions/AutomationResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"}},"required":["creator","created","message"],"$id":"#/definitions/Comment"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"}},"required":["message"],"$id":"#/definitions/CommentForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifact":{"$ref":"#/definitions/Artifact"},"playbook":{"$ref":"#/definitions/PlaybookResponse"},"task":{"$ref":"#/definitions/TaskResponse"},"ticket":{"$ref":"#/definitions/TicketResponse"}},"$id":"#/definitions/Context"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"name":{"type":"string"},"widgets":{"items":{"$ref":"#/definitions/Widget"},"type":"array"}},"required":["name","widgets"],"$id":"#/definitions/Dashboard"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"widgets":{"items":{"$ref":"#/definitions/Widget"},"type":"array"}},"required":["id","name","widgets"],"$id":"#/definitions/DashboardResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"created":{"format":"date-time","type":"string"},"data":{"type":"object"},"name":{"type":"string"}},"required":["name","data","created"],"$id":"#/definitions/Enrichment"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"data":{"type":"object"},"name":{"type":"string"}},"required":["name","data"],"$id":"#/definitions/EnrichmentForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"key":{"type":"string"},"name":{"type":"string"}},"required":["key","name"],"$id":"#/definitions/File"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"automation":{"type":"string"},"container":{"type":"string"},"log":{"type":"string"},"origin":{"$ref":"#/definitions/Origin"},"output":{"type":"object"},"payload":{},"running":{"type":"boolean"},"status":{"type":"string"}},"required":["automation","running","status"],"$id":"#/definitions/Job"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"automation":{"type":"string"},"origin":{"$ref":"#/definitions/Origin"},"payload":{}},"required":["automation"],"$id":"#/definitions/JobForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"automation":{"type":"string"},"container":{"type":"string"},"id":{"type":"string"},"log":{"type":"string"},"origin":{"$ref":"#/definitions/Origin"},"output":{"type":"object"},"payload":{},"status":{"type":"string"}},"required":["id","automation","status"],"$id":"#/definitions/JobResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"container":{"type":"string"},"log":{"type":"string"},"output":{"type":"object"},"running":{"type":"boolean"},"status":{"type":"string"}},"required":["running","status"],"$id":"#/definitions/JobUpdate"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"},"reference":{"type":"string"},"type":{"type":"string"}},"required":["type","reference","creator","created","message"],"$id":"#/definitions/LogEntry"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"context":{"$ref":"#/definitions/Context"},"payload":{},"secrets":{"type":"object","additionalProperties":{"type":"string"}}},"$id":"#/definitions/Message"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"blocked":{"type":"boolean"},"id":{"type":"string"},"roles":{"items":{"type":"string"},"type":"array"},"secret":{"type":"string"}},"required":["id","blocked","roles"],"$id":"#/definitions/NewUserResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifact_origin":{"$ref":"#/definitions/ArtifactOrigin"},"task_origin":{"$ref":"#/definitions/TaskOrigin"}},"$id":"#/definitions/Origin"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"name":{"type":"string"},"tasks":{"type":"object","additionalProperties":{"$ref":"#/definitions/Task"}}},"required":["name","tasks"],"$id":"#/definitions/Playbook"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"name":{"type":"string"},"tasks":{"type":"object","additionalProperties":{"$ref":"#/definitions/TaskResponse"}}},"required":["name","tasks"],"$id":"#/definitions/PlaybookResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"name":{"type":"string"},"yaml":{"type":"string"}},"required":["name","yaml"],"$id":"#/definitions/PlaybookTemplate"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"yaml":{"type":"string"}},"required":["yaml"],"$id":"#/definitions/PlaybookTemplateForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"yaml":{"type":"string"}},"required":["id","name","yaml"],"$id":"#/definitions/PlaybookTemplateResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"href":{"type":"string"},"name":{"type":"string"}},"required":["name","href"],"$id":"#/definitions/Reference"}`),
		gojsonschema.NewStringLoader(`{"items":{"$ref":"#/definitions/Reference"},"type":"array","$id":"#/definitions/ReferenceArray"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifactKinds":{"title":"Artifact Kinds","items":{"$ref":"#/definitions/Type"},"type":"array"},"artifactStates":{"title":"Artifact States","items":{"$ref":"#/definitions/Type"},"type":"array"},"timeformat":{"title":"Time Format","type":"string"}},"required":["timeformat","artifactKinds","artifactStates"],"$id":"#/definitions/Settings"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifactKinds":{"title":"Artifact Kinds","items":{"$ref":"#/definitions/Type"},"type":"array"},"artifactStates":{"title":"Artifact States","items":{"$ref":"#/definitions/Type"},"type":"array"},"roles":{"title":"Roles","items":{"type":"string"},"type":"array"},"ticketTypes":{"title":"Ticket Types","items":{"$ref":"#/definitions/TicketTypeResponse"},"type":"array"},"tier":{"title":"Tier","type":"string","enum":["community","enterprise"]},"timeformat":{"title":"Time Format","type":"string"},"version":{"title":"Version","type":"string"}},"required":["version","tier","timeformat","ticketTypes","artifactKinds","artifactStates"],"$id":"#/definitions/SettingsResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"open_tickets_per_user":{"type":"object","additionalProperties":{"type":"integer"}},"tickets_per_type":{"type":"object","additionalProperties":{"type":"integer"}},"tickets_per_week":{"type":"object","additionalProperties":{"type":"integer"}},"unassigned":{"type":"integer"}},"required":["unassigned","open_tickets_per_user","tickets_per_week","tickets_per_type"],"$id":"#/definitions/Statistics"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"automation":{"type":"string"},"closed":{"format":"date-time","type":"string"},"created":{"format":"date-time","type":"string"},"data":{"type":"object"},"done":{"type":"boolean"},"join":{"type":"boolean"},"name":{"type":"string"},"next":{"type":"object","additionalProperties":{"type":"string"}},"owner":{"type":"string"},"payload":{"type":"object","additionalProperties":{"type":"string"}},"schema":{"type":"object"},"type":{"type":"string","enum":["task","input","automation"]}},"required":["name","type","done","created"],"$id":"#/definitions/Task"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"playbook_id":{"type":"string"},"task_id":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["ticket_id","playbook_id","task_id"],"$id":"#/definitions/TaskOrigin"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"active":{"type":"boolean"},"automation":{"type":"string"},"closed":{"format":"date-time","type":"string"},"created":{"format":"date-time","type":"string"},"data":{"type":"object"},"done":{"type":"boolean"},"join":{"type":"boolean"},"name":{"type":"string"},"next":{"type":"object","additionalProperties":{"type":"string"}},"order":{"format":"int64","type":"number"},"owner":{"type":"string"},"payload":{"type":"object","additionalProperties":{"type":"string"}},"schema":{"type":"object"},"type":{"type":"string","enum":["task","input","automation"]}},"required":["name","type","done","created","order","active"],"$id":"#/definitions/TaskResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"playbook_id":{"type":"string"},"playbook_name":{"type":"string"},"task":{"$ref":"#/definitions/TaskResponse"},"task_id":{"type":"string"},"ticket_id":{"format":"int64","type":"number"},"ticket_name":{"type":"string"}},"required":["ticket_id","ticket_name","playbook_id","playbook_name","task_id","task"],"$id":"#/definitions/TaskWithContext"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifacts":{"items":{"$ref":"#/definitions/Artifact"},"type":"array"},"comments":{"items":{"$ref":"#/definitions/Comment"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"type":"object"},"files":{"items":{"$ref":"#/definitions/File"},"type":"array"},"modified":{"format":"date-time","type":"string"},"name":{"type":"string"},"owner":{"type":"string"},"playbooks":{"type":"object","additionalProperties":{"$ref":"#/definitions/Playbook"}},"read":{"items":{"type":"string"},"type":"array"},"references":{"items":{"$ref":"#/definitions/Reference"},"type":"array"},"schema":{"type":"string"},"status":{"type":"string"},"type":{"type":"string"},"write":{"items":{"type":"string"},"type":"array"}},"required":["name","type","status","created","modified","schema"],"$id":"#/definitions/Ticket"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifacts":{"items":{"$ref":"#/definitions/Artifact"},"type":"array"},"comments":{"items":{"$ref":"#/definitions/Comment"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"type":"object"},"files":{"items":{"$ref":"#/definitions/File"},"type":"array"},"id":{"format":"int64","type":"integer"},"modified":{"format":"date-time","type":"string"},"name":{"type":"string"},"owner":{"type":"string"},"playbooks":{"items":{"$ref":"#/definitions/PlaybookTemplateForm"},"type":"array"},"read":{"items":{"type":"string"},"type":"array"},"references":{"items":{"$ref":"#/definitions/Reference"},"type":"array"},"schema":{"type":"string"},"status":{"type":"string"},"type":{"type":"string"},"write":{"items":{"type":"string"},"type":"array"}},"required":["name","type","status"],"$id":"#/definitions/TicketForm"}`),
		gojsonschema.NewStringLoader(`{"items":{"$ref":"#/definitions/TicketForm"},"type":"array","$id":"#/definitions/TicketFormArray"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"count":{"type":"number"},"tickets":{"items":{"$ref":"#/definitions/TicketSimpleResponse"},"type":"array"}},"required":["tickets","count"],"$id":"#/definitions/TicketList"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifacts":{"items":{"$ref":"#/definitions/Artifact"},"type":"array"},"comments":{"items":{"$ref":"#/definitions/Comment"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"type":"object"},"files":{"items":{"$ref":"#/definitions/File"},"type":"array"},"id":{"format":"int64","type":"integer"},"modified":{"format":"date-time","type":"string"},"name":{"type":"string"},"owner":{"type":"string"},"playbooks":{"type":"object","additionalProperties":{"$ref":"#/definitions/PlaybookResponse"}},"read":{"items":{"type":"string"},"type":"array"},"references":{"items":{"$ref":"#/definitions/Reference"},"type":"array"},"schema":{"type":"string"},"status":{"type":"string"},"type":{"type":"string"},"write":{"items":{"type":"string"},"type":"array"}},"required":["id","name","type","status","created","modified","schema"],"$id":"#/definitions/TicketResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifacts":{"items":{"$ref":"#/definitions/Artifact"},"type":"array"},"comments":{"items":{"$ref":"#/definitions/Comment"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"type":"object"},"files":{"items":{"$ref":"#/definitions/File"},"type":"array"},"id":{"format":"int64","type":"integer"},"modified":{"format":"date-time","type":"string"},"name":{"type":"string"},"owner":{"type":"string"},"playbooks":{"type":"object","additionalProperties":{"$ref":"#/definitions/Playbook"}},"read":{"items":{"type":"string"},"type":"array"},"references":{"items":{"$ref":"#/definitions/Reference"},"type":"array"},"schema":{"type":"string"},"status":{"type":"string"},"type":{"type":"string"},"write":{"items":{"type":"string"},"type":"array"}},"required":["id","name","type","status","created","modified","schema"],"$id":"#/definitions/TicketSimpleResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"name":{"type":"string"},"schema":{"type":"string"}},"required":["name","schema"],"$id":"#/definitions/TicketTemplate"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"schema":{"type":"string"}},"required":["name","schema"],"$id":"#/definitions/TicketTemplateForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"schema":{"type":"string"}},"required":["id","name","schema"],"$id":"#/definitions/TicketTemplateResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"default_groups":{"items":{"type":"string"},"type":"array"},"default_playbooks":{"items":{"type":"string"},"type":"array"},"default_template":{"type":"string"},"icon":{"type":"string"},"name":{"type":"string"}},"required":["name","icon","default_template","default_playbooks"],"$id":"#/definitions/TicketType"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"default_groups":{"items":{"type":"string"},"type":"array"},"default_playbooks":{"items":{"type":"string"},"type":"array"},"default_template":{"type":"string"},"icon":{"type":"string"},"id":{"type":"string"},"name":{"type":"string"}},"required":["name","icon","default_template","default_playbooks"],"$id":"#/definitions/TicketTypeForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"default_groups":{"items":{"type":"string"},"type":"array"},"default_playbooks":{"items":{"type":"string"},"type":"array"},"default_template":{"type":"string"},"icon":{"type":"string"},"id":{"type":"string"},"name":{"type":"string"}},"required":["id","name","icon","default_template","default_playbooks"],"$id":"#/definitions/TicketTypeResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"artifacts":{"items":{"$ref":"#/definitions/Artifact"},"type":"array"},"comments":{"items":{"$ref":"#/definitions/Comment"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"type":"object"},"files":{"items":{"$ref":"#/definitions/File"},"type":"array"},"id":{"format":"int64","type":"integer"},"logs":{"items":{"$ref":"#/definitions/LogEntry"},"type":"array"},"modified":{"format":"date-time","type":"string"},"name":{"type":"string"},"owner":{"type":"string"},"playbooks":{"type":"object","additionalProperties":{"$ref":"#/definitions/PlaybookResponse"}},"read":{"items":{"type":"string"},"type":"array"},"references":{"items":{"$ref":"#/definitions/Reference"},"type":"array"},"schema":{"type":"string"},"status":{"type":"string"},"tickets":{"items":{"$ref":"#/definitions/TicketSimpleResponse"},"type":"array"},"type":{"type":"string"},"write":{"items":{"type":"string"},"type":"array"}},"required":["id","name","type","status","created","modified","schema"],"$id":"#/definitions/TicketWithTickets"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"color":{"title":"Color","type":"string","enum":["error","info","success","warning"]},"icon":{"title":"Icon (https://materialdesignicons.com)","type":"string"},"id":{"title":"ID","type":"string"},"name":{"title":"Name","type":"string"}},"required":["id","name","icon"],"$id":"#/definitions/Type"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"apikey":{"type":"boolean"},"blocked":{"type":"boolean"},"roles":{"items":{"type":"string"},"type":"array"},"sha256":{"type":"string"}},"required":["blocked","apikey","roles"],"$id":"#/definitions/User"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"email":{"type":"string"},"image":{"type":"string"},"name":{"type":"string"},"timeformat":{"title":"Time Format (https://moment.github.io/luxon/docs/manual/formatting.html#table-of-tokens)","type":"string"}},"$id":"#/definitions/UserData"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"email":{"type":"string"},"id":{"type":"string"},"image":{"type":"string"},"name":{"type":"string"},"timeformat":{"title":"Time Format (https://moment.github.io/luxon/docs/manual/formatting.html#table-of-tokens)","type":"string"}},"required":["id"],"$id":"#/definitions/UserDataResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"apikey":{"type":"boolean"},"blocked":{"type":"boolean"},"id":{"type":"string"},"roles":{"items":{"type":"string"},"type":"array"}},"required":["id","blocked","roles","apikey"],"$id":"#/definitions/UserForm"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"apikey":{"type":"boolean"},"blocked":{"type":"boolean"},"id":{"type":"string"},"roles":{"items":{"type":"string"},"type":"array"}},"required":["id","blocked","roles","apikey"],"$id":"#/definitions/UserResponse"}`),
		gojsonschema.NewStringLoader(`{"type":"object","properties":{"aggregation":{"type":"string"},"filter":{"type":"string"},"name":{"type":"string"},"type":{"type":"string","enum":["bar","line","pie"]},"width":{"maximum":12,"type":"integer"}},"required":["name","type","aggregation","width"],"$id":"#/definitions/Widget"}`),
	)
	if err != nil {
		panic(err)
	}

	ArtifactSchema = mustCompile(`#/definitions/Artifact`)
	ArtifactOriginSchema = mustCompile(`#/definitions/ArtifactOrigin`)
	AutomationSchema = mustCompile(`#/definitions/Automation`)
	AutomationFormSchema = mustCompile(`#/definitions/AutomationForm`)
	AutomationResponseSchema = mustCompile(`#/definitions/AutomationResponse`)
	CommentSchema = mustCompile(`#/definitions/Comment`)
	CommentFormSchema = mustCompile(`#/definitions/CommentForm`)
	ContextSchema = mustCompile(`#/definitions/Context`)
	DashboardSchema = mustCompile(`#/definitions/Dashboard`)
	DashboardResponseSchema = mustCompile(`#/definitions/DashboardResponse`)
	EnrichmentSchema = mustCompile(`#/definitions/Enrichment`)
	EnrichmentFormSchema = mustCompile(`#/definitions/EnrichmentForm`)
	FileSchema = mustCompile(`#/definitions/File`)
	JobSchema = mustCompile(`#/definitions/Job`)
	JobFormSchema = mustCompile(`#/definitions/JobForm`)
	JobResponseSchema = mustCompile(`#/definitions/JobResponse`)
	JobUpdateSchema = mustCompile(`#/definitions/JobUpdate`)
	LogEntrySchema = mustCompile(`#/definitions/LogEntry`)
	MessageSchema = mustCompile(`#/definitions/Message`)
	NewUserResponseSchema = mustCompile(`#/definitions/NewUserResponse`)
	OriginSchema = mustCompile(`#/definitions/Origin`)
	PlaybookSchema = mustCompile(`#/definitions/Playbook`)
	PlaybookResponseSchema = mustCompile(`#/definitions/PlaybookResponse`)
	PlaybookTemplateSchema = mustCompile(`#/definitions/PlaybookTemplate`)
	PlaybookTemplateFormSchema = mustCompile(`#/definitions/PlaybookTemplateForm`)
	PlaybookTemplateResponseSchema = mustCompile(`#/definitions/PlaybookTemplateResponse`)
	ReferenceSchema = mustCompile(`#/definitions/Reference`)
	ReferenceArraySchema = mustCompile(`#/definitions/ReferenceArray`)
	SettingsSchema = mustCompile(`#/definitions/Settings`)
	SettingsResponseSchema = mustCompile(`#/definitions/SettingsResponse`)
	StatisticsSchema = mustCompile(`#/definitions/Statistics`)
	TaskSchema = mustCompile(`#/definitions/Task`)
	TaskOriginSchema = mustCompile(`#/definitions/TaskOrigin`)
	TaskResponseSchema = mustCompile(`#/definitions/TaskResponse`)
	TaskWithContextSchema = mustCompile(`#/definitions/TaskWithContext`)
	TicketSchema = mustCompile(`#/definitions/Ticket`)
	TicketFormSchema = mustCompile(`#/definitions/TicketForm`)
	TicketFormArraySchema = mustCompile(`#/definitions/TicketFormArray`)
	TicketListSchema = mustCompile(`#/definitions/TicketList`)
	TicketResponseSchema = mustCompile(`#/definitions/TicketResponse`)
	TicketSimpleResponseSchema = mustCompile(`#/definitions/TicketSimpleResponse`)
	TicketTemplateSchema = mustCompile(`#/definitions/TicketTemplate`)
	TicketTemplateFormSchema = mustCompile(`#/definitions/TicketTemplateForm`)
	TicketTemplateResponseSchema = mustCompile(`#/definitions/TicketTemplateResponse`)
	TicketTypeSchema = mustCompile(`#/definitions/TicketType`)
	TicketTypeFormSchema = mustCompile(`#/definitions/TicketTypeForm`)
	TicketTypeResponseSchema = mustCompile(`#/definitions/TicketTypeResponse`)
	TicketWithTicketsSchema = mustCompile(`#/definitions/TicketWithTickets`)
	TypeSchema = mustCompile(`#/definitions/Type`)
	UserSchema = mustCompile(`#/definitions/User`)
	UserDataSchema = mustCompile(`#/definitions/UserData`)
	UserDataResponseSchema = mustCompile(`#/definitions/UserDataResponse`)
	UserFormSchema = mustCompile(`#/definitions/UserForm`)
	UserResponseSchema = mustCompile(`#/definitions/UserResponse`)
	WidgetSchema = mustCompile(`#/definitions/Widget`)
}

type Artifact struct {
	Enrichments map[string]*Enrichment `json:"enrichments,omitempty"`
	Kind        *string                `json:"kind,omitempty"`
	Name        string                 `json:"name"`
	Status      *string                `json:"status,omitempty"`
	Type        *string                `json:"type,omitempty"`
}

type ArtifactOrigin struct {
	Artifact string `json:"artifact"`
	TicketId int64  `json:"ticket_id"`
}

type Automation struct {
	Image  string   `json:"image"`
	Schema *string  `json:"schema,omitempty"`
	Script string   `json:"script"`
	Type   []string `json:"type"`
}

type AutomationForm struct {
	ID     string   `json:"id"`
	Image  string   `json:"image"`
	Schema *string  `json:"schema,omitempty"`
	Script string   `json:"script"`
	Type   []string `json:"type"`
}

type AutomationResponse struct {
	ID     string   `json:"id"`
	Image  string   `json:"image"`
	Schema *string  `json:"schema,omitempty"`
	Script string   `json:"script"`
	Type   []string `json:"type"`
}

type Comment struct {
	Created time.Time `json:"created"`
	Creator string    `json:"creator"`
	Message string    `json:"message"`
}

type CommentForm struct {
	Created *time.Time `json:"created,omitempty"`
	Creator *string    `json:"creator,omitempty"`
	Message string     `json:"message"`
}

type Context struct {
	Artifact *Artifact         `json:"artifact,omitempty"`
	Playbook *PlaybookResponse `json:"playbook,omitempty"`
	Task     *TaskResponse     `json:"task,omitempty"`
	Ticket   *TicketResponse   `json:"ticket,omitempty"`
}

type Dashboard struct {
	Name    string    `json:"name"`
	Widgets []*Widget `json:"widgets"`
}

type DashboardResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Widgets []*Widget `json:"widgets"`
}

type Enrichment struct {
	Created time.Time      `json:"created"`
	Data    map[string]any `json:"data"`
	Name    string         `json:"name"`
}

type EnrichmentForm struct {
	Data map[string]any `json:"data"`
	Name string         `json:"name"`
}

type File struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Job struct {
	Automation string         `json:"automation"`
	Container  *string        `json:"container,omitempty"`
	Log        *string        `json:"log,omitempty"`
	Origin     *Origin        `json:"origin,omitempty"`
	Output     map[string]any `json:"output,omitempty"`
	Payload    any            `json:"payload,omitempty"`
	Running    bool           `json:"running"`
	Status     string         `json:"status"`
}

type JobForm struct {
	Automation string  `json:"automation"`
	Origin     *Origin `json:"origin,omitempty"`
	Payload    any     `json:"payload,omitempty"`
}

type JobResponse struct {
	Automation string         `json:"automation"`
	Container  *string        `json:"container,omitempty"`
	ID         string         `json:"id"`
	Log        *string        `json:"log,omitempty"`
	Origin     *Origin        `json:"origin,omitempty"`
	Output     map[string]any `json:"output,omitempty"`
	Payload    any            `json:"payload,omitempty"`
	Status     string         `json:"status"`
}

type JobUpdate struct {
	Container *string        `json:"container,omitempty"`
	Log       *string        `json:"log,omitempty"`
	Output    map[string]any `json:"output,omitempty"`
	Running   bool           `json:"running"`
	Status    string         `json:"status"`
}

type LogEntry struct {
	Created   time.Time `json:"created"`
	Creator   string    `json:"creator"`
	Message   string    `json:"message"`
	Reference string    `json:"reference"`
	Type      string    `json:"type"`
}

type Message struct {
	Context *Context          `json:"context,omitempty"`
	Payload any               `json:"payload,omitempty"`
	Secrets map[string]string `json:"secrets,omitempty"`
}

type NewUserResponse struct {
	Blocked bool     `json:"blocked"`
	ID      string   `json:"id"`
	Roles   []string `json:"roles"`
	Secret  *string  `json:"secret,omitempty"`
}

type Origin struct {
	ArtifactOrigin *ArtifactOrigin `json:"artifact_origin,omitempty"`
	TaskOrigin     *TaskOrigin     `json:"task_origin,omitempty"`
}

type Playbook struct {
	Name  string           `json:"name"`
	Tasks map[string]*Task `json:"tasks"`
}

type PlaybookResponse struct {
	Name  string                   `json:"name"`
	Tasks map[string]*TaskResponse `json:"tasks"`
}

type PlaybookTemplate struct {
	Name string `json:"name"`
	Yaml string `json:"yaml"`
}

type PlaybookTemplateForm struct {
	ID   *string `json:"id,omitempty"`
	Yaml string  `json:"yaml"`
}

type PlaybookTemplateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Yaml string `json:"yaml"`
}

type Reference struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type ReferenceArray []*Reference

type Settings struct {
	ArtifactKinds  []*Type `json:"artifactKinds"`
	ArtifactStates []*Type `json:"artifactStates"`
	Timeformat     string  `json:"timeformat"`
}

type SettingsResponse struct {
	ArtifactKinds  []*Type               `json:"artifactKinds"`
	ArtifactStates []*Type               `json:"artifactStates"`
	Roles          []string              `json:"roles,omitempty"`
	TicketTypes    []*TicketTypeResponse `json:"ticketTypes"`
	Tier           string                `json:"tier"`
	Timeformat     string                `json:"timeformat"`
	Version        string                `json:"version"`
}

type Statistics struct {
	OpenTicketsPerUser map[string]int `json:"open_tickets_per_user"`
	TicketsPerType     map[string]int `json:"tickets_per_type"`
	TicketsPerWeek     map[string]int `json:"tickets_per_week"`
	Unassigned         int            `json:"unassigned"`
}

type Task struct {
	Automation *string           `json:"automation,omitempty"`
	Closed     *time.Time        `json:"closed,omitempty"`
	Created    time.Time         `json:"created"`
	Data       map[string]any    `json:"data,omitempty"`
	Done       bool              `json:"done"`
	Join       *bool             `json:"join,omitempty"`
	Name       string            `json:"name"`
	Next       map[string]string `json:"next,omitempty"`
	Owner      *string           `json:"owner,omitempty"`
	Payload    map[string]string `json:"payload,omitempty"`
	Schema     map[string]any    `json:"schema,omitempty"`
	Type       string            `json:"type"`
}

type TaskOrigin struct {
	PlaybookId string `json:"playbook_id"`
	TaskId     string `json:"task_id"`
	TicketId   int64  `json:"ticket_id"`
}

type TaskResponse struct {
	Active     bool              `json:"active"`
	Automation *string           `json:"automation,omitempty"`
	Closed     *time.Time        `json:"closed,omitempty"`
	Created    time.Time         `json:"created"`
	Data       map[string]any    `json:"data,omitempty"`
	Done       bool              `json:"done"`
	Join       *bool             `json:"join,omitempty"`
	Name       string            `json:"name"`
	Next       map[string]string `json:"next,omitempty"`
	Order      int64             `json:"order"`
	Owner      *string           `json:"owner,omitempty"`
	Payload    map[string]string `json:"payload,omitempty"`
	Schema     map[string]any    `json:"schema,omitempty"`
	Type       string            `json:"type"`
}

type TaskWithContext struct {
	PlaybookId   string        `json:"playbook_id"`
	PlaybookName string        `json:"playbook_name"`
	Task         *TaskResponse `json:"task"`
	TaskId       string        `json:"task_id"`
	TicketId     int64         `json:"ticket_id"`
	TicketName   string        `json:"ticket_name"`
}

type Ticket struct {
	Artifacts  []*Artifact          `json:"artifacts,omitempty"`
	Comments   []*Comment           `json:"comments,omitempty"`
	Created    time.Time            `json:"created"`
	Details    map[string]any       `json:"details,omitempty"`
	Files      []*File              `json:"files,omitempty"`
	Modified   time.Time            `json:"modified"`
	Name       string               `json:"name"`
	Owner      *string              `json:"owner,omitempty"`
	Playbooks  map[string]*Playbook `json:"playbooks,omitempty"`
	Read       []string             `json:"read,omitempty"`
	References []*Reference         `json:"references,omitempty"`
	Schema     string               `json:"schema"`
	Status     string               `json:"status"`
	Type       string               `json:"type"`
	Write      []string             `json:"write,omitempty"`
}

type TicketForm struct {
	Artifacts  []*Artifact             `json:"artifacts,omitempty"`
	Comments   []*Comment              `json:"comments,omitempty"`
	Created    *time.Time              `json:"created,omitempty"`
	Details    map[string]any          `json:"details,omitempty"`
	Files      []*File                 `json:"files,omitempty"`
	ID         *int64                  `json:"id,omitempty"`
	Modified   *time.Time              `json:"modified,omitempty"`
	Name       string                  `json:"name"`
	Owner      *string                 `json:"owner,omitempty"`
	Playbooks  []*PlaybookTemplateForm `json:"playbooks,omitempty"`
	Read       []string                `json:"read,omitempty"`
	References []*Reference            `json:"references,omitempty"`
	Schema     *string                 `json:"schema,omitempty"`
	Status     string                  `json:"status"`
	Type       string                  `json:"type"`
	Write      []string                `json:"write,omitempty"`
}

type TicketFormArray []*TicketForm

type TicketList struct {
	Count   int                     `json:"count"`
	Tickets []*TicketSimpleResponse `json:"tickets"`
}

type TicketResponse struct {
	Artifacts  []*Artifact                  `json:"artifacts,omitempty"`
	Comments   []*Comment                   `json:"comments,omitempty"`
	Created    time.Time                    `json:"created"`
	Details    map[string]any               `json:"details,omitempty"`
	Files      []*File                      `json:"files,omitempty"`
	ID         int64                        `json:"id"`
	Modified   time.Time                    `json:"modified"`
	Name       string                       `json:"name"`
	Owner      *string                      `json:"owner,omitempty"`
	Playbooks  map[string]*PlaybookResponse `json:"playbooks,omitempty"`
	Read       []string                     `json:"read,omitempty"`
	References []*Reference                 `json:"references,omitempty"`
	Schema     string                       `json:"schema"`
	Status     string                       `json:"status"`
	Type       string                       `json:"type"`
	Write      []string                     `json:"write,omitempty"`
}

type TicketSimpleResponse struct {
	Artifacts  []*Artifact          `json:"artifacts,omitempty"`
	Comments   []*Comment           `json:"comments,omitempty"`
	Created    time.Time            `json:"created"`
	Details    map[string]any       `json:"details,omitempty"`
	Files      []*File              `json:"files,omitempty"`
	ID         int64                `json:"id"`
	Modified   time.Time            `json:"modified"`
	Name       string               `json:"name"`
	Owner      *string              `json:"owner,omitempty"`
	Playbooks  map[string]*Playbook `json:"playbooks,omitempty"`
	Read       []string             `json:"read,omitempty"`
	References []*Reference         `json:"references,omitempty"`
	Schema     string               `json:"schema"`
	Status     string               `json:"status"`
	Type       string               `json:"type"`
	Write      []string             `json:"write,omitempty"`
}

type TicketTemplate struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type TicketTemplateForm struct {
	ID     *string `json:"id,omitempty"`
	Name   string  `json:"name"`
	Schema string  `json:"schema"`
}

type TicketTemplateResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type TicketType struct {
	DefaultGroups    []string `json:"default_groups,omitempty"`
	DefaultPlaybooks []string `json:"default_playbooks"`
	DefaultTemplate  string   `json:"default_template"`
	Icon             string   `json:"icon"`
	Name             string   `json:"name"`
}

type TicketTypeForm struct {
	DefaultGroups    []string `json:"default_groups,omitempty"`
	DefaultPlaybooks []string `json:"default_playbooks"`
	DefaultTemplate  string   `json:"default_template"`
	Icon             string   `json:"icon"`
	ID               *string  `json:"id,omitempty"`
	Name             string   `json:"name"`
}

type TicketTypeResponse struct {
	DefaultGroups    []string `json:"default_groups,omitempty"`
	DefaultPlaybooks []string `json:"default_playbooks"`
	DefaultTemplate  string   `json:"default_template"`
	Icon             string   `json:"icon"`
	ID               string   `json:"id"`
	Name             string   `json:"name"`
}

type TicketWithTickets struct {
	Artifacts  []*Artifact                  `json:"artifacts,omitempty"`
	Comments   []*Comment                   `json:"comments,omitempty"`
	Created    time.Time                    `json:"created"`
	Details    map[string]any               `json:"details,omitempty"`
	Files      []*File                      `json:"files,omitempty"`
	ID         int64                        `json:"id"`
	Logs       []*LogEntry                  `json:"logs,omitempty"`
	Modified   time.Time                    `json:"modified"`
	Name       string                       `json:"name"`
	Owner      *string                      `json:"owner,omitempty"`
	Playbooks  map[string]*PlaybookResponse `json:"playbooks,omitempty"`
	Read       []string                     `json:"read,omitempty"`
	References []*Reference                 `json:"references,omitempty"`
	Schema     string                       `json:"schema"`
	Status     string                       `json:"status"`
	Tickets    []*TicketSimpleResponse      `json:"tickets,omitempty"`
	Type       string                       `json:"type"`
	Write      []string                     `json:"write,omitempty"`
}

type Type struct {
	Color *string `json:"color,omitempty"`
	Icon  string  `json:"icon"`
	ID    string  `json:"id"`
	Name  string  `json:"name"`
}

type User struct {
	Apikey  bool     `json:"apikey"`
	Blocked bool     `json:"blocked"`
	Roles   []string `json:"roles"`
	Sha256  *string  `json:"sha256,omitempty"`
}

type UserData struct {
	Email      *string `json:"email,omitempty"`
	Image      *string `json:"image,omitempty"`
	Name       *string `json:"name,omitempty"`
	Timeformat *string `json:"timeformat,omitempty"`
}

type UserDataResponse struct {
	Email      *string `json:"email,omitempty"`
	ID         string  `json:"id"`
	Image      *string `json:"image,omitempty"`
	Name       *string `json:"name,omitempty"`
	Timeformat *string `json:"timeformat,omitempty"`
}

type UserForm struct {
	Apikey  bool     `json:"apikey"`
	Blocked bool     `json:"blocked"`
	ID      string   `json:"id"`
	Roles   []string `json:"roles"`
}

type UserResponse struct {
	Apikey  bool     `json:"apikey"`
	Blocked bool     `json:"blocked"`
	ID      string   `json:"id"`
	Roles   []string `json:"roles"`
}

type Widget struct {
	Aggregation string  `json:"aggregation"`
	Filter      *string `json:"filter,omitempty"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Width       int     `json:"width"`
}

func mustCompile(uri string) *gojsonschema.Schema {
	s, err := schemaLoader.Compile(gojsonschema.NewReferenceLoader(uri))
	if err != nil {
		panic(err)
	}
	return s
}

const (
	SettingsResponseTierCommunity = "community"

	SettingsResponseTierEnterprise = "enterprise"

	TaskTypeTask = "task"

	TaskTypeInput = "input"

	TaskTypeAutomation = "automation"

	TaskResponseTypeTask = "task"

	TaskResponseTypeInput = "input"

	TaskResponseTypeAutomation = "automation"

	TypeColorError = "error"

	TypeColorInfo = "info"

	TypeColorSuccess = "success"

	TypeColorWarning = "warning"

	WidgetTypeBar = "bar"

	WidgetTypeLine = "line"

	WidgetTypePie = "pie"
)
