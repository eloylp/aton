// +build e2e

package components_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/eloylp/aton/pkg/test/helper"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/ctl"
	ctlConfig "github.com/eloylp/aton/components/ctl/config"
	"github.com/eloylp/aton/components/ctl/www"
	"github.com/eloylp/aton/components/node"
	nodeConfig "github.com/eloylp/aton/components/node/config"
)

const (
	ModelsDir = "../models"
	imagesDir = "../samples/images"
)

var (
	faceBona1 = filepath.Join(imagesDir, "bona.jpg")
	faceBona2 = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3 = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4 = filepath.Join(imagesDir, "bona4.jpg")
	group     = filepath.Join(imagesDir, "bonaAndLuda.jpg")
)

func TestRun(t *testing.T) {
	t.Skip() // TODO make this green
	ctlLog := bytes.NewBuffer(nil)
	ctlAddr := "localhost:9090"
	cplane, err := ctl.New(
		ctlConfig.WithListenAddress(ctlAddr),
		ctlConfig.WithLogOutput(ctlLog),
	)
	assert.NoError(t, err)

	nodeLog := bytes.NewBuffer(nil)
	node1Addr := "localhost:9091"
	node1, err := node.New(
		nodeConfig.WithListenAddress(node1Addr),
		nodeConfig.WithMetricsAddress("localhost:9092"),
		nodeConfig.WithModelDir(ModelsDir),
		nodeConfig.WithLogOutput(nodeLog),
	)
	node2Addr := "localhost:9093"
	node2, err := node.New(
		nodeConfig.WithListenAddress(node2Addr),
		nodeConfig.WithMetricsAddress("localhost:9094"),
		nodeConfig.WithModelDir(ModelsDir),
		nodeConfig.WithLogOutput(nodeLog),
	)
	assert.NoError(t, err)
	go cplane.Start()
	defer cplane.Shutdown()

	go node1.Start()
	defer node1.Shutdown()

	go node2.Start()
	defer node2.Shutdown()

	images1 := []string{faceBona1, faceBona2}
	videoSource1 := helper.ReplayedVideoStream(t, images1, "/", 10, 25)
	defer videoSource1.Close()

	images2 := []string{faceBona3, faceBona4}
	videoSource2 := helper.ReplayedVideoStream(t, images2, "/", 10, 25)
	defer videoSource2.Close()

	helper.TryConnectTo(t, ctlAddr, time.Second)

	addNode(t, ctlAddr, node1Addr)

	addTarget(t, ctlAddr, "images-1-camera", videoSource1.URL)
	addTarget(t, ctlAddr, "images-2-camera", videoSource2.URL)

	loadCategories(t, ctlAddr, []string{"bona", "luda"}, helper.ReadFile(t, group))

	assert.Contains(t, ctlLog.String(), "recognized")

}

func addNode(t *testing.T, ctlAddr string, nodeAddr string) {
	t.Helper()
	r := &www.AddNodeRequest{
		Addr: nodeAddr,
	}
	sendDataAndExpect(t, "http://"+ctlAddr+"/nodes", r, http.StatusOK)
}

func sendDataAndExpect(t *testing.T, addr string, data interface{}, expectStatus int) {
	t.Helper()
	requestData := bytes.NewBuffer(nil)
	if err := json.NewEncoder(requestData).Encode(data); err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(addr, "application/json", bytes.NewReader(requestData.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectStatus {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("expected status code when adding targetAddr is %d, found %d: %s", expectStatus, resp.StatusCode, data)
	}
}

func loadCategories(t *testing.T, ctlAddr string, categoriesNames []string, categoriesImage []byte) {
	t.Helper()
	r := &www.LoadCategoriesRequest{
		Categories: categoriesNames,
		Image:      categoriesImage,
	}
	sendDataAndExpect(t, "http://"+ctlAddr+"/categories", r, http.StatusOK)
}

func addTarget(t *testing.T, ctlAddr string, targetUuid, targetAddr string) {
	t.Helper()
	r := &www.AddTargetRequest{
		UUID:       targetUuid,
		TargetAddr: targetAddr,
	}
	sendDataAndExpect(t, "http://"+ctlAddr+"/targets", r, http.StatusCreated)
}
