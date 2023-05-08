package test_engine

import (
	"fmt"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/service_network/partition_topology"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/builtin_argument"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/kurtosis_types/packet_delay_distribution"
	"github.com/stretchr/testify/require"
	"testing"
)

type normalPacketDelayDistributionFullTestCase struct {
	*testing.T
}

func newNormalPacketDelayDistributionFullTestCase(t *testing.T) *normalPacketDelayDistributionFullTestCase {
	return &normalPacketDelayDistributionFullTestCase{
		T: t,
	}
}

func (t *normalPacketDelayDistributionFullTestCase) GetId() string {
	return fmt.Sprintf("%s_%s", packet_delay_distribution.NormalPacketDelayDistributionTypeName, "Full")
}

func (t *normalPacketDelayDistributionFullTestCase) GetStarlarkCode() string {
	return fmt.Sprintf("%s(%s=%d, %s=%d, %s=%s)", packet_delay_distribution.NormalPacketDelayDistributionTypeName, packet_delay_distribution.MeanAttr, 110, packet_delay_distribution.StdDevAttr, 16, packet_delay_distribution.CorrelationAttr, "13.4")
}

func (t *normalPacketDelayDistributionFullTestCase) Assert(typeValue builtin_argument.KurtosisValueType) {
	normalPacketDelayDistributionStarlark, ok := typeValue.(*packet_delay_distribution.NormalPacketDelayDistribution)
	require.True(t, ok)
	normalPacketDelayDistribution, err := normalPacketDelayDistributionStarlark.ToKurtosisType()
	require.Nil(t, err)

	expectedNormalPacketDelayDistribution := partition_topology.NewNormalPacketDelayDistribution(110, 16, 13.4)
	require.Equal(t, expectedNormalPacketDelayDistribution, *normalPacketDelayDistribution)
}
