package agent

import (
	memstorage "github.com/chikiryau3/garbage-collector/internal/memStorage"
	"github.com/chikiryau3/garbage-collector/internal/mocks"
	"github.com/golang/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

//go:generate mockgen -source=../metricsCollector/metrics-collector.go -destination=../mocks/metrics-collector_mock.go -package=mocks -mock_names=MetricsCollector=MetricsCollectorMock
//go:generate mockgen -source=../clients/garbage-collector/client.go -destination=../mocks/garbage-collector_mock.go -package=mocks -mock_names=Client=GarbageCollectorMock

// покрыл функции, содержащие основную логику (как тестить то что в горутинах пока не знаю)
// тесты пока только на обычные случаи, без корнер-кейсов с ошибками (не успеваю ес чесн)

func Test_agent_pollMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectorMock := mocks.NewMetricsCollectorMock(ctrl)
	collectionServiceClientMock := mocks.NewGarbageCollectorMock(ctrl)
	pollInterval := time.Second * 2
	reportInterval := time.Second * 10

	tests := []struct {
		name           string
		wantErr        bool
		expectationsFn func()
	}{
		{
			name:    `success`,
			wantErr: false,
			expectationsFn: func() {
				// много gauge из runtime
				collectorMock.EXPECT().SetGauge(gomock.Any(), gomock.Any()).AnyTimes()
				// 2 count
				collectorMock.EXPECT().SetCount(gomock.Any(), gomock.Any()).Times(2)
			},
		},
	}

	a := &agent{
		collector:               collectorMock,
		collectionServiceClient: collectionServiceClientMock,
		config: Config{
			`/`,
			pollInterval,
			reportInterval,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectationsFn()
			if err := a.pollMetrics(); (err != nil) != tt.wantErr {
				t.Errorf("pollMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agent_sendReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectorMock := mocks.NewMetricsCollectorMock(ctrl)
	collectionServiceClientMock := mocks.NewGarbageCollectorMock(ctrl)
	pollInterval := time.Second * 2
	reportInterval := time.Second * 10

	storageData := &memstorage.StorageData{
		"gauge:someGaugeMetric":  rand.Float64(),
		"gauge:someGaugeMetric2": rand.Float64(),
		"count:someCountMetric":  rand.Int63(),
		"count:someCountMetric2": rand.Int63(),
	}

	tests := []struct {
		name           string
		wantErr        bool
		expectationsFn func()
	}{
		{
			name:    `success`,
			wantErr: false,
			expectationsFn: func() {
				collectorMock.EXPECT().ReadStorage().Return(storageData, nil).Times(1)
				collectionServiceClientMock.EXPECT().SendGauge(gomock.Any(), gomock.Any()).Times(2)
				collectionServiceClientMock.EXPECT().SendCounter(gomock.Any(), gomock.Any()).Times(2)
			},
		},
	}

	a := &agent{
		collector:               collectorMock,
		collectionServiceClient: collectionServiceClientMock,
		config: Config{
			`/`,
			pollInterval,
			reportInterval,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectationsFn()
			if err := a.sendReport(); (err != nil) != tt.wantErr {
				t.Errorf("pollMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
