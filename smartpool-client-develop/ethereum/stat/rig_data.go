package stat

import (
	"fmt"
	"github.com/SmartPool/smartpool-client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"time"
)

type PeriodRigData struct {
	MinedShare               uint64    `json:"mined_share"`
	ValidShare               uint64    `json:"valid_share"`
	TotalValidDifficulty     *big.Int  `json:"-"`
	AverageShareDifficulty   *big.Int  `json:"average_share_difficulty"`
	RejectedShare            uint64    `json:"rejected_share"`
	TotalHashrate            *big.Int  `json:"-"`
	NoHashrateSubmission     uint64    `json:"-"`
	AverageReportedHashrate  *big.Int  `json:"reported_hashrate"`
	AverageEffectiveHashrate *big.Int  `json:"effective_hashrate"`
	BlockFound               uint64    `json:"block_found"`
	TimePeriod               uint64    `json:"time_period"`
	StartTime                time.Time `json:"start_time"`
}

func NewPeriodRigData(timePeriod uint64) *PeriodRigData {
	return &PeriodRigData{
		TotalHashrate:            big.NewInt(0),
		TotalValidDifficulty:     big.NewInt(0),
		AverageShareDifficulty:   big.NewInt(0),
		AverageReportedHashrate:  big.NewInt(0),
		AverageEffectiveHashrate: big.NewInt(0),
		TimePeriod:               timePeriod,
	}
}

func (prd *PeriodRigData) updateAvgHashrate(t time.Time) {
	if prd.NoHashrateSubmission > 0 {
		prd.AverageReportedHashrate.Div(
			prd.TotalHashrate,
			big.NewInt(int64(prd.NoHashrateSubmission)),
		)
	}
}

func (prd *PeriodRigData) updateAvgEffHashrate(t time.Time) {
	prd.AverageEffectiveHashrate.Div(
		prd.TotalValidDifficulty,
		big.NewInt(BaseTimePeriod),
	)
}

func (prd *PeriodRigData) updateAvgShareDifficulty(t time.Time) {
	if prd.ValidShare > 0 {
		prd.AverageShareDifficulty.Div(
			prd.TotalValidDifficulty,
			big.NewInt(int64(prd.ValidShare)),
		)
	}
}

type OverallRigData struct {
	LastMinedShare           time.Time `json:"last_mined_share"`
	LastValidShare           time.Time `json:"last_valid_share"`
	LastRejectedShare        time.Time `json:"last_rejected_share"`
	LastBlock                time.Time `json:"last_block"`
	MinedShare               uint64    `json:"total_submitted_share"`
	ValidShare               uint64    `json:"total_accepted_share"`
	TotalValidDifficulty     *big.Int  `json:"total_accepted_difficulty"`
	AverageShareDifficulty   *big.Int  `json:"average_share_difficulty"`
	RejectedShare            uint64    `json:"total_rejected_share"`
	TotalHashrate            *big.Int  `json:"total_hashrate"`
	NoHashrateSubmission     uint64    `json:"no_hashrate_submission"`
	AverageReportedHashrate  *big.Int  `json:"reported_hashrate"`
	AverageEffectiveHashrate *big.Int  `json:"effective_hashrate"`
	BlockFound               uint64    `json:"total_block_found"`
	StartTime                time.Time `json:"start_time"`
}

type RigData struct {
	RigID string
	Datas map[uint64]*PeriodRigData
	*OverallRigData
}

func NewRigData(rigID string) *RigData {
	return &RigData{
		RigID: rigID,
		Datas: map[uint64]*PeriodRigData{},
		OverallRigData: &OverallRigData{
			TotalHashrate:            big.NewInt(0),
			TotalValidDifficulty:     big.NewInt(0),
			AverageShareDifficulty:   big.NewInt(0),
			AverageReportedHashrate:  big.NewInt(0),
			AverageEffectiveHashrate: big.NewInt(0),
		},
	}
}

