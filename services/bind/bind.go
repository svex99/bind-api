package bind

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/svex99/bind-api/pkg/setting"
	"github.com/svex99/bind-api/schemas"
	"github.com/svex99/bind-api/services/bind/parser"
)

type BindService struct {
	ctx           context.Context
	Mutex         *sync.Mutex
	DockerCli     *client.Client
	ContainerId   string
	ZonesFilePath string
	BindConf      *parser.BindConf
	Zones         map[string]*parser.ZoneConf
}

var Service = &BindService{}

func (bs *BindService) Init() {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	Service.ctx = context.Background()
	Service.Mutex = &sync.Mutex{}
	Service.DockerCli = cli
	Service.ContainerId = setting.Bind.ContainerId
	Service.ZonesFilePath = setting.Bind.ConfPath + "named.conf.local"

	Service.Load()
}

func (bs *BindService) Load() {
	var err error

	bs.BindConf, err = bs.parseBindConf(Service.ZonesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">>> Loaded %d zone(s) from %s\n", len(Service.BindConf.Zones), Service.ZonesFilePath)

	bs.Zones = make(map[string]*parser.ZoneConf)

	fmt.Println(">>> Loading BIND9 zone files")
	for _, zone := range Service.BindConf.Zones {
		filename := zone.File[strings.LastIndex(zone.File, "/")+1:]

		zConf, err := Service.parseZoneConf(setting.Bind.LibPath + filename)
		if err != nil {
			log.Printf("Error loading %s: %s\n", filename, err)
			continue
		}

		Service.Zones[zConf.Origin] = zConf
		fmt.Println("Loaded domain file", filename)
	}
}

func (bs *BindService) parseBindConf(filename string) (*parser.BindConf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bindConf, err := parser.ConfParser.Parse(filename, file)
	if err != nil {
		return nil, err
	}

	return bindConf, nil
}

func (bs *BindService) parseZoneConf(filename string) (*parser.ZoneConf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	zConf, err := parser.ZoneParser.Parse(filename, file)
	if err != nil {
		return nil, err
	}

	return zConf, nil
}

func (bs *BindService) CreateZone(data *schemas.ZoneData) (*parser.ZoneConf, error) {
	// Get write access to the filesystem and release it when done
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	// Validate that the new zone is not defined already
	if _, ok := bs.Zones[data.Origin]; ok {
		return nil, fmt.Errorf("zone %s exists already", data.Origin)
	}

	// Create the new zone from received data
	zConf := &parser.ZoneConf{
		Origin: data.Origin,
		Ttl:    data.Ttl,
		SOARecord: &parser.SOARecord{
			NameServer: data.NameServer,
			Admin:      data.Admin,
			Refresh:    data.Refresh,
			Retry:      data.Retry,
			Expire:     data.Expire,
			Minimum:    data.Minimum,
		},
		Records: []parser.Record{},
	}

	bindConf := *bs.BindConf

	if err := bindConf.AddZone(zConf); err != nil {
		return nil, err
	}

	// Write new changes to BIND files and rollback on error
	rollbackZConf, err := zConf.WriteToDisk(zConf.GetFilename())
	if err != nil {
		rollbackZConf()
		return nil, err
	}

	rollbackBindConf, err := bindConf.WriteToDisk(bs.ZonesFilePath)
	if err != nil {
		rollbackZConf()
		rollbackBindConf()
		return nil, err
	}

	// Notify BIND about the new update
	if err := bs.Reconfig(); err != nil {
		rollbackZConf()
		rollbackBindConf()
		return nil, err
	}

	// Sync changes on memory
	bs.BindConf = &bindConf
	bs.Zones[zConf.Origin] = zConf

	return zConf, nil
}

func (bs *BindService) UpdateZone(targetOrigin string, data *schemas.ZoneData) (*parser.ZoneConf, error) {
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	zConfPointer, ok := bs.Zones[targetOrigin]
	if !ok {
		return nil, fmt.Errorf("zone %s does not exist", targetOrigin)
	}

	ZConf := *zConfPointer

	ZConf.Ttl = data.Ttl
	ZConf.SOARecord.NameServer = data.NameServer
	ZConf.SOARecord.Admin = data.Admin
	ZConf.SOARecord.Refresh = data.Refresh
	ZConf.SOARecord.Retry = data.Retry
	ZConf.SOARecord.Expire = data.Expire
	ZConf.SOARecord.Minimum = data.Minimum

	rollback, err := ZConf.WriteToDisk(ZConf.GetFilename())
	if err != nil {
		rollback()
		return nil, err
	}

	if err := bs.ReloadZone(targetOrigin); err != nil {
		rollback()
		return nil, err
	}

	bs.Zones[targetOrigin] = &ZConf

	return &ZConf, nil
}

