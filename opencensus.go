package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	//"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	//"go.opencensus.io/resource"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

func openCensusTest(ctx context.Context, projectID string) error {
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: "demo-prefix",
		// ReportingInterval sets the frequency of reporting metrics
		// to stackdriver backend.
		ReportingInterval: 60 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
		return err
	}
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	// Start the metrics exporter
	err = sd.StartMetricsExporter()
	if err != nil {
		log.Fatalf("Failed to start metrics exporter: %v", err)
		return err
	}
	defer sd.StopMetricsExporter()

	reader := metricexport.NewReader()

	// Encounters the number of non EOF(end-of-file) errors.
	var mErrors = stats.Int64("opencensus-example/errors", "The number of errors encountered", stats.UnitDimensionless)
	var errorCountView = &view.View{
		Name:        "opencensus-demo/oc-errors",
		Measure:     mErrors,
		Description: "The number of errors encountered",
		Aggregation: view.Count(),
	}
	// Register the views
	if err := view.Register(errorCountView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
	stats.Record(ctx, mErrors.M(1))

	reader.ReadAndExport(sd)
	//err = sd.ExportMetrics(ctx, []*metricdata.Metric{
	//	&metricdata.Metric{
	//		Descriptor: metricdata.Descriptor{
	//			Name:        "opencensus-descriptor",
	//			Description: "",
	//			Unit:        "",
	//			Type:        0,
	//			LabelKeys:   nil,
	//		},
	//		Resource:   &resource.Resource{
	//			Type:   "",
	//			Labels:  map[string]string{
	//				"project_id": projectID,
	//			},
	//		},
	//		TimeSeries: []*metricdata.TimeSeries{
	//			&metricdata.TimeSeries{
	//				LabelValues: []metricdata.LabelValue{
	//					metricdata.LabelValue{
	//						Value:   "project_id:" + projectID,
	//						Present: true,
	//					},
	//				},
	//				Points:      []metricdata.Point{},
	//				StartTime:   time.Time{},
	//			},
	//		},
	//	},
	//})

	if err != nil {
		log.Fatalf("Failed to export metrics from OpenCensus to Stackdriver: %v", err)
		return err
	}

	fmt.Printf("Done writing time series data via OpenCensus.\n")
	return nil
}