func (rd *RigData) TruncateData(storage smartpool.PersistentStorage) error {
	curPeriod := TimeToPeriod(time.Now())
	var err error
	for period, rigData := range rd.Datas {
		if int64(curPeriod-period) > LongWindow/BaseTimePeriod {
			if err = storage.Persist(rigData, fmt.Sprintf("rig-%s-data-%d", rd.RigID, period)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rd *RigData) getData(t time.Time) *PeriodRigData {
	timePeriod := TimeToPeriod(t)
	data := rd.Datas[timePeriod]
	if data == nil {
		data = NewPeriodRigData(timePeriod)
		rd.Datas[timePeriod] = data
	}
	return data
}

func (rd *RigData) AddShare(status string, share smartpool.Share, t time.Time) {
	if rd.StartTime.IsZero() {
		rd.StartTime = t
	}
	curPeriodData := rd.getData(t)
	if curPeriodData.StartTime.IsZero() {
		curPeriodData.StartTime = t
	}
	if status == "submitted" {
		rd.LastMinedShare = t
		rd.MinedShare++
		curPeriodData.MinedShare++
	} else if status == "accepted" {
		rd.LastValidShare = t
		rd.ValidShare++
		rd.TotalValidDifficulty.Add(rd.TotalValidDifficulty, share.ShareDifficulty())
		rd.updateAvgShareDifficulty(t)
		rd.updateAvgEffHashrate(t)
		curPeriodData.ValidShare++
		curPeriodData.TotalValidDifficulty.Add(curPeriodData.TotalValidDifficulty, share.ShareDifficulty())
		curPeriodData.updateAvgShareDifficulty(t)
		curPeriodData.updateAvgEffHashrate(t)
	} else if status == "rejected" {
		rd.LastRejectedShare = t
		rd.RejectedShare++
		curPeriodData.RejectedShare++
	} else if status == "fullsolution" {
		rd.LastBlock = t
		rd.LastValidShare = t
		rd.ValidShare++
		rd.TotalValidDifficulty.Add(rd.TotalValidDifficulty, share.ShareDifficulty())
		rd.updateAvgShareDifficulty(t)
		rd.updateAvgEffHashrate(t)
		rd.BlockFound++
		curPeriodData.ValidShare++
		curPeriodData.TotalValidDifficulty.Add(curPeriodData.TotalValidDifficulty, share.ShareDifficulty())
		curPeriodData.updateAvgShareDifficulty(t)
		curPeriodData.updateAvgEffHashrate(t)
		curPeriodData.BlockFound++
	}
}

func (rd *RigData) PeriodReportedHashrate(t time.Time) *big.Int {
	curPeriodData := rd.getData(t)
	return curPeriodData.AverageReportedHashrate
}

func (rd *RigData) PeriodEffectiveHashrate(t time.Time) *big.Int {
	curPeriodData := rd.getData(t)
	return curPeriodData.AverageEffectiveHashrate
}

func (rd *RigData) AddHashrate(hashrate hexutil.Uint64, id common.Hash, t time.Time) {
	if rd.StartTime.IsZero() {
		rd.StartTime = t
	}
	curPeriodData := rd.getData(t)
	if curPeriodData.StartTime.IsZero() {
		curPeriodData.StartTime = t
	}
	rd.TotalHashrate.Add(rd.TotalHashrate, big.NewInt(int64(hashrate)))
	rd.NoHashrateSubmission++
	rd.updateAvgHashrate(t)
	// fmt.Printf("Updated total hashrate: %d\n", rd.TotalHashrate.Uint64())
	// fmt.Printf("Updated reported hashrate: %d\n", rd.AverageReportedHashrate.Uint64())
	curPeriodData.TotalHashrate.Add(curPeriodData.TotalHashrate, big.NewInt(int64(hashrate)))
	curPeriodData.NoHashrateSubmission++
	curPeriodData.updateAvgHashrate(t)
}

func (rd *RigData) updateAvgHashrate(t time.Time) {
	if rd.NoHashrateSubmission > 0 {
		rd.AverageReportedHashrate.Div(
			rd.TotalHashrate,
			big.NewInt(int64(rd.NoHashrateSubmission)),
		)
	}
}

func (rd *RigData) updateAvgEffHashrate(t time.Time) {
	rd.AverageEffectiveHashrate.Div(
		rd.TotalValidDifficulty,
		big.NewInt(BaseTimePeriod),
	)
}

func (rd *RigData) updateAvgShareDifficulty(t time.Time) {
	if rd.ValidShare > 0 {
		rd.AverageShareDifficulty.Div(
			rd.TotalValidDifficulty,
			big.NewInt(int64(rd.ValidShare)),
		)
	}
}
