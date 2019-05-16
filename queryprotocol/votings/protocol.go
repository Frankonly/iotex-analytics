// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package votings

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-analytics/indexprotocol"
	"github.com/iotexproject/iotex-analytics/indexprotocol/votings"
	"github.com/iotexproject/iotex-analytics/indexservice"
	s "github.com/iotexproject/iotex-analytics/sql"
)

// Protocol defines the protocol of querying tables
type Protocol struct {
	indexer *indexservice.Indexer
}

// VotingInfo defines voting infos
type VotingInfo struct {
	VoterAddress  string
	WeightedVotes string
}

// NewProtocol creates a new protocol
func NewProtocol(idx *indexservice.Indexer) *Protocol {
	return &Protocol{indexer: idx}
}

// GetVotingInformation gets voting infos
func (p *Protocol) GetVotingInformation(epochNum int, delegateName string) (votingInfos []*VotingInfo, err error) {
	if _, ok := p.indexer.Registry.Find(votings.ProtocolID); !ok {
		err = errors.New("votings protocol is unregistered")
		return
	}
	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf("SELECT voter_address,weighted_votes FROM %s WHERE epoch_number = ? and candidate_name = ?",
		votings.VotingHistoryTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare get query")
		return
	}

	rows, err := stmt.Query(epochNum, delegateName)
	if err != nil {
		err = errors.Wrap(err, "failed to execute get query")
		return
	}

	var votingHistory VotingInfo
	parsedRows, err := s.ParseSQLRows(rows, &votingHistory)
	if err != nil {
		err = errors.Wrap(err, "failed to parse results")
		return
	}

	if len(parsedRows) == 0 {
		err = indexprotocol.ErrNotExist
		return
	}

	for _, parsedRow := range parsedRows {
		voting := parsedRow.(*VotingInfo)
		votingInfos = append(votingInfos, voting)
	}
	return
}
