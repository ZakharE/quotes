package daemons

import (
	"context"
	"fmt"
	"plata_card_quotes/internal/quotes/models"
	"reflect"
	"testing"
	"time"
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
			name: "Split zero length",
			args: args{make([]models.TaskDTO, 0)},
			want: make(map[models.CurrencyPair][]int64),
		},
		{
			name: "Split by currency",
			args: args{[]models.TaskDTO{
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEur,
						Counter: models.CurrencyUsd,
					},
					TaskID: 1,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEur,
						Counter: models.CurrencyUsd,
					},
					TaskID: 2,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEur,
						Counter: models.CurrencyUsd,
					},
					TaskID: 3,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyEur,
						Counter: models.CurrencyMxn,
					},
					TaskID: 4,
				},
				{
					CurrencyPair: models.CurrencyPair{
						Base:    models.CurrencyMxn,
						Counter: models.CurrencyEur,
					},
					TaskID: 5,
				},
			}},
			want: map[models.CurrencyPair][]int64{
				models.CurrencyPair{
					Base:    models.CurrencyEur,
					Counter: models.CurrencyUsd,
				}: {1, 2, 3},
				models.CurrencyPair{
					Base:    models.CurrencyEur,
					Counter: models.CurrencyMxn,
				}: {4},
				models.CurrencyPair{
					Base:    models.CurrencyMxn,
					Counter: models.CurrencyEur,
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func(ctx context.Context) {
		time.Sleep(10 * time.Second)
		fmt.Print("LOL")
	}(ctx)
	select {
	case <-ctx.Done():
		fmt.Print(ctx.Err())
	}
	fmt.Print("NON BLOCK")

}
