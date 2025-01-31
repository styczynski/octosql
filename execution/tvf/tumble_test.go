package tvf

import (
	"testing"
	"time"

	"github.com/cube2222/octosql"
	"github.com/cube2222/octosql/execution"
)

func TestTumble_Get(t *testing.T) {
	baseTime := time.Date(2019, 9, 3, 12, 0, 0, 0, time.UTC)

	type fields struct {
		source       execution.Node
		timeField    octosql.VariableName
		windowLength execution.Expression
		offset       execution.Expression
	}
	type args struct {
		variables octosql.Variables
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    execution.RecordStream
		wantErr bool
	}{
		{
			name: "tumble with offset",
			fields: fields{
				source: execution.NewDummyNode([]*execution.Record{
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{1, baseTime},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{2, baseTime.Add(time.Second * 10)},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{3, baseTime.Add(time.Second * 8)},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{4, baseTime.Add(time.Second * 13)},
					),
				}),
				timeField:    "time",
				windowLength: execution.NewVariable(octosql.NewVariableName("window_length")),
				offset:       execution.NewVariable(octosql.NewVariableName("offset")),
			},
			args: args{
				variables: octosql.NewVariables(map[octosql.VariableName]octosql.Value{
					"window_length": octosql.MakeDuration(time.Second * 10),
					"offset":        octosql.MakeDuration(time.Second * 3),
				}),
			},
			want: execution.NewInMemoryStream([]*execution.Record{
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{1, baseTime, baseTime.Add(time.Second * -7), baseTime.Add(time.Second * 3)},
					execution.WithEventTime(baseTime.Add(time.Second*3)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{2, baseTime.Add(time.Second * 10), baseTime.Add(time.Second * 3), baseTime.Add(time.Second * 13)},
					execution.WithEventTime(baseTime.Add(time.Second*13)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{3, baseTime.Add(time.Second * 8), baseTime.Add(time.Second * 3), baseTime.Add(time.Second * 13)},
					execution.WithEventTime(baseTime.Add(time.Second*13)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{4, baseTime.Add(time.Second * 13), baseTime.Add(time.Second * 13), baseTime.Add(time.Second * 23)},
					execution.WithEventTime(baseTime.Add(time.Second*23)),
				),
			}),
			wantErr: false,
		},
		{
			name: "simple tumble",
			fields: fields{
				source: execution.NewDummyNode([]*execution.Record{
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{1, baseTime},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{2, baseTime.Add(time.Second * 10)},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{3, baseTime.Add(time.Second * 8)},
					),
					execution.NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "time"},
						[]interface{}{4, baseTime.Add(time.Second * 13)},
					),
				}),
				timeField:    "time",
				windowLength: execution.NewVariable(octosql.NewVariableName("window_length")),
				offset:       execution.NewVariable(octosql.NewVariableName("offset")),
			},
			args: args{
				variables: octosql.NewVariables(map[octosql.VariableName]octosql.Value{
					"window_length": octosql.MakeDuration(time.Second * 10),
					"offset":        octosql.MakeDuration(0),
				}),
			},
			want: execution.NewInMemoryStream([]*execution.Record{
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{1, baseTime, baseTime, baseTime.Add(time.Second * 10)},
					execution.WithEventTime(baseTime.Add(time.Second*10)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{2, baseTime.Add(time.Second * 10), baseTime.Add(time.Second * 10), baseTime.Add(time.Second * 20)},
					execution.WithEventTime(baseTime.Add(time.Second*20)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{3, baseTime.Add(time.Second * 8), baseTime, baseTime.Add(time.Second * 10)},
					execution.WithEventTime(baseTime.Add(time.Second*10)),
				),
				execution.NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "time", "window_start", "window_end"},
					[]interface{}{4, baseTime.Add(time.Second * 13), baseTime.Add(time.Second * 10), baseTime.Add(time.Second * 20)},
					execution.WithEventTime(baseTime.Add(time.Second*20)),
				),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Tumble{
				source:       tt.fields.source,
				timeField:    tt.fields.timeField,
				windowLength: tt.fields.windowLength,
				offset:       tt.fields.offset,
			}
			got, err := r.Get(tt.args.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tumble.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			eq, err := execution.AreStreamsEqual(got, tt.want)
			if err != nil {
				t.Errorf("Tumble.Get() AreStreamsEqual error = %v", err)
			}
			if !eq {
				t.Errorf("Tumble.Get() streams not equal")
			}
		})
	}
}
