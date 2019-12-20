// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ibm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils/exec"
	"github.com/appleboy/easyssh-proxy"
)

type MakeConfig struct{
	User string
	Server string
	Password string
	Port   string
	Timeout time.Duration
}

func Executer() *easyssh.MakeConfig{
  ssh := &easyssh.MakeConfig{
        User:   username,
        Server: defaultTgtBindIp,
        Password: password,
        Port:    port,
        Timeout: 60 * time.Second,
}
return ssh
}


type Cli struct {
	// Command executer
	BaseExecuter exec.Executer
	// Command Root executer
	RootExecuter exec.Executer
}

func login() (error) {
	stdout, stderr, done, err := Executer().Run("uname", 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err!=nil{
	    return err
	}
	return nil
}

func NewCli() (*Cli, error) {
	return &Cli{
		BaseExecuter: exec.NewBaseExecuter(),
		RootExecuter: exec.NewRootExecuter(),
	}, nil
}

func (c *Cli) execute(cmd ...string) (string, error) {
	return c.RootExecuter.Run(cmd[0], cmd[1:]...)
}

// get the spectrumscale cluster status
func (c *Cli) GetSpectrumScaleStatus() error {
	createCmd := "mmgetstate"
	stdout, stderr, done, err := Executer().Run(createCmd, 10*time.Second)
        glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
	lines := strings.Split(stdout, "\n")
	var bool = strings.Contains(lines[2], "active")
	if bool != true{
	    return err
	}
	return nil
}

// get spectrumscale mount point
func (c *Cli) GetSpectrumScaleMountPoint() (string, error) {
	createCmd := "mmlsfs all -T"
	stdout, stderr, done, err := Executer().Run(createCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return "", err
	}
	var mountPoint string
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "-T") != true {
		    continue
		}
		field := strings.Fields(line)
		mountPoint = field[1]
}
	return mountPoint, nil
}

// create volume
func (c *Cli) CreateVolume(name string, size string) error {
	createCmd := "mmcrfileset" + " " + "fs1" + " " + name + " " + "--inode-space" + " " + "new"
	stdout, stderr, done, err := Executer().Run(createCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
	linkCmd := "mmlinkfileset" + " " + "fs1" + " " + name + " " + "-J /gpfs/fs1/" + name
	stdout, stderr, done, err = Executer().Run(linkCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
        //mmsetquota fs1:vol8 --block 1G:2G --files 10K:11K
	quotaCmd := "mmsetquota" + " " + "fs1" + ":" + name + " --block" + " " + size + "G" + ":" + size + "G"
	stdout, stderr, done, err = Executer().Run(quotaCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
	return err
}

// delete volume
func (c *Cli) Delete(name string) error {
	unlinkCmd := "mmunlinkfileset" + " " + "fs1" + " " + name
	stdout, stderr, done, err := Executer().Run(unlinkCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}

	delCmd := "mmdelfileset" + " " + "fs1" + " " + name + " " + "-f"
	stdout, stderr, done, err = Executer().Run(delCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
	return nil
}

// this is function for extending the volume size
func (c *Cli) ExtendVolume(name string, newSize string) error {
	quotaCmd := "mmsetquota" + " " + "fs1" + ":" + name + " --block" + " " + newSize + "G" + ":" + newSize + "G"
	fmt.Println("setquota:",quotaCmd)
	stdout, stderr, done, err := Executer().Run(quotaCmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
	return nil
}

// this is function for creating the snapshot
func (c *Cli) CreateSnapshot(snapName, volName string) error {
	cmd := "mmcrsnapshot" + " " + "fs1" + " " + snapName +  " " + "-j" + " " + volName
	stdout, stderr, done, err := Executer().Run(cmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
  return nil
}

// this is function for deleting the snapshot
func (c *Cli) DeleteSnapshot(volName, snapName string) error {
	cmd := "mmdelsnapshot" + " " + "fs1" + " " + volName + ":" + snapName
	stdout, stderr, done, err := Executer().Run(cmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return err
	}
  return nil
}

type Pools struct {
	Name          string
	TotalCapacity int64
	FreeCapacity  int64
	UUID          string
}

// this function is for discover all the pool from spectrumscale cluster
func (c *Cli) ListPools(mountPoint string) (*[]Pools, error) {
	field := strings.Split(mountPoint, "/")
        length := len(field)
	filesystem := field[length-1]
	cmd := "mmlspool" + " " + filesystem
	stdout, stderr, done, err := Executer().Run(cmd, 10*time.Second)
	glog.Infof("stdout:%v stderr:%v done:%v", stdout, stderr, done)
	if err != nil {
	    return nil, err
	}
	lines := strings.Split(stdout, "\n")
	var pols []Pools
	for _, line := range lines {
		if len(line) == 0 {
		    continue
		}
		fields := strings.Fields(line)
		if fields[0] == "Storage" {
		    continue
		}
		if fields[0] == "Name" {
		    continue
		}

		total, _ := strconv.ParseFloat(fields[6], 64)
		free, _ := strconv.ParseFloat(fields[7], 64)
		pool := Pools{
			Name:          fields[0],
			TotalCapacity: int64(total/1000000),
			FreeCapacity:  int64(free/1000000),
			UUID:          fields[1],
		}
		pols = append(pols, pool)
	}
	return &pols, nil
}
