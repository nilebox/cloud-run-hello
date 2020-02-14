package main

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func stackdriverTest(ctx context.Context, projectID string) error {
	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}

	// Prepares an individual data point
	dataPoint := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_DoubleValue{
				DoubleValue: 123.45,
			},
		},
	}

	// Writes time series data.
	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/stores/daily_sales",
					Labels: map[string]string{
						"store_id": "Pittsburg",
					},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": projectID,
					},
				},
				Points: []*monitoringpb.Point{
					dataPoint,
				},
			},
		},
	}); err != nil {
		log.Fatalf("Failed to write time series data: %v", err)
		return err
	}

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
		return err
	}

	fmt.Printf("Done writing time series data directly.\n")
	return nil
}
