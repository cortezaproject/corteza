package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/rest/handlers"
	"github.com/crusttech/human/server/system/service"
)

func MountRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			handlers.NewLocale(Locale{}.New()).MountRoutes(r)

			handlers.NewAttachment(Attachment{}.New()).MountRoutes(r)
			handlers.NewAuth((Auth{}).New()).MountRoutes(r)

			// A special case that, we do not add this through standard request, handlers & controllers
			// combo but directly -- we need access to r.Body
			r.Handle(service.SinkBaseURL+"*", &Sink{
				svc:  service.DefaultSink,
				sign: auth.DefaultSigner,
			})
		})

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.HttpTokenValidator("api"))

			handlers.NewAuthClient(AuthClient{}.New()).MountRoutes(r)
			handlers.NewAutomation(Automation{}.New()).MountRoutes(r)
			handlers.NewUser(User{}.New()).MountRoutes(r)
			handlers.NewDalConnection(DalConnection{}.New()).MountRoutes(r)
			handlers.NewDalSensitivityLevel(SensitivityLevel{}.New()).MountRoutes(r)
			handlers.NewDalDriver(DalDriver{}.New()).MountRoutes(r)
			handlers.NewDalSchemaAlteration(DalSchemaAlteration{}.New()).MountRoutes(r)
			handlers.NewRole(Role{}.New()).MountRoutes(r)
			handlers.NewUserGroup(UserGroup{}.New()).MountRoutes(r)
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
			handlers.NewLabel(Label{}.New()).MountRoutes(r)
			handlers.NewApplication(Application{}.New()).MountRoutes(r)
			handlers.NewTemplate(Template{}.New()).MountRoutes(r)
			handlers.NewReport(Report{}.New()).MountRoutes(r)
			handlers.NewSettings(Settings{}.New()).MountRoutes(r)
			handlers.NewStats(Stats{}.New()).MountRoutes(r)
			handlers.NewReminder(Reminder{}.New()).MountRoutes(r)
			handlers.NewActionlog(Actionlog{}.New()).MountRoutes(r)
			handlers.NewNotification(Notification{}.New()).MountRoutes(r)
			handlers.NewQueues(Queue{}.New()).MountRoutes(r)
			handlers.NewApigwRoute(ApigwRoute{}.New()).MountRoutes(r)
			handlers.NewApigwFilter(ApigwFilter{}.New()).MountRoutes(r)
			handlers.NewApigwProfiler(ApigwProfiler{}.New()).MountRoutes(r)
			handlers.NewDataPrivacy(DataPrivacy{}.New()).MountRoutes(r)
			handlers.NewConnection(Connection{}.New()).MountRoutes(r)
			handlers.NewConfiguredConnection(ConfiguredConnection{}.New()).MountRoutes(r)
			handlers.NewSmtpConfigurationChecker(SmtpConfigurationChecker{}.New()).MountRoutes(r)
			handlers.NewExpression(Expression{}.New()).MountRoutes(r)
			handlers.NewMcp(Mcp{}.New()).MountRoutes(r)
			handlers.NewAgent(Agent{}.New()).MountRoutes(r)
			handlers.NewAiConversation(AiConversation{}.New()).MountRoutes(r)
			handlers.NewLlmProvider(LlmProvider{}.New()).MountRoutes(r)
			handlers.NewKnowledgeBase(KnowledgeBase{}.New()).MountRoutes(r)
			handlers.NewChatbot(Chatbot{}.New()).MountRoutes(r)

			handlers.NewChatbotSession(ChatbotSession{}.New()).MountRoutes(r)

			handlers.NewProject(Project{}.New()).MountRoutes(r)
			handlers.NewProjectGroup(ProjectGroup{}.New()).MountRoutes(r)

			handlers.NewTenant(Tenant{}.New()).MountRoutes(r)

			// Admin-authed chatbot preview API. Mirrors /api/widget/v1 but
			// uses an in-memory FIFO store of inline chatbot configs so
			// drafts can be exercised without persisting.
			//
			// @todo can/should we pipe this via our standard approach instead?
			NewChatbotPreviewController().MountRoutes(r)
		})
	}
}
