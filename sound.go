package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/huyinghuan/aliyun-voice/asr"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type InVoice struct {
	buff []byte
	seq  int
	end  int
	vid  string
}

func recordSeg() {
	var seg = make([]int16, 0)
	defer func() {
		fmt.Println("-----over-----")
	}()
	inputChannels := 1
	outputChannels := 0
	sampleRate := 8000
	framesPerBuffer := make([]int16, 1024)

	// init PortAudio

	_ = portaudio.Initialize()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)

	// setup Wave file writer
	//audioSeg := "tmp/" + time.Now().Format("2006-01-02-15-04-05") + ".wav"
	//waveFile, err := os.Create(audioSeg)
	//param := wave.WriterParam{
	//	Out:           waveFile,
	//	Channel:       inputChannels,
	//	SampleRate:    sampleRate,
	//	BitsPerSample: 16, // if 16, change to WriteSample16()
	//}
	//waveWriter, err := wave.NewWriter(param)
	//errCheck(err)
	//rand.Seed(time.Now().UnixNano())
	// start reading from microphone
	fmt.Println("----------init------")
	errCheck(stream.Start())
	tick := time.Tick(2000 * time.Millisecond)
	//client := getTcClient()
	var seq int
	vid := RandStringBytesMaskImpr(16)
	for {
		errCheck(stream.Read())
		select {
		case <-tick:
			fmt.Println("---------TICK-------")
			fmt.Printf("---seg len is %d---\n", len(seg))
			transeg := make([]int16, len(seg))
			copy(transeg, seg)
			//go baiduTranslate(transeg)
			//auKey := RandStringBytesMaskImpr(8)
			b := int16ToByte(transeg)
			//TcTrans(b, auKey, client)
			in := InVoice{buff: b, seq: seq, end: 0, vid: vid}
			seq++
			go tecentTranslate(in)
			seg = make([]int16, 0)
			errCheck(err)
		case <-endChan:
			fmt.Printf("---ctx done seg len is %d---\n", len(seg))
			transeg := make([]int16, len(seg))
			copy(transeg, seg)
			//go baiduTranslate(transeg)
			//auKey := RandStringBytesMaskImpr(8)
			b := int16ToByte(transeg)
			//TcTrans(b, auKey, client)
			in := InVoice{buff: b, seq: seq, end: 1, vid: vid}
			seq++
			go tecentTranslate(in)
			_ = stream.Close()
			_ = portaudio.Terminate()
			return
		default:
			// write to wave file
			//_, err := waveWriter.WriteSample16(framesPerBuffer) // WriteSample16 for 16 bits
			//seg = append(seg, framesPerBuffer...)
			seg = append(seg, framesPerBuffer...)
			//in := InVoice{seg, seq, 0}
			//seq++
			errCheck(err)
		}
	}
}

