// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	IndexPostgresDB  = insmetrics.MustTagKey("index_postgres_db")
	RecordPostgresDB = insmetrics.MustTagKey("record_postgres_db")
)

var (
	statIndexesAddedCount = stats.Int64(
		"object_indexes_added_count",
		"How many bucket have been created on a node",
		stats.UnitDimensionless,
	)
	statIndexesRemovedCount = stats.Int64(
		"object_indexes_removed_count",
		"How many bucket have been removed from a node",
		stats.UnitDimensionless,
	)
	statRecordInMemoryAddedCount = stats.Int64(
		"record_storage_added_count",
		"How many records have been saved to a in-memory storage",
		stats.UnitDimensionless,
	)
	statRecordInMemoryRemovedCount = stats.Int64(
		"record_storage_removed_count",
		"How many records have been removed from a in-memory storage",
		stats.UnitDimensionless,
	)
	SetIndexRetries = stats.Int64(
		"index_set_retries",
		"retries while SetIndex",
		stats.UnitDimensionless,
	)
	SetIndexTime = stats.Float64(
		"index_set_time",
		"time spent on SetIndex",
		stats.UnitMilliseconds,
	)
	UpdateLastKnownPulseRetries = stats.Int64(
		"index_set_retries",
		"retries while UpdateLastKnownPulse",
		stats.UnitDimensionless,
	)
	UpdateLastKnownPulseTime = stats.Float64(
		"index_set_time",
		"time spent on UpdateLastKnownPulse",
		stats.UnitMilliseconds,
	)
	ForIDTime = stats.Float64(
		"index_set_time",
		"time spent on ForID",
		stats.UnitMilliseconds,
	)
	ForPulseTime = stats.Float64(
		"index_set_time",
		"time spent on ForPulse",
		stats.UnitMilliseconds,
	)
	LastKnownForIDTime = stats.Float64(
		"index_set_time",
		"time spent on LastKnownForID",
		stats.UnitMilliseconds,
	)
	TruncateHeadIndexTime = stats.Float64(
		"index_set_time",
		"time spent on TruncateHead",
		stats.UnitMilliseconds,
	)
	TruncateHeadIndexRetries = stats.Int64(
		"index_set_retries",
		"retries while TruncateHead",
		stats.UnitDimensionless,
	)
	SetRecordTime = stats.Float64(
		"record_set_time",
		"time spent on SetRecord",
		stats.UnitMilliseconds,
	)
	SetRecordsRetries = stats.Int64(
		"record_set_retries",
		"retries while SetRecords",
		stats.UnitDimensionless,
	)
	BatchRecordsRetries = stats.Int64(
		"record_set_retries",
		"retries while BatchRecords",
		stats.UnitDimensionless,
	)
	BatchRecordTime = stats.Float64(
		"record_set_time",
		"time spent on BatchRecord",
		stats.UnitMilliseconds,
	)
	ForIDRecordTime = stats.Float64(
		"record_set_time",
		"time spent on ForID",
		stats.UnitMilliseconds,
	)
	AtPositionTime = stats.Float64(
		"record_set_time",
		"time spent on AtPosition",
		stats.UnitMilliseconds,
	)
	LastKnownPositionTime = stats.Float64(
		"record_set_time",
		"time spent on LastKnownPosition",
		stats.UnitMilliseconds,
	)
	TruncateHeadRecordRetries = stats.Int64(
		"record_set_retries",
		"retries while TruncateHeadRecord",
		stats.UnitDimensionless,
	)
	TruncateHeadRecordTime = stats.Float64(
		"record_set_time",
		"time spent on TruncateHeadRecord",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statIndexesAddedCount.Name(),
			Description: statIndexesAddedCount.Description(),
			Measure:     statIndexesAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statIndexesRemovedCount.Name(),
			Description: statIndexesRemovedCount.Description(),
			Measure:     statIndexesRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryAddedCount.Name(),
			Description: statRecordInMemoryAddedCount.Description(),
			Measure:     statRecordInMemoryAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryRemovedCount.Name(),
			Description: statRecordInMemoryRemovedCount.Description(),
			Measure:     statRecordInMemoryRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryRemovedCount.Name(),
			Description: statRecordInMemoryRemovedCount.Description(),
			Measure:     statRecordInMemoryRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        SetIndexTime.Name(),
			Description: SetIndexTime.Description(),
			Measure:     SetIndexTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetIndexRetries.Name(),
			Description: SetIndexRetries.Description(),
			Measure:     SetIndexRetries,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        UpdateLastKnownPulseTime.Name(),
			Description: UpdateLastKnownPulseTime.Description(),
			Measure:     UpdateLastKnownPulseTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        UpdateLastKnownPulseRetries.Name(),
			Description: UpdateLastKnownPulseRetries.Description(),
			Measure:     UpdateLastKnownPulseRetries,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        ForIDTime.Name(),
			Description: ForIDTime.Description(),
			Measure:     ForIDTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        ForPulseTime.Name(),
			Description: ForPulseTime.Description(),
			Measure:     ForPulseTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        LastKnownForIDTime.Name(),
			Description: LastKnownForIDTime.Description(),
			Measure:     LastKnownForIDTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadIndexTime.Name(),
			Description: TruncateHeadIndexTime.Description(),
			Measure:     TruncateHeadIndexTime,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadIndexRetries.Name(),
			Description: TruncateHeadIndexRetries.Description(),
			Measure:     TruncateHeadIndexRetries,
			TagKeys:     []tag.Key{IndexPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        SetRecordTime.Name(),
			Description: SetRecordTime.Description(),
			Measure:     SetRecordTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetRecordsRetries.Name(),
			Description: SetRecordsRetries.Description(),
			Measure:     SetRecordsRetries,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        BatchRecordTime.Name(),
			Description: BatchRecordTime.Description(),
			Measure:     BatchRecordTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        BatchRecordsRetries.Name(),
			Description: BatchRecordsRetries.Description(),
			Measure:     BatchRecordsRetries,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        ForIDRecordTime.Name(),
			Description: ForIDRecordTime.Description(),
			Measure:     ForIDRecordTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        AtPositionTime.Name(),
			Description: AtPositionTime.Description(),
			Measure:     AtPositionTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        LastKnownPositionTime.Name(),
			Description: LastKnownPositionTime.Description(),
			Measure:     LastKnownPositionTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRecordTime.Name(),
			Description: TruncateHeadRecordTime.Description(),
			Measure:     TruncateHeadRecordTime,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRecordRetries.Name(),
			Description: TruncateHeadRecordRetries.Description(),
			Measure:     TruncateHeadRecordRetries,
			TagKeys:     []tag.Key{RecordPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
	)
	if err != nil {
		panic(err)
	}
}
