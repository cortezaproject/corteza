package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/rest/handlers"
	"github.com/cortezaproject/corteza-server/system/service"
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
			handlers.NewRole(Role{}.New()).MountRoutes(r)
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
			handlers.NewApplication(Application{}.New()).MountRoutes(r)
			handlers.NewTemplate(Template{}.New()).MountRoutes(r)
			handlers.NewReport(Report{}.New()).MountRoutes(r)
			handlers.NewSettings(Settings{}.New()).MountRoutes(r)
			handlers.NewStats(Stats{}.New()).MountRoutes(r)
			handlers.NewReminder(Reminder{}.New()).MountRoutes(r)
			handlers.NewActionlog(Actionlog{}.New()).MountRoutes(r)
			handlers.NewQueues(Queue{}.New()).MountRoutes(r)
			handlers.NewApigwRoute(ApigwRoute{}.New()).MountRoutes(r)
			handlers.NewApigwFilter(ApigwFilter{}.New()).MountRoutes(r)
			handlers.NewApigwProfiler(ApigwProfiler{}.New()).MountRoutes(r)
			handlers.NewDataPrivacy(DataPrivacy{}.New()).MountRoutes(r)
		})
	}
}
