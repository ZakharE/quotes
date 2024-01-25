package daemons

import (
	"plata_card_quotes/internal/quotes/models"
	"reflect"
	"testing"
)

func Test_splitByCurrency(t *testing.T) {
	type args struct {
		tasks []models.TaskDTO
	}
	tests := []struct {
		name string
		args args
		want map[models.CurrencyPair][]int64
	}{
		{
			name: "Splitting zero length slice returns empty map",
			args: args{make([]models.TaskDTO, 0)},
			want: make(map[models.CurrencyPair][]int64),
		},
		{
			name: "Splits by currency",
			args: args{[]models.TaskDTO{
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEUR,
						Counter: models.CurrencyUSD,
					},
					TaskID: 1,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEUR,
						Counter: models.CurrencyUSD,
					},
					TaskID: 2,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEUR,
						Counter: models.CurrencyUSD,
					},
					TaskID: 3,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEUR,
						Counter: models.CurrencyMXN,
					},
					TaskID: 4,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyMXN,
						Counter: models.CurrencyEUR,
					},
					TaskID: 5,
				},
			}},
			want: map[models.CurrencyPair][]int64{
				models.CurrencyPair{
					Base:    models.CurrencyEUR,
					Counter: models.CurrencyUSD,
				}: {1, 2, 3},
				models.CurrencyPair{
					Base:    models.CurrencyEUR,
					Counter: models.CurrencyMXN,
				}: {4},
				models.CurrencyPair{
					Base:    models.CurrencyMXN,
					Counter: models.CurrencyEUR,
				}: {5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByCurrency(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}
