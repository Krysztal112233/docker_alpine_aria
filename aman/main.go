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
bt-save-metadata=true
bt-tracker=http://104.238.198.186:8000/announce,http://1337.abcvg.info:80/announce,http://185.148.3.231:80/announce,http://193.37.214.12:6969/announce,http://195.201.31.194:80/announce,http://51.15.55.204:1337/announce,http://51.79.71.167:80/announce,http://51.81.46.170:6969/announce,http://54.36.126.137:6969/announce,http://54.39.179.91:6699/announce,http://60-fps.org:80/bt:80/announce,.phphttp://95.107.48.115:80/announce,http://[2001:1b10:1000:8101:0:242:ac11:2]:6969/announce,http://[2001:470:1:189:0:1:2:3]:6969/announce,http://[2a04:ac00:1:3dd8::1:2710]:2710/announce,http://all4nothin.net:80/announce,.phphttp://alltorrents.net:80/bt:80/announce,.phphttp://atrack.pow7.com:80/announce,http://baibako.tv:80/announce,http://big-boss-tracker.net:80/announce,.phphttp://bithq.org:80/announce,.phphttp://bt.1000.pet:2712/announce,http://bt.10000.pet:2714/announce,http://bt.3dmgame.com:2710/announce,http://bt.ali213.net:8080/announce,http://bt.firebit.org:2710/announce,http://bt.okmp3.ru:2710/announce,http://bt.unionpeer.org:777/announce,http://bt.zlofenix.org:81/announce,http://bttracker.debian.org:6969/announce,http://btx.anifilm.tv:80/announce,.phphttp://concen.org:6969/announce,http://data-bg.net:80/announce,.phphttp://datascene.net:80/announce,.phphttp://explodie.org:6969/announce,http://finbytes.org:80/announce,.phphttp://freerainbowtables.com:6969/announce,http://h4.trakx.nibba.trade:80/announce,http://ipv4announce,.sktorrent.eu:6969/announce,http://irrenhaus.dyndns.dk:80/announce,.phphttp://kinorun.com:80/announce,.phphttp://masters-tb.com:80/announce,.phphttp://mediaclub.tv:80/announce,.phphttp://mixfiend.com:6969/announce,http://music-torrent.net:2710/announce,http://mvgroup.org:2710/announce,http://ns3107607.ip-54-36-126.eu:6969/announce,http://ns349743.ip-91-121-106.eu:80/announce,http://nyaa.tracker.wf:7777/announce,http://open.acgnxtracker.com:80/announce,http://openbittorrent.com:80/announce,http://opentracker.i2p.rocks:6969/announce,http://opentracker.xyz:80/announce,http://p4p.arenabg.com:1337/announce,http://pow7.com:80/announce,http://proaudiotorrents.org:80/announce,.phphttp://retracker.hotplug.ru:2710/announce,http://retracker.sevstar.net:2710/announce,http://retracker.spark-rostov.ru:80/announce,http://retracker.telecom.by:80/announce,http://rt.tace.ru:80/announce,http://secure.pow7.com:80/announce,http://share.camoe.cn:8080/announce,http://siambit.com:80/announce,.phphttp://t.acg.rip:6699/announce,http://t.nyaatracker.com:80/announce,http://t.overflow.biz:6969/announce,http://t1.pow7.com:80/announce,http://t2.pow7.com:80/announce,http://torrent-team.net:80/announce,.phphttp://torrent.arjlover.net:2710/announce,http://torrent.fedoraproject.org:6969/announce,http://torrent.mp3quran.net:80/announce,.phphttp://torrent.resonatingmedia.com:6969/announce,http://torrent.rus.ec:2710/announce,http://torrent.ubuntu.com:6969/announce,http://torrentclub.online:54123/announce,http://torrentsmd.com:8080/announce,http://torrenttracker.nwc.acsalaska.net:6969/announce,http://torrentzilla.org:80/announce,http://tr.cili001.com:8070/announce,http://tr.kxmp.cf:80/announce,http://tracker-cdn.moeking.me:2095/announce,http://tracker.ali213.net:8080/announce,http://tracker.anirena.com:80/announce,http://tracker.anirena.com:80/b16a15d9a238d1f59178d3614b857290/announce,http://tracker.anonwebz.xyz:8080/announce,http://tracker.birkenwald.de:6969/announce,http://tracker.bittor.pw:1337/announce,http://tracker.breizh.pm:6969/announce,http://tracker.bt4g.com:2095/announce,http://tracker.ccp.ovh:6969/announce,http://tracker.dler.org:6969/announce,http://tracker.etree.org:6969/announce,http://tracker.fdn.fr:6969/announce,http://tracker.files.fm:6969/announce,http://tracker.frozen-layer.net:6969/announce,.phphttp://tracker.gbitt.info:80/announce,http://tracker.gcvchp.com:2710/announce,http://tracker.gigatorrents.ws:2710/announce,http://tracker.grepler.com:6969/announce,http://tracker.internetwarriors.net:1337/announce,http://tracker.ipv6tracker.ru:80/announce,http://tracker.lelux.fi:80/announce,http://tracker.loadbt.com:6969/announce,http://tracker.minglong.org:8080/announce,http://tracker.noobsubs.net:80/announce,http://tracker.openbittorrent.com:80/announce,http://tracker.opentrackr.org:1337/announce,http://tracker.pow7.com:80/announce,http://tracker.pussytorrents.org:3000/announce,http://tracker.shittyurl.org:80/announce,http://tracker.sloppyta.co:80/announce,http://tracker.tambovnet.org:80/announce,.phphttp://tracker.tasvideos.org:6969/announce,http://tracker.tfile.co:80/announce,http://tracker.tfile.me:80/announce,http://tracker.trackerfix.com:80/announce,http://tracker.xdvdz.com:2710/announce,http://tracker.zerobytes.xyz:1337/announce,http://tracker1.bt.moack.co.kr:80/announce,http://tracker2.dler.org:80/announce,http://tracker3.dler.org:2710/announce,http://tracker4.itzmx.com:2710/announce,http://trk.publictracker.xyz:6969/announce,http://vpn.flying-datacenter.de:6969/announce,http://vps02.net.orel.ru:80/announce,http://www.all4nothin.net:80/announce,.phphttp://www.freerainbowtables.com:6969/announce,http://www.legittorrents.info:80/announce,.phphttp://www.thetradersden.org/forums/tracker:80/announce,.phphttp://www.tribalmixes.com:80/announce,.phphttp://www.tvnihon.com:6969/announce,http://www.wareztorrent.com:80/announce,http://www.worldboxingvideoarchive.com:80/announce,.phphttp://www.xwt-classics.net:80/announce,.phphttp:/:80/announce,.partis.si:80:80/announce,https://1337.abcvg.info:443/announce,https://carapax.net:443/announce,https://open.kickasstracker.com:443/announce,https://opentracker.acgnx.se:443/announce,https://torrents.linuxmint.com:443/announce,.phphttps://tr.ready4.icu:443/announce,https://tracker.bt-hash.com:443/announce,https://tracker.coalition.space:443/announce,https://tracker.foreverpirates.co:443/announce,https://tracker.gbitt.info:443/announce,https://tracker.iriseden.fr:443/announce,https://tracker.lelux.fi:443/announce,https://tracker.lilithraws.cf:443/announce,https://tracker.nanoha.org:443/announce,https://tracker.nitrix.me:443/announce,https://tracker.parrotsec.org:443/announce,https://tracker.shittyurl.org:443/announce,https://tracker.sloppyta.co:443/announce,https://tracker.tamersunion.org:443/announce,https://trakx.herokuapp.com:443/announce,https://w.wwwww.wtf:443/announce,udp://103.196.36.31:6969/announce,udp://103.30.17.23:6969/announce,udp://104.238.159.144:6969/announce,udp://104.238.198.186:8000/announce,udp://104.244.153.245:6969/announce,udp://104.244.72.77:1337/announce,udp://109.248.43.36:6969/announce,udp://119.28.134.203:6969/announce,udp://138.68.171.1:6969/announce,udp://144.76.35.202:6969/announce,udp://144.76.82.110:6969/announce,udp://148.251.53.72:6969/announce,udp://149.28.47.87:1738/announce,udp://151.236.218.182:6969/announce,udp://157.90.161.74:6969/announce,udp://157.90.169.123:80/announce,udp://159.65.202.134:6969/announce,udp://159.69.208.124:6969/announce,udp://163.172.170.127:6969/announce,udp://167.179.77.133:1/announce,udp://173.212.223.237:6969/announce,udp://176.123.5.238:3391/announce,udp://178.159.40.252:6969/announce,udp://185.181.60.67:80/announce,udp://185.21.216.185:6969/announce,udp://185.243.215.40:6969/announce,udp://185.8.156.2:6969/announce,udp://193.34.92.5:80/announce,udp://195.201.94.195:6969/announce,udp://198.100.149.66:6969/announce,udp://198.50.195.216:7777/announce,udp://199.195.249.193:1337/announce,udp://199.217.118.72:6969/announce,udp://205.185.121.146:6969/announce,udp://208.83.20.20:6969/announce,udp://209.141.45.244:1337/announce,udp://209.141.59.16:6969/announce,udp://212.1.226.176:2710/announce,udp://212.47.227.58:6969/announce,udp://212.83.181.109:6969/announce,udp://217.12.218.177:2710/announce,udp://37.1.205.89:2710/announce,udp://37.235.174.46:2710/announce,udp://37.59.48.81:6969/announce,udp://45.33.83.49:6969/announce,udp://45.56.65.82:54123/announce,udp://45.76.92.209:6969/announce,udp://46.101.244.237:6969/announce,udp://46.148.18.250:2710/announce,udp://46.148.18.254:2710/announce,udp://47.ip-51-68-199.eu:6969/announce,udp://5.206.31.154:6969/announce,udp://51.15.2.221:6969/announce,udp://51.15.41.46:6969/announce,udp://51.68.199.47:6969/announce,udp://51.68.34.33:6969/announce,udp://51.79.81.233:6969/announce,udp://52.58.128.163:6969/announce,udp://62.168.229.166:6969/announce,udp://6ahddutb1ucc3cp.ru:6969/announce,udp://6rt.tace.ru:80/announce,udp://78.30.254.12:2710/announce,udp://88.99.142.4:8000/announce,udp://9.rarbg.me:2710/announce,udp://9.rarbg.to:2710/announce,udp://91.121.145.207:6969/announce,udp://91.149.192.31:6969/announce,udp://91.216.110.52:451/announce,udp://[2001:1b10:1000:8101:0:242:ac11:2]:6969/announce,udp://[2001:470:1:189:0:1:2:3]:6969/announce,udp://[2a03:7220:8083:cd00::1]:451/announce,udp://[2a04:ac00:1:3dd8::1:2710]:2710/announce,udp://[2a0f:e586:f:f::220]:6969/announce,udp://admin.videoenpoche.info:6969/announce,udp://anidex.moe:6969/announce,udp://app.icon256.com:8000/announce,udp://bt.firebit.org:2710/announce,udp://bt.okmp3.ru:2710/announce,udp://bt2.54new.com:8080/announce,udp://bubu.mapfactor.com:6969/announce,udp://cdn-1.gamecoast.org:6969/announce,udp://code2chicken.nl:6969/announce,udp://concen.org:6969/announce,udp://cutiegirl.ru:6969/announce,udp://daveking.com:6969/announce,udp://discord.heihachi.pw:6969/announce,udp://drumkitx.com:6969/announce,udp://edu.uifr.ru:6969/announce,udp://engplus.ru:6969/announce,udp://exodus.desync.com:6969/announce,udp://explodie.org:6969/announce,udp://fe.dealclub.de:6969/announce,udp://free.publictracker.xyz:6969/announce,udp://inferno.demonoid.is:3391/announce,udp://ipv4.tracker.harry.lu:80/announce,udp://ipv6.tracker.zerobytes.xyz:16661/announce,udp://johnrosen1.com:6969/announce,udp://line-net.ru:6969/announce,udp://ln.mtahost.co:6969/announce,udp://mail.realliferpg.de:6969/announce,udp://movies.zsw.ca:6969/announce,udp://mts.tvbit.co:6969/announce,udp://nagios.tks.sumy.ua:80/announce,udp://newtoncity.org:6969/announce,udp://open.demonii.com:1337/announce,udp://open.publictracker.xyz:6969/announce,udp://open.stealth.si:80/announce,udp://openbittorrent.com:80/announce,udp://opentor.org:2710/announce,udp://opentracker.i2p.rocks:6969/announce,udp://opentrackr.org:1337/announce,udp://p4p.arenabg.ch:1337/announce,udp://p4p.arenabg.com:1337/announce,udp://peerfect.org:6969/announce,udp://pow7.com:80/announce,udp://public-tracker.zooki.xyz:6969/announce,udp://public.publictracker.xyz:6969/announce,udp://public.tracker.vraphim.com:6969/announce,udp://qg.lorzl.gq:6969/announce,udp://retracker.hotplug.ru:2710/announce,udp://retracker.lanta-net.ru:2710/announce,udp://retracker.netbynet.ru:2710/announce,udp://retracker.nts.su:2710/announce,udp://retracker.sevstar.net:2710/announce,udp://sugoi.pomf.se:80/announce,udp://t1.leech.ie:1337/announce,udp://t2.leech.ie:1337/announce,udp://t3.leech.ie:1337/announce,udp://thetracker.org:80/announce,udp://torrentclub.online:54123/announce,udp://tr.bangumi.moe:6969/announce,udp://tr2.ysagin.top:2710/announce,udp://tracker-de.ololosh.space:6969/announce,udp://tracker.0x.tf:6969/announce,udp://tracker.altrosky.nl:6969/announce,udp://tracker.army:6969/announce,udp://tracker.beeimg.com:6969/announce,udp://tracker.birkenwald.de:6969/announce,udp://tracker.bittor.pw:1337/announce,udp://tracker.blacksparrowmedia.net:6969/announce,udp://tracker.breizh.pm:6969/announce,udp://tracker.btsync.gq:2710/announce,udp://tracker.ccp.ovh:6969/announce,udp://tracker.coppersurfer.tk:6969/announce,udp://tracker.cyberia.is:6969/announce,udp://tracker.dler.com:6969/announce,udp://tracker.dler.org:6969/announce,udp://tracker.edkj.club:6969/announce,udp://tracker.filemail.com:6969/announce,udp://tracker.grepler.com:6969/announce,udp://tracker.halfchub.club:6969/announce,udp://tracker.internetwarriors.net:1337/announce,udp://tracker.kali.org:6969/announce,udp://tracker.kuroy.me:5944/announce,udp://tracker.lelux.fi:6969/announce,udp://tracker.moeking.me:6969/announce,udp://tracker.monitorit4.me:6969/announce,udp://tracker.nrx.me:6969/announce,udp://tracker.ololosh.space:6969/announce,udp://tracker.open-internet.nl:6969/announce,udp://tracker.openbittorrent.com:6969/announce,udp://tracker.openbittorrent.com:80/announce,udp://tracker.opentrackr.org:1337/announce,udp://tracker.sbsub.com:2710/announce,udp://tracker.shkinev.me:6969/announce,udp://tracker.sktorrent.net:6969/announce,udp://tracker.skyts.net:6969/announce,udp://tracker.swateam.org.uk:2710/announce,udp://tracker.theoks.net:6969/announce,udp://tracker.tiny-vps.com:6969/announce,udp://tracker.torrent.eu.org:451/announce,udp://tracker.uw0.xyz:6969/announce,udp://tracker.zemoj.com:6969/announce,udp://tracker.zerobytes.xyz:1337/announce,udp://tracker0.ufibox.com:6969/announce,udp://tracker2.dler.com:80/announce,udp://tracker2.dler.org:80/announce,udp://tracker4.itzmx.com:2710/announce,udp://u.wwwww.wtf:1/announce,udp://udp-tracker.shittyurl.org:6969/announce,udp://us-tracker.publictracker.xyz:6969/announce,udp://valakas.rollo.dnsabr.com:2710/announce,udp://vibe.community:6969/announce,udp://vibe.sleepyinternetfun.xyz:1738/announce,udp://wassermann.online:6969/announce,udp://www.mvgroup.org:2710/announce,udp://www.torrent.eu.org:451/announce,udp://z.mercax.com:53/announce,ws://tracker.sloppyta.co:80/announce,wss://tracker.openwebtorrent.com:443/announce
`

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
