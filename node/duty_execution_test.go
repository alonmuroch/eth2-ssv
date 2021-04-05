package node

import (
	"context"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/ibft/proto"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConsensusOnInputValue(t *testing.T) {
	tests := []struct {
		name                        string
		decided                     bool
		signaturesCount             int
		role                        beacon.Role
		expectedAttestationDataByts []byte
		expectedError               string
	}{
		{
			"valid consensus",
			true,
			3,
			beacon.RoleAttester,
			refAttestationDataByts,
			"",
		},
		{
			"not decided",
			false,
			3,
			beacon.RoleAttester,
			refAttestationDataByts,
			"ibft did not decide, not executing role",
		},
		{
			"non supported role",
			false,
			3,
			beacon.RoleAggregator,
			refAttestationDataByts,
			"unknown role: AGGREGATOR",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := testingSSVNode(test.decided, test.signaturesCount)
			signaturesCount, inputValue, _, err := node.comeToConsensusOnInputValue(context.Background(), node.logger, []byte("id"), 0, test.role, &ethpb.DutiesResponse_Duty{
				Committee:      nil,
				CommitteeIndex: 0,
				AttesterSlot:   0,
				ProposerSlots:  nil,
				PublicKey:      nil,
				Status:         0,
				ValidatorIndex: 0,
			})
			if !test.decided {
				require.EqualError(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, 3, signaturesCount)
			require.NotNil(t, inputValue)

			byts, err := inputValue.GetAttestationData().Marshal()
			require.NoError(t, err)
			require.EqualValues(t, test.expectedAttestationDataByts, byts)
		})
	}
}

func TestPostConsensusSignatureAndAggregation(t *testing.T) {
	tests := []struct {
		name                        string
		sigs                        map[uint64][]byte
		expectedSignaturesCount     int
		expectedAttestationDataByts []byte
		expectedReconstructedSig    []byte
		expectedError               string
	}{
		{
			"valid 4/4",
			map[uint64][]byte{
				1: refAttestationSplitSigs[0],
				2: refAttestationSplitSigs[1],
				3: refAttestationSplitSigs[2],
				4: refAttestationSplitSigs[3],
			},
			4,
			refAttestationDataByts,
			refAttestationSig,
			"",
		},
		{
			"valid 3/4",
			map[uint64][]byte{
				1: refAttestationSplitSigs[0],
				2: refAttestationSplitSigs[1],
				3: refAttestationSplitSigs[2],
			},
			3,
			refAttestationDataByts,
			refAttestationSig,
			"",
		},
		{
			"invalid 3/4",
			map[uint64][]byte{
				1: refAttestationSplitSigs[0],
				2: refAttestationSplitSigs[0],
				3: refAttestationSplitSigs[2],
			},
			3,
			refAttestationDataByts,
			refAttestationSig,
			"timed out waiting for post consensus signatures, received 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := testingSSVNode(true, test.expectedSignaturesCount)
			// wait for for listeners to spin up
			time.Sleep(time.Millisecond * 100)

			// construct value
			attData := &ethpb.AttestationData{}
			require.NoError(t, attData.Unmarshal(test.expectedAttestationDataByts))
			inputValue := &proto.InputValue{
				Data: &proto.InputValue_AttestationData{
					AttestationData: attData,
				},
				SignedData: nil,
			}

			duty := &ethpb.DutiesResponse_Duty{
				Committee:      nil,
				CommitteeIndex: 0,
				AttesterSlot:   0,
				ProposerSlots:  nil,
				PublicKey:      nil,
				Status:         0,
				ValidatorIndex: 0,
			}

			// send sigs
			for index, sig := range test.sigs {
				node.network.BroadcastSignature(&proto.SignedMessage{
					Message: &proto.Message{
						Lambda: []byte("id"),
					},
					Signature: sig,
					SignerIds: []uint64{index},
				})
			}

			err := node.postConsensusDutyExecution(context.Background(), node.logger, []byte("id"), inputValue, test.expectedSignaturesCount, beacon.RoleAttester, duty)
			if len(test.expectedError) > 0 {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, test.expectedReconstructedSig, node.beacon.(*testBeacon).LastSubmittedAttestation.GetSignature())
			}
		})
	}
}
