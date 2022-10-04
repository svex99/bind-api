package bind

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
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
	ctx          context.Context
	Mutex        *sync.Mutex
	DockerCli    *client.Client
	ContainerId  string
	ZoneFilePath string
	Domains      map[string]*parser.DomainConf
}

var Service = &BindService{}

func init() {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	Service.ctx = context.Background()
	Service.Mutex = &sync.Mutex{}
	Service.DockerCli = cli
	Service.ContainerId = setting.Bind.ContainerId
	Service.ZoneFilePath = setting.Bind.ConfPath + "named.conf.local"
	Service.Domains = make(map[string]*parser.DomainConf)

	recordsDir, err := fs.ReadDir(os.DirFS(setting.Bind.RecordsPath), ".")
	if err != nil {
		panic(err)
	}

	fmt.Println(">>> Loading BIND9 domain files")
	for _, entry := range recordsDir {
		name := entry.Name()

		if !entry.IsDir() && strings.HasPrefix(name, "db.") && !strings.HasSuffix(name, ".bak") {
			dConf, err := Service.parseDomainConf(setting.Bind.RecordsPath + entry.Name())
			if err != nil {
				fmt.Printf("Error loading %s: %s", name, err)
				continue
			}

			Service.Domains[dConf.Origin] = dConf
			fmt.Println("Loaded domain file", name)
		}
	}
}

func (bs *BindService) parseDomainConf(filename string) (*parser.DomainConf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dConf, err := parser.Parser.Parse(filename, file)
	if err != nil {
		return nil, err
	}

	return dConf, nil
}

func (bm *BindService) CreateDomain(data *schemas.DomainData) (*parser.DomainConf, error) {
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()

	if _, ok := bm.Domains[data.Origin]; ok {
		return nil, fmt.Errorf("domain %s exists already", data.Origin)
	}

	dConf := &parser.DomainConf{
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
		Records: []parser.Record{
			parser.NSRecord{NameServer: data.NameServer},
			parser.ARecord{Name: data.NameServer, Ip: data.NSIp},
		},
	}

	if err := dConf.WriteToDisk(dConf.GetFilename()); err != nil {
		return nil, err
	}

	if err := bm.Reload(); err != nil {
		return nil, err
	}

	bm.Domains[dConf.Origin] = dConf

	return dConf, nil
}

func (bm *BindService) UpdateDomain(targetOrigin string, data *schemas.DomainData) (*parser.DomainConf, error) {
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()

	targetDConfPointer, ok := bm.Domains[targetOrigin]
	if !ok {
		return nil, fmt.Errorf("domain %s does not exist", targetOrigin)
	}

	targetDConf := *targetDConfPointer

	indexNS, err := targetDConf.GetRecordIndex(parser.NSRecord{NameServer: targetDConf.SOARecord.NameServer}.GetHash())
	if err != nil {
		return nil, err
	}
	targetNS := targetDConf.Records[indexNS].(parser.NSRecord)

	for i, record := range targetDConf.Records {
		switch record.(type) {
		case parser.ARecord:
			aRecord := targetDConf.Records[i].(parser.ARecord)
			if aRecord.Name == targetDConf.SOARecord.NameServer {
				aRecord.Name = data.NameServer
			}
		}
	}

	targetDConf.Ttl = data.Ttl

	targetDConf.SOARecord.NameServer = data.NameServer
	targetDConf.SOARecord.Admin = data.Admin
	targetDConf.SOARecord.Refresh = data.Refresh
	targetDConf.SOARecord.Retry = data.Retry
	targetDConf.SOARecord.Expire = data.Expire
	targetDConf.SOARecord.Minimum = data.Minimum

	targetNS.NameServer = data.NameServer

	if err := targetDConf.WriteToDisk(targetDConf.GetFilename()); err != nil {
		return nil, err
	}

	if err := bm.Reload(); err != nil {
		return nil, err
	}

	bm.Domains[targetOrigin] = &targetDConf

	return &targetDConf, nil
}

func (bm *BindService) DeleteDomain(targetOrigin string) error {
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()

	targetDConf, ok := bm.Domains[targetOrigin]
	if !ok {
		return fmt.Errorf("domain %s does not exist", targetOrigin)
	}

	if err := targetDConf.DeleteFromDisk(targetDConf.GetFilename()); err != nil {
		return err
	}

	if err := bm.Reload(); err != nil {
		return err
	}

	delete(bm.Domains, targetOrigin)

	return nil
}

func (bs *BindService) exec(command []string) error {
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
		AttachStderr: false,
		AttachStdout: false,
		Detach:       true,
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

// Reloads the configuration files for the bind service running inside a container
func (bs *BindService) Reload() error {
	return bs.exec([]string{"service", "named", "reload"})
}
