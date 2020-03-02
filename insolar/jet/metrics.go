// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package jet

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	JetPostgresDB = insmetrics.MustTagKey("jet_postgres_db")
)

var (
	TruncateHeadRetries = stats.Int64(
		"jet_truncate_head_retries",
		"retries while truncating head",
		stats.UnitDimensionless,
	)
	TruncateHeadTime = stats.Float64(
		"jet_truncate_head_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
	SetTime = stats.Float64(
		"jet_set_time",
		"time spent on setting jet",
		stats.UnitMilliseconds,
	)
	SetRetries = stats.Int64(
		"jet_set_retries",
		"retries while setting jets",
		stats.UnitDimensionless,
	)
	GetTime = stats.Float64(
		"jet_get_time",
		"time spent on getting jet",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        TruncateHeadTime.Name(),
			Description: TruncateHeadTime.Description(),
			Measure:     TruncateHeadTime,
			TagKeys:     []tag.Key{JetPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetTime.Name(),
			Description: SetTime.Description(),
			Measure:     SetTime,
			TagKeys:     []tag.Key{JetPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRetries.Name(),
			Description: TruncateHeadRetries.Description(),
			Measure:     TruncateHeadRetries,
			TagKeys:     []tag.Key{JetPostgresDB},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        SetRetries.Name(),
			Description: SetRetries.Description(),
			Measure:     SetRetries,
			TagKeys:     []tag.Key{JetPostgresDB},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        GetTime.Name(),
			Description: GetTime.Description(),
			Measure:     GetTime,
			TagKeys:     []tag.Key{JetPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
	)
	if err != nil {
		panic(err)
	}
}