type BaiduAccessTokenRes struct {
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	SessionSecret    string `json:"session_secret"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type BaiduReqBody struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	DevPid  int    `json:"dev_pid"`
	Channel int    `json:"channel"`
	Cuid    string `json:"cuid"`
	Speech  string `json:"speech"`
	Len     int    `json:"len"`
}

var accessToken string

func baiduTranslate(seg []int16) {

	at := getAccessToken()
	uri := `https://aip.baidubce.com/rpc/2.0/bicc/v1/general?access_token=` + at
	cli := &http.Client{}
	cfg := &tls.Config{
		MaxVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}
	cli.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
	b := make([]byte, 0)
	for _, i16 := range seg {
		var h, l = uint8(i16 >> 8), uint8(i16 & 0xff)
		b = append(b, l)
		b = append(b, h)
	}
	speech := base64.StdEncoding.EncodeToString(b)
	var body = BaiduReqBody{
		Format:  "pcm",
		Rate:    8000,
		DevPid:  0,
		Channel: 1,
		Cuid:    "selinplus@163.com",
		Speech:  speech,
		Len:     len(seg),
	}
	bb, err := json.Marshal(&body)
	res, err := cli.Post(uri, "Content-Type:application/json", bytes.NewReader(bb))
	if err != nil {
		errCheck(err)
	}
	errCheck(err)
	defer res.Body.Close()
	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ReadAll err=%v\n", err)
		return
	}
	fmt.Printf("BODY = %v\n", string(bd))
}
func getAccessToken() string {
	if accessToken != "" {
		fmt.Println("----factory access token return----")
		return accessToken
	} else {
		var bats = BaiduAccessTokenRes{}
		accessTokenUrl := `https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=gLPovGQHXUL47so5qGCyG0Fu&client_secret=1boFq32mwPGG86AlSEf7DHezphRu5ylA`
		cli := &http.Client{}
		cfg := &tls.Config{
			MaxVersion:               tls.VersionTLS11, // try tls.VersionTLS10 if this doesn't work
			PreferServerCipherSuites: true,
		}
		cli.Transport = &http.Transport{
			TLSClientConfig: cfg,
		}
		req, err := http.NewRequest("POST", accessTokenUrl, nil)
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		res, err := cli.Do(req)
		if err != nil {
			fmt.Printf("EE:%v\n", err)
			return ""
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		err = json.Unmarshal(body, &bats)
		if err != nil {
			fmt.Printf("unmarshal error err=%v\n", err)
			return ""
		}
		accessToken = bats.AccessToken
		return bats.AccessToken
	}
}
func aliTranslate(in InVoice) {
	auth := asr.GetAuth("", "")
	fw, _ := filepath.Abs("in.file")
	bytesOfFile, _ := ioutil.ReadFile(fw)
	result, e := auth.GetOneWord(bytesOfFile)
	if e != nil {
		fmt.Errorf("err:%v", e)
	}
	fmt.Println(result)
}
func tecentTranslate(i InVoice) {
	//fw, _ := filepath.Abs(i.file)

	signTemplate := `POSTasr.cloud.tencent.com/asr/v1/1258311667?end=%d&engine_model_type=8k_0&expired=%d&needvad=1&nonce=52811334&projectid=0&res_type=1&result_text_format=0&secretid=AKID5HNNI69U2UCqPs4dv7FZlsqeIn9LujFf&seq=%d&source=0&sub_service_type=1&timeout=5000&timestamp=%d&voice_format=1&voice_id=%s`
	timeStamp := time.Now().Unix()
	expiredStamp := time.Now().AddDate(0, 0, 2).Unix()
	//unique := RandStringBytesMaskImpr(16)
	signStr := fmt.Sprintf(signTemplate, i.end, expiredStamp, i.seq, timeStamp, i.vid)
	secretKey := `Zy94k5E8pHUSAdw1rSqzsXxHL9hz2pMd`

	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	reqTemplate := `http://asr.cloud.tencent.com/asr/v1/1258311667?end=%d&engine_model_type=8k_0&expired=%d&needvad=1&nonce=52811334&projectid=0&res_type=1&result_text_format=0&secretid=AKID5HNNI69U2UCqPs4dv7FZlsqeIn9LujFf&seq=%d&source=0&sub_service_type=1&timeout=5000&timestamp=%d&voice_format=1&voice_id=%s`
	reqStr := fmt.Sprintf(reqTemplate, i.end, expiredStamp, i.seq, timeStamp, i.vid)

	fmt.Printf("URL: %s\n", reqStr)
	fmt.Printf("LEN is : %d\n", len(i.buff))
	reader := bytes.NewReader(i.buff)
	cli := &http.Client{}

	req, err := http.NewRequest("POST", reqStr, reader)
	req.Header.Add("Host", "asr.cloud.tencent.com")
	req.Header.Add("Authorization", sign)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(i.buff)))

	res, err := cli.Do(req)
	if err != nil {
		fmt.Printf("EE:%v\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ReadAll err=%v\n", err)
		return
	}
	fmt.Printf("seq is %d,BODY = %v\n", i.seq, string(body))
}
func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}
