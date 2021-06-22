package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearMessagebusQueueSettings() {
	h.noError(store.TruncateMessagebusQueueSettings(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeMessagebusQueueSetting(consumer ...string) *messagebus.QueueSettings {
	res := &messagebus.QueueSettings{
		ID:        id.Next(),
		Queue:     rs(),
		CreatedAt: time.Now(),
	}

	if len(consumer) == 0 {
		res.Consumer = string(messagebus.ConsumerCorteza)
	} else {
		res.Consumer = consumer[0]
	}

	h.a.NoError(store.CreateMessagebusQueueSetting(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupByID(id uint64) *messagebus.QueueSettings {
	res, err := store.LookupMessagebusQueueSettingByID(context.Background(), service.DefaultStore, id)
	h.noError(err)
	return res
}

func (h helper) lookupByQueue(queue string) *messagebus.QueueSettings {
	res, err := store.LookupMessagebusQueueSettingByQueue(context.Background(), service.DefaultStore, queue)
	h.noError(err)
	return res
}

func TestQueueList(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	h.repoMakeMessagebusQueueSetting()
	h.repoMakeMessagebusQueueSetting()

	h.allow(types.MessagebusQueueRBACResource.AppendWildcard(), "read")

	h.apiInit().
		Get("/queues/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestQueueRead(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	res := h.repoMakeMessagebusQueueSetting()

	h.allow(types.MessagebusQueueRBACResource.AppendWildcard(), "read")

	h.apiInit().
		Get(fmt.Sprintf("/queues/%d", res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestQueueCreate(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	h.allow(types.SystemRBACResource, "messagebus-queue.create")

	consumer := string(messagebus.ConsumerStore)
	queue := rs()

	h.apiInit().
		Put("/queues").
		FormData("consumer", consumer).
		FormData("queue", queue).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupByQueue(queue)
	h.a.NotNil(res)
	h.a.Equal(consumer, res.Consumer)
}

func TestQueueUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	consumer := string(messagebus.ConsumerRedis)
	res := h.repoMakeMessagebusQueueSetting()
	res.Consumer = consumer

	h.allow(types.MessagebusQueueRBACResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/queues/%d", res.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(res)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(consumer, res.Consumer)
}

func TestQueueDelete(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	res := h.repoMakeMessagebusQueueSetting()

	h.allow(types.MessagebusQueueRBACResource.AppendWildcard(), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/queues/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestQueueUnDelete(t *testing.T) {
	h := newHelper(t)
	h.clearMessagebusQueueSettings()

	res := h.repoMakeMessagebusQueueSetting()

	h.allow(types.MessagebusQueueRBACResource.AppendWildcard(), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/queues/%d/undelete", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupByID(res.ID)
	h.a.NotNil(res)
	h.a.Nil(res.DeletedAt)
}
