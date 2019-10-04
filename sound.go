package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
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
	file string
	seq  int
	end  int
}

func recordSeg(ctx context.Context) {

	defer func() {
		fmt.Println("-----over-----")
	}()
	// www.people.csail.mit.edu/hubert/pyaudio/  - under the Record tab
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	//framesPerBuffer := make([]byte, 64)
	framesPerBuffer := make([]byte, 1000)
	//var seg []byte
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
	var seq int
	//var tick = time.Tick(3000 * time.Millisecond)
	for {
		errCheck(stream.Read())
		select {
		//case <-tick:
		//	in := InVoice{"", seq, 0}
		//	fmt.Println("---------TICK-------")
		//	//go aliTranslate(in)
		//	go tecentTranslate(in, seg)
		//	seq++
		//
		//	errCheck(err)
		case <-ctx.Done():
			in := InVoice{"", seq, 1}
			go tecentTranslate(in, framesPerBuffer)
			_ = stream.Close()
			_ = portaudio.Terminate()
			return
		default:
			// write to wave file
			//_, err := waveWriter.WriteSample16(framesPerBuffer) // WriteSample16 for 16 bits
			//seg = append(seg, framesPerBuffer...)
			in := InVoice{"", seq, 0}
			seq++
			go tecentTranslate(in, framesPerBuffer)
			errCheck(err)
		}
	}
}

type BaiduAccessTokenRes struct {
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        string `json:"expires_in"`
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

func baiduTranslate(in InVoice) {

	at := getAccessToken()
	uri := `https://aip.baidubce.com/rpc/2.0/bicc/v1/general?access_token=` + at
	cli := &http.Client{}

	var body = BaiduReqBody{Format: "pcm", Rate: 8000, DevPid: 0, Channel: 1, Cuid: "IYTU@selinplus", Speech: in.file, Len: 22}
	bb, err := json.Marshal(&body)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(bb))
	if err != nil {
		errCheck(err)
	}
	cli.Do(req)

}
func getAccessToken() string {
	var bats = BaiduAccessTokenRes{}
	accessTokenUrl := `https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=cKpjeMvYmO0dLEHOH9KYRR0O&client_secret=1D9iSmFR2kIHinbuMbEAdeHfluhTYott&`
	cli := &http.Client{}
	req, err := http.NewRequest("POST", accessTokenUrl, nil)
	res, err := cli.Do(req)
	if err != nil {
		fmt.Printf("EE:%v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &bats)
	if err != nil {
		fmt.Printf("unmarshal error err=%v\n", err)
		return ""
	}
	return bats.AccessToken
}
func aliTranslate(in InVoice) {
	auth := asr.GetAuth("", "")
	fw, _ := filepath.Abs(in.file)
	bytesOfFile, _ := ioutil.ReadFile(fw)
	result, e := auth.GetOneWord(bytesOfFile)
	if e != nil {
		fmt.Errorf("err:%v", e)
	}
	fmt.Println(result)
}
func tecentTranslate(i InVoice, seg []byte) {
	//fw, _ := filepath.Abs(i.file)

	signTemplate := `POSTasr.cloud.tencent.com/asr/v1/1258311667?end=%d&engine_model_type=16k_0&expired=%d&needvad=0&nonce=52811334&projectid=0&res_type=1&result_text_format=0&secretid=AKID5HNNI69U2UCqPs4dv7FZlsqeIn9LujFf&seq=%d&source=0&sub_service_type=1&timeout=5000&timestamp=%d&voice_format=1&voice_id=%s`
	timeStamp := time.Now().Unix()
	expiredStamp := time.Now().AddDate(0, 0, 2).Unix()
	unique := RandStringBytesMaskImpr(16)
	signStr := fmt.Sprintf(signTemplate, i.end, expiredStamp, i.seq, timeStamp, unique)
	secretKey := `Zy94k5E8pHUSAdw1rSqzsXxHL9hz2pMd`

	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	reqTemplate := `http://asr.cloud.tencent.com/asr/v1/1258311667?end=%d&engine_model_type=16k_0&expired=%d&needvad=0&nonce=52811334&projectid=0&res_type=1&result_text_format=0&secretid=AKID5HNNI69U2UCqPs4dv7FZlsqeIn9LujFf&seq=%d&source=0&sub_service_type=1&timeout=5000&timestamp=%d&voice_format=1&voice_id=%s`
	reqStr := fmt.Sprintf(reqTemplate, i.end, expiredStamp, i.seq, timeStamp, unique)

	fmt.Printf("URL: %s\n", reqStr)
	fmt.Printf("LEN is : %d\n", len(seg))
	reader := bytes.NewReader(seg)
	cli := &http.Client{}
	req, err := http.NewRequest("POST", reqStr, reader)
	req.Header.Add("Host", "asr.cloud.tencent.com")
	req.Header.Add("Authorization", sign)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(seg)))
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
	seg = make([]byte, 0)
	fmt.Printf("seq is %d,BODY = %v\n", i.seq, string(body))
}
func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}
