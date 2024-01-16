package etcd

import (
	"etcd-sync/config"
	"log"
	"testing"
)

func setupEtcd() *Etcd {
	sourceCfg := config.Config{
		Host:     "http://127.0.0.1:12379",
		Username: "",
		Password: "",
	}

	targetCfg := config.Config{
		Host:     "http://127.0.0.1:22379",
		Username: "",
		Password: "",
	}

	return NewEtcd(sourceCfg, targetCfg)
}

func TestGet(t *testing.T) {
	etcd := setupEtcd()

	etcd.Get("/message")
}

func TestSync(t *testing.T) {
	etcd := setupEtcd()

	etcd.Get("/message")
	err := etcd.Sync("/")
	if err != nil {
		log.Fatal("err: ", err)
	}
	etcd.Get("/message")

}

func TestExport(t *testing.T) {
	etcd := setupEtcd()

	etcd.Export(SourceType, "etcd_export")
	etcd.Export(TargetType, "etcd_export")
}

func TestImport(t *testing.T) {
	etcd := setupEtcd()
	//etcd.ImportData(SourceType, "etcd_export_source.json")
	etcd.ImportData(TargetType, "etcd_export_target.json")
}
