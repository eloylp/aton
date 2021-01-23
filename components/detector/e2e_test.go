// +build e2e detector

package detector_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eloylp/aton/pkg/test/helper"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/components/detector"
	"github.com/eloylp/aton/components/detector/config"
	"github.com/eloylp/aton/components/proto"
)

func TestStartStopSequence(t *testing.T) {
	logOutput := bytes.NewBuffer(nil)
	d, err := detector.New(
		config.WithListenAddress("0.0.0.0:10002"),
		config.WithMetricsAddress("0.0.0.0:10003"),
		config.WithLogOutput(logOutput),
		config.WithLogFormat("text"),
		config.WithModelDir("../../models"),
	)
	assert.NoError(t, err)

	go d.Start()
	helper.TryConnectTo(t, "127.0.0.1:10002", time.Second)
	helper.TryConnectTo(t, "127.0.0.1:10003", time.Second)
	d.Shutdown()

	logO := logOutput.String()
	assert.Contains(t, logO, "starting detector service at 0.0.0.0:10002")
	assert.Contains(t, logO, "starting detector metrics at 0.0.0.0:10003")
	assert.Contains(t, logO, "gracefully shutdown started.")
	assert.Contains(t, logO, "stopped detector service at 0.0.0.0:10002")
	assert.Contains(t, logO, "stopped detector metrics at 0.0.0.0:10003")
	assert.NotContains(t, logO, "level=error")

	t.Log(logO)
}

func TestMatchingCapturingRound(t *testing.T) {
	video := helper.ReplayedVideoStream(t, []string{faceBona1, faceBona2}, "/", 100)
	defer video.Close()

	logOutput := bytes.NewBuffer(nil)
	d, err := detector.New(
		config.WithUUID("UUID"),
		config.WithListenAddress("0.0.0.0:10002"),
		config.WithMetricsAddress("0.0.0.0:10003"),
		config.WithLogOutput(logOutput),
		config.WithLogFormat("text"),
		config.WithModelDir("../../models"),
	)
	assert.NoError(t, err)

	go d.Start()
	defer d.Shutdown()
	helper.TryConnectTo(t, "127.0.0.1:10002", time.Second)
	helper.TryConnectTo(t, "127.0.0.1:10003", time.Second)

	con, err := grpc.Dial("127.0.0.1:10002", grpc.WithInsecure(), grpc.WithBlock())
	assert.NoError(t, err)
	defer con.Close()

	assert.NoError(t, err)
	client := proto.NewDetectorClient(con)
	_, err = client.LoadCategories(context.Background(), &proto.LoadCategoriesRequest{
		Categories: []string{"bona"},
		Image:      helper.ReadFile(t, faceBona3),
	})
	assert.NoError(t, err)
	_, err = client.AddCapturer(context.Background(), &proto.AddCapturerRequest{CapturerUuid: "UUID", CapturerUrl: video.URL})
	assert.NoError(t, err)
	clientProcess, err := client.ProcessResults(context.Background(), &empty.Empty{})
	assert.NoError(t, err)

	recv, err := clientProcess.Recv()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recv.Recognized))
	assert.Equal(t, int32(1), recv.TotalEntities)
	assert.Equal(t, "UUID", recv.DetectorUuid)
	now := time.Now().Unix()
	assert.InDelta(t, now, recv.CapturedAt.AsTime().Unix(), 5)
	assert.InDelta(t, now, recv.RecognizedAt.AsTime().Unix(), 5)
	assert.True(t, recv.CapturedAt.AsTime().Before(recv.RecognizedAt.AsTime()))

	logO := logOutput.String()
	assert.Contains(t, logO, "capturerHandler: added target with UUID: UUID")
	assert.Contains(t, logO, "capturerHandler: starting target with UUID: UUID")
	assert.NotContains(t, logO, "level=error")

	resp, err := http.Get("http://127.0.0.1:10003/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	metricsData, err := ioutil.ReadAll(resp.Body)
	metricsO := string(metricsData)
	assert.Contains(t, metricsO, `aton_detector_capturer_received_frames_total{capturer_url="`+video.URL+`",capturer_uuid="UUID",uuid="UUID"}`)
	assert.Contains(t, metricsO, `aton_detector_entities_total{uuid="UUID"} 1`)
	assert.Contains(t, metricsO, `aton_detector_unrecognized_entities_total{uuid="UUID"} 0`)
	assert.Contains(t, metricsO, `aton_detector_processed_frames_total{uuid="UUID"} 2`)
	assert.Contains(t, metricsO, `grpc_server_msg_sent_total{grpc_method="AddCapturer",grpc_service="proto.Detector",grpc_type="unary"} 1`)
}

func TestNonMatchingCapturingRound(t *testing.T) {
	video := helper.ReplayedVideoStream(t, []string{faceBona1, faceBona2}, "/", 100)
	defer video.Close()

	logOutput := bytes.NewBuffer(nil)
	d, err := detector.New(
		config.WithUUID("UUID"),
		config.WithListenAddress("0.0.0.0:10002"),
		config.WithMetricsAddress("0.0.0.0:10003"),
		config.WithLogOutput(logOutput),
		config.WithLogFormat("text"),
		config.WithModelDir("../../models"),
	)
	assert.NoError(t, err)

	go d.Start()
	defer d.Shutdown()
	helper.TryConnectTo(t, "127.0.0.1:10002", time.Second)
	helper.TryConnectTo(t, "127.0.0.1:10003", time.Second)

	con, err := grpc.Dial("127.0.0.1:10002", grpc.WithInsecure(), grpc.WithBlock())
	assert.NoError(t, err)
	defer con.Close()

	assert.NoError(t, err)
	client := proto.NewDetectorClient(con)
	_, err = client.AddCapturer(context.Background(), &proto.AddCapturerRequest{CapturerUuid: "UUID", CapturerUrl: video.URL})
	assert.NoError(t, err)
	clientProcess, err := client.ProcessResults(context.Background(), &empty.Empty{})
	assert.NoError(t, err)

	recv, err := clientProcess.Recv()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recv.Recognized))
	assert.Equal(t, int32(1), recv.TotalEntities)
	assert.Equal(t, "UUID", recv.DetectorUuid)
	now := time.Now().Unix()
	assert.InDelta(t, now, recv.CapturedAt.AsTime().Unix(), 5)
	assert.InDelta(t, now, recv.RecognizedAt.AsTime().Unix(), 5)
	assert.True(t, recv.CapturedAt.AsTime().Before(recv.RecognizedAt.AsTime()))

	logO := logOutput.String()
	assert.NotContains(t, logO, "level=error")

	resp, err := http.Get("http://127.0.0.1:10003/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	metricsData, err := ioutil.ReadAll(resp.Body)
	metricsO := string(metricsData)
	assert.Contains(t, metricsO, `aton_detector_capturer_received_frames_total{capturer_url="`+video.URL+`",capturer_uuid="UUID",uuid="UUID"}`)
	assert.Contains(t, metricsO, `aton_detector_entities_total{uuid="UUID"} 1`)
	assert.Contains(t, metricsO, `aton_detector_processed_frames_total{uuid="UUID"} 2`)
	assert.Contains(t, metricsO, `aton_detector_unrecognized_entities_total{uuid="UUID"} 1`)
	assert.Contains(t, metricsO, `grpc_server_msg_sent_total{grpc_method="AddCapturer",grpc_service="proto.Detector",grpc_type="unary"} 1`)
}
