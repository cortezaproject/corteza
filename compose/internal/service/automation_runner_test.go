package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	service_mocks "github.com/cortezaproject/corteza-server/compose/internal/service/mocks"
	"github.com/cortezaproject/corteza-server/internal/test"
	"github.com/cortezaproject/corteza-server/pkg/automation"
)

func Test_automationRunner_findImplicitScripts(t *testing.T) {
	var (
		// Should be a part of users-scripts, with only one (manual) trigger
		s1   = &automation.Script{ID: 1000, Enabled: true}
		s1t1 = &automation.Trigger{ID: 1001, Enabled: true, Resource: "res", Condition: "5555", Event: "manual"}
		s1t2 = &automation.Trigger{ID: 1002, Enabled: true, Resource: "res", Condition: "5555", Event: "beforeCreate"}
		s1t3 = &automation.Trigger{ID: 1003, Enabled: true, Resource: "res", Condition: "5555", Event: "afterDelete"}

		// Should be a part of user-scripts, with all triggers
		s2   = &automation.Script{ID: 2000, Enabled: true, RunInUA: true}
		s2t1 = &automation.Trigger{ID: 2001, Enabled: true, Resource: "res", Condition: "5555", Event: "manual"}
		s2t2 = &automation.Trigger{ID: 2002, Enabled: true, Resource: "res", Condition: "5555", Event: "beforeCreate"}
		s2t3 = &automation.Trigger{ID: 2003, Enabled: true, Resource: "res", Condition: "5555", Event: "afterDelete"}

		// Should not be a part of the user-scripts
		s3 = &automation.Script{ID: 3000, Enabled: true}
	)

	s1.AddTrigger(automation.STMS_REPLACE, s1t1, s1t2, s1t3)
	s2.AddTrigger(automation.STMS_REPLACE, s2t1, s2t2, s2t3)

	var runnables = automation.ScriptSet{s1, s2, s3}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sfMock := service_mocks.NewMockautomationScriptsFinder(mockCtrl)
	sfMock.EXPECT().
		FindRunnableScripts(gomock.Eq(AutomationResourceRecord), gomock.Eq("beforeCreate"), gomock.Any()).
		Return(runnables)

	runner := automationRunner{
		scriptFinder: sfMock,
	}

	ss := runner.findRecordScripts("beforeCreate", 5555)

	test.Assert(t, len(ss) == 2, "Received user scripts do not match")

	test.Assert(t, len(runnables) == 3, "Expected runnable scriptSet to be intact")
	test.Assert(t, len(runnables.FindByID(1000).Triggers()) == 3, "Expected runnable scriptSet (triggers from first script) to be intact")
}

func Test_automationRunner_UserScripts(t *testing.T) {
	var (
		// Should be a part of users-scripts, with only one (manual) trigger
		s1   = &automation.Script{ID: 1000, Enabled: true}
		s1t1 = &automation.Trigger{ID: 1001, Enabled: true, Event: "manual"}
		s1t2 = &automation.Trigger{ID: 1002, Enabled: true, Event: "beforeCreate"}
		s1t3 = &automation.Trigger{ID: 1003, Enabled: true, Event: "afterDelete"}

		// Should be a part of user-scripts, with all triggers
		s2   = &automation.Script{ID: 2000, Enabled: true, RunInUA: true}
		s2t1 = &automation.Trigger{ID: 2001, Enabled: true, Event: "manual"}
		s2t2 = &automation.Trigger{ID: 2002, Enabled: true, Event: "beforeCreate"}
		s2t3 = &automation.Trigger{ID: 2003, Enabled: true, Event: "afterDelete"}

		// Should not be a part of the user-scripts
		s3 = &automation.Script{ID: 3000, Enabled: true}
	)

	s1.AddTrigger(automation.STMS_REPLACE, s1t1, s1t2, s1t3)
	s2.AddTrigger(automation.STMS_REPLACE, s2t1, s2t2, s2t3)

	var runnables = automation.ScriptSet{s1, s2, s3}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	acAllowMock := service_mocks.NewMockautomationRunnerAccessControler(mockCtrl)
	acAllowMock.EXPECT().
		CanRunAutomationTrigger(gomock.Any(), gomock.Any()).
		// Should be equal to number of matching triggers
		Times(4).
		Return(true)

	sfMock := service_mocks.NewMockautomationScriptsFinder(mockCtrl)
	sfMock.EXPECT().
		FindRunnableScripts(gomock.Eq(""), gomock.Eq("")).
		Return(runnables)

	runner := automationRunner{
		ac:           acAllowMock,
		scriptFinder: sfMock,
	}

	ss := runner.UserScripts(context.Background())

	test.Assert(t, len(ss) == 2, "Received user scripts do not match")
	test.Assert(t, len(ss.FindByID(1000).Triggers()) == 1, "Received user script triggers do not match")
	test.Assert(t, len(ss.FindByID(2000).Triggers()) == 3, "Received user script triggers do not match")

	test.Assert(t, len(runnables) == 3, "Expected runnable scriptSet to be intact")
	test.Assert(t, len(runnables.FindByID(1000).Triggers()) == 3, "Expected runnable scriptSet (triggers from first script) to be intact")
}
