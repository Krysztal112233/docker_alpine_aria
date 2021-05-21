package main

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kpango/glg"
	"github.com/mitchellh/go-ps"
)

var confMap map[string]string = make(map[string]string)
var configPath = os.Getenv("ARIA2_CONFIG_PATH")

var statu chan bool = make(chan bool)
var errorsChan chan error = make(chan error)
var dockerSingal chan int = make(chan int)

var defaultConfig = `#save-session=aria2.session
#save-session-interval=60
dir=/tmp/Downloads
#disk-cache=32M
#file-allocation=none
continue=true
#max-concurrent-downloads=5
max-connection-per-server=5
min-split-size=10M
#split=5
#max-overall-download-limit=0
#max-download-limit=0
#max-overall-upload-limit=0
#max-upload-limit=0
#disable-ipv6=true
#timeout=60
#max-tries=5
#retry-wait=0
enable-rpc=true
rpc-allow-origin-all=true
rpc-listen-all=true
#event-poll=select
#rpc-listen-port=6800
#rpc-secret=<TOKEN>
#rpc-user=<USER>
#rpc-passwd=<PASSWD>
#rpc-secure=true
#rpc-certificate=/path/to/certificate.pem
#rpc-private-key=/path/to/certificate.key
#follow-torrent=true
listen-port=51413
#bt-max-peers=55
enable-dht=false
#enable-dht6=false9
#dht-listen-port=6881-6999
#bt-enable-lpd=false
enable-peer-exchange=false
#bt-request-peer-speed-limit=50K
peer-id-prefix=-TR2770-
user-agent=Transmission/2.77
seed-ratio=0
#force-save=false
#bt-hash-check-seed=true
bt-seed-unverified=true
bt-save-metadata=true`

func init() {

	//Process for the config
	tmp := strings.Split(strings.ReplaceAll(defaultConfig, "-", "_"), "\n")
	for _, v := range tmp {
		s := strings.Split(v, "=")
		if s[0][0] == '#' {
			if t := os.Getenv(s[0][1:]); len(t) == 0 {
				confMap[strings.ReplaceAll(s[0], "_", "-")] = s[1]
			} else {
				confMap[strings.ReplaceAll(s[0][1:], "_", "-")] = t
			}
			continue
		} else {
			if t := os.Getenv(s[0]); len(t) == 0 {
				confMap[strings.ReplaceAll(s[0], "_", "-")] = s[1]
			} else {
				confMap[strings.ReplaceAll(s[0], "_", "-")] = t
			}
		}
	}
}

func genConfigFile() (b []byte) {
	tmp := make([]string, 0)
	for k, v := range confMap {
		tmp = append(tmp, k+"="+v)
	}
	return []byte(strings.Join(tmp, "\n"))
}

func main() {
	b := genConfigFile()

	if len(configPath) == 0 {
		configPath = "./aria2.conf"
		fp, err := os.Create(configPath)
		if err != nil {
			glg.Error(err)
			os.Exit(-1)
		}
		fp.Write(b)
	}

	s, _ := exec.LookPath("aria2c")
	aria2Cmd := s + " " + "--conf-path=" + configPath + " " + "-D"
	cmd := exec.Command(s, "--conf-path="+configPath, "-D")
	glg.Info(aria2Cmd)
	for k, v := range confMap {
		glg.Logf("%-30s:%s", k, v)
	}
	cmd.Start()
	time.Sleep(time.Second * 1)
	go checkStatus()

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	for {
		select {
		case <-statu:
			glg.Info("Aria2c had exited!")
			return
		case erro := <-errorsChan:
			glg.Error(erro)
			return
		case sig := <-c:
			switch sig {
			case os.Kill:
				glg.Error("container has been killed!")
				os.Exit(0)
			case syscall.SIGTERM:
				glg.Info("container was been stoped by docker!")
				os.Exit(0)
			case os.Interrupt:
				glg.Info("container has been interrupted!")
				os.Exit(0)
			}
		}
	}

}

func checkStatus() {
	var isAria2cRunning bool = false
	for {
		l, err := ps.Processes()
		if err != nil {
			glg.Error(err)
			statu <- false
			errorsChan <- err
			break
		}

		for _, v := range l {
			if v.Executable() == "aria2c" {
				isAria2cRunning = true
			}
		}

		if !isAria2cRunning {
			statu <- false
			errorsChan <- errors.New("aria2c has stoped!")
		}
	}
}
