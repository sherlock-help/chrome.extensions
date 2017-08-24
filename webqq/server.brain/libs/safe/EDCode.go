package safe

import(
	//bulit-in
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"crypto/hmac"
	
	"crypto/sha1"

	"crypto/md5"
	
	"encoding/hex"
	"errors"
	"sync"
	"io"
	"fmt"
	//extends
	. "bakerstreet-club/logs"
)


// key
const sCodeKey = "bakerstreet.club"

var (
    block cipher.Block
    mutex sync.Mutex
    iCBCBlockSize int
)

func init() {

	mutex.Lock()
    defer mutex.Unlock()
 
    if block != nil {
        return
    }
 
    cblock, err := aes.NewCipher([]byte(sCodeKey))
    if err != nil {
        Error("aes.NewCipher: " + err.Error())
    }
 
    block = cblock

    iCBCBlockSize = block.BlockSize()
    //fmt.Println(iCBCBlockSize)
}

//des

//aes (CBC, ECB, CTR, OCF, CFB)
//=============================================================
func EnCodeCBC(byteSrc []byte) ([]byte, error) {
	// 验证输入参数
    // 必须为aes.Blocksize的倍数
    if len(byteSrc) % iCBCBlockSize != 0 {
        return nil, errors.New("crypto/cipher: input not full blocks")
    }
    encryptText := make([]byte, iCBCBlockSize + len(byteSrc))
    iv := encryptText[:iCBCBlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(encryptText[iCBCBlockSize:], byteSrc)
    
    return encryptText, nil
}

func DeCodeCBC(byteEnSrc []byte) ([]byte, error) {
	// hex
    decryptText, err := hex.DecodeString(fmt.Sprintf("%x", string(byteEnSrc)))
    if err != nil {
        return nil, err
    }
    // 长度不能小于aes.Blocksize
    if len(decryptText) < iCBCBlockSize {
        return nil, errors.New("crypto/cipher: ciphertext too short")
    }
    iv := decryptText[:iCBCBlockSize]
    decryptText = decryptText[iCBCBlockSize:]
    // 验证输入参数
    // 必须为aes.Blocksize的倍数
    if len(decryptText) % iCBCBlockSize != 0 {
        return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
    }
    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(decryptText, decryptText)

    return decryptText, nil 
}

//=============================================================

//res
//=============================================================

//=============================================================


//hash
//=============================================================
//sha
func EnCodeSHA(sInitStr string) string{

	//sha1
	h := sha1.New()
	io.WriteString(h, sInitStr)
	return fmt.Sprintf("%x", h.Sum(nil))	
}
//hmac
func EncodeHMAC(sInitStr string) string {
	//hmac ,use sha1
	key := []byte(sCodeKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(sInitStr))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
//md5
func EnCodeMD5(sInitStr string) string {
 	//
 	md5Ctx := md5.New()
    md5Ctx.Write([]byte(sInitStr))
    cipherStr := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherStr)
}
//=============================================================
