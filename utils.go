package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"github.com/qingni918/utils/zaplogger"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func EncodeGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	// 确保所有数据都被写入
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CalcFuncCostTime(logKey string, f func()) {
	timeBegin := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	defer func() {
		log.Printf("calcFuncCostTime processed, logKey: %s, cost time: %s", logKey, time.Since(timeBegin).String())
	}()
	f()
}

type Logger struct {
	*log.Logger
}

func (l *Logger) PrintErr(err error, exit ...bool) {
	if err != nil {
		l.Println(err.Error())
	}

	if len(exit) > 0 && exit[0] == true {
		os.Exit(1)
	}
}

func GetLogger() *Logger {
	lg := log.Default()
	lg.SetFlags(log.Ldate | log.Ltime | log.Lshortfile /* | log.LUTC*/)
	return &Logger{lg}
}

func GetZapLoggerWith(serviceID string, level zapcore.Level, outputMode int, bdInfo string, devMode bool) *zaplogger.ZapLogger {
	logger := zaplogger.NewLogger(serviceID, level, outputMode, devMode)
	return logger
}

func GetCaller() string {
	_, file, line, ok := runtime.Caller(2)
	callerStr := "GetCaller - failed"
	if ok {
		callerStr = fmt.Sprintf("%s:%d", file, line)
	}
	return callerStr
}

func GetCallerWithSkip(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	callerStr := "GetCaller - failed"
	if ok {
		callerStr = fmt.Sprintf("%s:%d", file, line)
	}
	return callerStr
}

// ReadFile 去除BOM文件格式
func ReadFile(fp string) ([]byte, error) {
	fileContent, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	fileContent = bytes.TrimPrefix(fileContent, []byte{239, 187, 191})
	return fileContent, err
}

func Padding(str, pad string, length int, prev bool) string {
	lenStr := len(str)
	if lenStr >= length {
		return str
	}
	padStr := ""
	for i := lenStr; i < length; i++ {
		padStr += pad
	}

	if prev {
		return padStr + str
	}

	return str + padStr
}

func a2bHex(hexStr string) ([]byte, error) {
	// DecodeString 将十六进制字符串解码为字节切片
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Des(pText, pKey string) string {
	plaintext := []byte(pText)
	key, _ := a2bHex(pKey)
	log.Println(string(key))

	ret, _ := encryptDES_ECB(key, plaintext, false)
	retStr := hex.EncodeToString(ret)
	log.Println("DesEncrypt:", retStr)
	return retStr
}

// DES ECB 加密
func encryptDES_ECB(key, plaintext []byte, pad bool) ([]byte, error) {
	// 创建 DES 加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充明文以满足块大小
	if pad {
		plaintext = pkcs5Padding(plaintext, block.BlockSize())
	}

	// ECB 模式：逐块加密
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], plaintext[i:i+block.BlockSize()])
	}

	ret, _ := decryptDES_ECB(key, ciphertext, pad)
	log.Println("decryptDES_ECB:", string(ret))

	return ciphertext, nil
}

// DES ECB 解密
func decryptDES_ECB(key, ciphertext []byte, pad bool) ([]byte, error) {
	// 创建 DES 解密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// ECB 模式：逐块解密
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(plaintext[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}

	// 去除填充
	if pad {
		plaintext = pkcs5UnPadding(plaintext)
	}

	return plaintext, nil
}

// PKCS5 填充
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5 去除填充
func pkcs5UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func LMHash(str string) string {
	// 明文转化成大写
	str = strings.ToUpper(str)
	// 转16进制
	str = fmt.Sprintf("%x", str)
	lenStr := len(str)
	// 字符串长度不足14B(28)，用0补全
	if lenStr < 28 {
		paddingStr := ""
		for i := lenStr; i < 28; i++ {
			paddingStr += "0"
		}
		str += paddingStr
	}
	// 将14B分为两组，每组7B，然后转二进制
	ga := str[:14]
	gb := str[14:]
	// 16 => 10
	gai, _ := strconv.ParseInt(ga, 16, 64)
	gbi, _ := strconv.ParseInt(gb, 16, 64)
	// 10 => 2 最高位补齐
	ga = Padding(strconv.FormatInt(gai, 2), "0", 56, true)
	gb = Padding(strconv.FormatInt(gbi, 2), "0", 56, true)
	// 将每组二进制数据按7bit一组，分为8组，每组末尾添0，再转16进制
	padAndToX := func(str string) string {
		stepLen := 7
		dest2 := ""
		for i := 0; i < 8; i++ {
			s := str[i*stepLen : (i+1)*stepLen]
			s += "0"
			si, _ := strconv.ParseInt(s, 2, 64)
			s = Padding(strconv.FormatInt(si, 2), "0", 8, true)
			log.Println(s)
			dest2 += s
		}
		// 2 => 10
		dest10, _ := strconv.ParseInt(dest2, 2, 64)
		// 10 => 16
		dest16 := Padding(strconv.FormatInt(dest10, 16), "0", 16, true)
		return dest16
	}

	ga = padAndToX(ga)
	gb = padAndToX(gb)
	ga = Des("KGS!@#$%", ga)
	gb = Des("KGS!@#$%", gb)
	return ga + gb
}
