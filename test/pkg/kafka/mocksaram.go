package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
)

// the default offset range is 0 ~ 100, topic is "my-topic", partition is 0
func MockSaramaCluster(t *testing.T, messages []string) *sarama.MockBroker {
	// mockFetchResponse := sarama.NewMockFetchResponse(t, 1).
	// 	SetMessage("my-topic", 0, 0, sarama.StringEncoder("foo")).
	// 	SetMessage("my-topic", 0, 1, sarama.StringEncoder("bar")).
	// 	SetMessage("my-topic", 0, 2, sarama.StringEncoder("baz")).
	// 	SetMessage("my-topic", 0, 3, sarama.StringEncoder("qux"))
	oldestOffset := int64(0)
	newestOffset := int64(100)
	mockFetchResponse := sarama.NewMockFetchResponse(t, 1)
	for i, msg := range messages {
		mockFetchResponse = mockFetchResponse.SetMessage("my-topic", 0, oldestOffset+int64(i), sarama.StringEncoder(msg))
	}

	broker0 := sarama.NewMockBroker(t, 0)
	broker0.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(broker0.Addr(), broker0.BrokerID()).
			SetLeader("my-topic", 0, broker0.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset("my-topic", 0, sarama.OffsetOldest, oldestOffset).
			SetOffset("my-topic", 0, sarama.OffsetNewest, newestOffset),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, "my-group", broker0),
		"HeartbeatRequest": sarama.NewMockHeartbeatResponse(t),
		"JoinGroupRequest": sarama.NewMockSequence(
			sarama.NewMockJoinGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockJoinGroupResponse(t).SetGroupProtocol(sarama.RangeBalanceStrategyName),
		),
		"SyncGroupRequest": sarama.NewMockSequence(
			sarama.NewMockSyncGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockSyncGroupResponse(t).SetMemberAssignment(
				&sarama.ConsumerGroupMemberAssignment{
					Version: 0,
					Topics: map[string][]int32{
						"my-topic": {0},
					},
				}),
		),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).SetOffset(
			"my-group", "my-topic", 0, oldestOffset, "", sarama.ErrNoError,
		).SetError(sarama.ErrNoError),
		"FetchRequest": sarama.NewMockSequence(
			mockFetchResponse,
		),
	})
	return broker0
}
