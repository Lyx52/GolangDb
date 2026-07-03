package backing

import (
	"github.com/Lyx52/GolangDb/models"
	"github.com/Lyx52/GolangDb/signals"
)

type RowWriteChannel chan<- models.DataRow

type DataStore struct {
	Values        []*models.DataRow
	AppendChannel chan models.DataRow
}

func NewDataStore() *DataStore {
	return &DataStore{
		Values:        []*models.DataRow{},
		AppendChannel: make(chan models.DataRow),
	}
}

func (store *DataStore) AddRow(row *models.DataRow) {
	store.Values = append(store.Values, row)
}

type IndexedRow struct {
	Index int
	Row   *models.DataRow
}

func (store *DataStore) FullRowScan() <-chan IndexedRow {
	channel := make(chan IndexedRow)
	go func() {
		for i, row := range store.Values {
			channel <- IndexedRow{i, row}
		}
		close(channel)
	}()
	return channel
}

func (store *DataStore) HandleDataStoreWrites(canceled signals.CancelSignal) {
	for {
		select {
		case <-canceled:
			return
		case row, ok := <-store.AppendChannel:
			if !ok {
				return
			}

			store.AddRow(&row)
		}
	}
}
