package etcd

import (
	"context"
	"encoding/json"
	"etcd-sync/config"
	"fmt"
	"go.etcd.io/etcd/client/v2"
	"io"
	"log"
	"os"
)

type Etcd struct {
	source     client.Client
	sourceKApi client.KeysAPI

	target     client.Client
	targetKApi client.KeysAPI
}

const SourceType = "source"
const TargetType = "target"

func NewEtcd(sourceCfg, targetCfg config.Config) *Etcd {
	source, err := client.New(client.Config{
		Endpoints: []string{sourceCfg.Host},
		Transport: client.DefaultTransport,
		Username:  sourceCfg.Username,
		Password:  sourceCfg.Password,
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	target, err := client.New(client.Config{
		Endpoints: []string{targetCfg.Host},
		Transport: client.DefaultTransport,
		Username:  targetCfg.Username,
		Password:  targetCfg.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &Etcd{
		source:     source,
		sourceKApi: client.NewKeysAPI(source),
		target:     target,
		targetKApi: client.NewKeysAPI(target),
	}
}

func (e *Etcd) Get(key string) {
	var source, target string

	sourceResp, err := e.sourceKApi.Get(context.Background(), key, &client.GetOptions{Recursive: false})
	if err != nil {
		log.Println("source: ", err)
	} else {
		source = sourceResp.Node.Value
	}

	targetResp, err := e.targetKApi.Get(context.Background(), key, &client.GetOptions{Recursive: false})
	if err != nil {
		log.Println("target: ", err)
	} else {
		target = targetResp.Node.Value
	}

	fmt.Printf("source: %s : %s \n", key, source)
	fmt.Printf("target: %s : %s \n", key, target)
}

func (e *Etcd) Sync(key string) (err error) {
	resp, err := e.sourceKApi.Get(context.Background(), key, &client.GetOptions{Recursive: true})
	err = e.syncNode(resp.Node, "", client.NewKeysAPI(e.target))
	return
}

func (e *Etcd) syncNode(node *client.Node, path string, targetKapi client.KeysAPI) (err error) {
	if !node.Dir {
		fmt.Println(path)
		_, err = targetKapi.Set(context.Background(), path, node.Value, nil)
	} else {
		for _, child := range node.Nodes {
			err = e.syncNode(child, child.Key, targetKapi)
		}
	}
	return
}

func (e *Etcd) Export(fType, filename string) {

	var resp *client.Response
	var err error

	switch fType {
	case SourceType:
		resp, err = e.sourceKApi.Get(context.Background(), "/", &client.GetOptions{Recursive: true})
		break
	case TargetType:
		resp, err = e.targetKApi.Get(context.Background(), "/", &client.GetOptions{Recursive: true})
		break
	}

	if err != nil {
		log.Fatal(err)
	}

	data := e.extractData(resp.Node)

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	e.writeFile(filename+"_"+fType+".json", dataBytes)

	fmt.Printf("Data exported to %s_%s.json\n", filename, fType)
}

func (e *Etcd) ExportAll(filename string) {
	e.Export(SourceType, filename)
	e.Export(TargetType, filename)

	fmt.Printf("Data exported to %s_source.json and %s_target.json \n", filename, filename)

}

func (e *Etcd) extractData(node *client.Node) map[string]string {
	data := make(map[string]string)

	if !node.Dir {
		data[node.Key] = node.Value
	} else {
		for _, child := range node.Nodes {
			childData := e.extractData(child)
			for k, v := range childData {
				data[k] = v
			}
		}
	}

	return data
}

func (e *Etcd) writeFile(filename string, dataBytes []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(dataBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Etcd) openFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	dataBytes := make([]byte, fileInfo.Size())

	_, err = io.ReadFull(file, dataBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return dataBytes
}

func (e *Etcd) ImportData(fType, filename string) {
	dataBytes := e.openFile(filename)

	data := make(map[string]string)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range data {
		switch fType {
		case SourceType:
			_, err = e.sourceKApi.Set(context.Background(), key, value, nil)
			break
		case TargetType:
			_, err = e.targetKApi.Set(context.Background(), key, value, nil)
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Data imported from %s \n", filename)
}
