package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/huyinghuan/aliyun-voice/asr"
	"github.com/zenwerk/go-wave"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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

	// www.people.csail.mit.edu/hubert/pyaudio/  - under the Record tab
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	//framesPerBuffer := make([]byte, 64)
	framesPerBuffer := make([]int16, 500)

	// init PortAudio

	_ = portaudio.Initialize()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)

	// setup Wave file writer
	audioSeg := "tmp/" + time.Now().Format("2006-01-02-15-04-05") + ".wav"
	waveFile, err := os.Create(audioSeg)
	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 16, // if 16, change to WriteSample16()
	}
	waveWriter, err := wave.NewWriter(param)
	errCheck(err)
	rand.Seed(time.Now().UnixNano())
	// start reading from microphone
	errCheck(stream.Start())
	var seq int
	var tick = time.Tick(2000 * time.Millisecond)
	for {
		errCheck(stream.Read())
		select {
		case <-tick:
			err = waveWriter.Close()
			sendAudio := audioSeg
			fmt.Println(audioSeg)
			fmt.Println("----TICK----")
			if err != nil {
				fmt.Printf("ERR:%v", err)
			} else {
				in := InVoice{sendAudio, seq, 0}
				go aliTranslate(in)
				seq++
			}

			audioSeg = "tmp/" + time.Now().Format("2006-01-02-15-04-05") + ".wav"
			waveFile, err = os.Create(audioSeg)
			param.Out = waveFile
			waveWriter, err = wave.NewWriter(param)
			errCheck(err)
		case <-ctx.Done():
			sendAudio := audioSeg
			in := InVoice{sendAudio, seq, 1}
			_ = waveWriter.Close()
			go aliTranslate(in)
			_ = stream.Close()
			_ = portaudio.Terminate()
			return
		default:
			// write to wave file
			_, err := waveWriter.WriteSample16(framesPerBuffer) // WriteSample16 for 16 bits
			errCheck(err)
		}
	}
}
func aliTranslate(in InVoice) {
	auth := asr.GetAuth("LTAI4FxURV2oXxSNfWEiSJ6K", "MXGeKguYwVlyw0gd5rMWuoSn6LVHrM")
	fw, _ := filepath.Abs(in.file)
	bytesOfFile, _ := ioutil.ReadFile(fw)
	result, e := auth.GetOneWord(bytesOfFile)
	if e != nil {
		fmt.Errorf("err:%v", e)
	}
	fmt.Println(result)
}
func tecentTranslate(i InVoice) {
	fw, _ := filepath.Abs(i.file)

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
	fmt.Printf("File: %s\n", fw)
	bytesOfFile, _ := ioutil.ReadFile(fw)
	fmt.Printf("LEN is : %d\n", len(bytesOfFile))
	reader := bytes.NewReader(bytesOfFile)
	cli := &http.Client{}
	req, err := http.NewRequest("POST", reqStr, reader)
	req.Header.Add("Host", "asr.cloud.tencent.com")
	req.Header.Add("Authorization", sign)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(bytesOfFile)))
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
	fmt.Printf("BODY = %v\n", string(body))
}
func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}