func (bs *BindService) DeleteZone(origin string) error {
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	targetZConf, ok := bs.Zones[origin]
	if !ok {
		return fmt.Errorf("domain %s does not exist", origin)
	}

	rollbackDConf, err := targetZConf.DeleteFromDisk(targetZConf.GetFilename())
	if err != nil {
		rollbackDConf()
		return err
	}

	bindConf := *bs.BindConf

	if err := bindConf.DeleteZone(targetZConf); err != nil {
		rollbackDConf()
		return err
	}

	rollbackBindConf, err := bindConf.WriteToDisk(bs.ZonesFilePath)
	if err != nil {
		rollbackDConf()
		rollbackBindConf()
		return err
	}

	if err := bs.Reconfig(); err != nil {
		rollbackDConf()
		rollbackBindConf()
		return err
	}

	bs.BindConf = &bindConf
	delete(bs.Zones, origin)

	return nil
}

func (bs *BindService) AddRecord(origin string, record parser.Record) error {
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	targetZConf, ok := bs.Zones[origin]
	if !ok {
		return errors.New("origin not found")
	}

	zConf := *targetZConf

	if err := zConf.AddRecord(record); err != nil {
		return err
	}

	rollback, err := zConf.WriteToDisk(zConf.GetFilename())
	if err != nil {
		rollback()
		return err
	}

	if err := bs.ReloadZone(origin); err != nil {
		rollback()
		return err
	}

	bs.Zones[origin] = &zConf

	return nil
}

func (bs *BindService) UpdateRecord(origin, target string, record parser.Record) error {
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	targetZConf, ok := bs.Zones[origin]
	if !ok {
		return errors.New("origin not found")
	}

	zConf := *targetZConf

	if err := zConf.UpdateRecord(target, record); err != nil {
		return err
	}

	rollback, err := zConf.WriteToDisk(zConf.GetFilename())
	if err != nil {
		rollback()
		return err
	}

	if err := bs.ReloadZone(origin); err != nil {
		rollback()
		return err
	}

	bs.Zones[origin] = &zConf

	return nil
}

func (bs *BindService) DeleteRecord(origin string, record parser.Record) error {
	bs.Mutex.Lock()
	defer bs.Mutex.Unlock()

	targetZConf, ok := bs.Zones[origin]
	if !ok {
		return errors.New("origin not found")
	}

	zConf := *targetZConf

	if err := zConf.DeleteRecord(record); err != nil {
		return err
	}

	rollback, err := zConf.WriteToDisk(zConf.GetFilename())
	if err != nil {
		rollback()
		return err
	}

	if err := bs.ReloadZone(origin); err != nil {
		rollback()
		return err
	}

	bs.Zones[origin] = &zConf

	return nil
}

func (bs *BindService) exec(command ...string) error {
	// was used as reference for this method the docker-cli exec command implementation
	// https://github.com/docker/cli/blob/1163b4609978e0e6f2b2629b59c4a62d348e1466/cli/command/container/exec.go#L99

	if _, err := bs.DockerCli.ContainerInspect(bs.ctx, bs.ContainerId); err != nil {
		return err
	}

	execCreateConfig := &types.ExecConfig{
		User:         "bind",
		Privileged:   false,
		Tty:          false,
		AttachStdin:  false,
		AttachStderr: true,
		AttachStdout: false,
		Detach:       false,
		DetachKeys:   "",
		Env:          []string{},
		WorkingDir:   "/",
		Cmd:          command,
	}

	response, err := bs.DockerCli.ContainerExecCreate(bs.ctx, bs.ContainerId, *execCreateConfig)
	if err != nil {
		return err
	}
	if response.ID == "" {
		return errors.New("exec ID empty")
	}

	execStartConfig := &types.ExecStartCheck{
		Detach: execCreateConfig.Detach,
		Tty:    execCreateConfig.Tty,
	}

	if err := bs.DockerCli.ContainerExecStart(bs.ctx, response.ID, *execStartConfig); err != nil {
		return err
	}

	return nil
}

// Runs `rndc reconfig` in the BIND server.
// Reloads the configuration file and loads new zones, but does not reload existing zone files even if they have changed.
func (bs *BindService) Reconfig() error {
	return bs.exec("rndc", "reconfig")
}

// Runs `rndc reload {zone}` in the BIND server
func (bs *BindService) ReloadZone(zone string) error {
	return bs.exec("rndc", "reload", zone)
}
