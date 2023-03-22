/*
kunpengsecl licensed under the Mulan PSL v2.
You can use this software according to the terms and conditions of
the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at:
    http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.

Author: wangli/wanghaijing
Create: 2022-04-01
Description: Implement specific services provided by AS
*/

package akissuer

/*
#cgo CFLAGS: -I../../tee/tverlib/verifier
#cgo LDFLAGS: -L${SRCDIR}/../../tee/tverlib/verifier -lteeverifier -Wl,-rpath=${SRCDIR}/../../tee/tverlib/verifier
#include "teeverifier.h"
*/
import "C"
import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"miracl/core"
	"miracl/core/FP512BN"
	"sync"
	"time"
	"unsafe"

	"gitee.com/openeuler/kunpengsecl/attestation/tas/config"
)

const (
	// ZERO_VALUE means the value of zero.
	ZERO_VALUE = 0
	// UINT32_BYTES_LENGTH means the length of uint32 bytes.
	UINT32_BYTES_LENGTH = 4
	// UINT64_BYTES_LENGTH means the length of uint64 bytes.
	UINT64_BYTES_LENGTH = 8
	// tag type
	// KEY_TAG_TYPE_MOVE_BITS means the bits of key tag type
	KEY_TAG_TYPE_MOVE_BITS = 28
	// RA_INTEGER means ra integer
	RA_INTEGER = (1 << KEY_TAG_TYPE_MOVE_BITS)
	// RA_BYTES means ra bytes
	RA_BYTES = (2 << KEY_TAG_TYPE_MOVE_BITS)
	// RA_TAG_SIGN_TYPE means ra tag sign type
	RA_TAG_SIGN_TYPE = RA_INTEGER
	// RA_TAG_HASH_TYPE means ra tag hash type
	RA_TAG_HASH_TYPE = RA_INTEGER | 1
	// RA_TAG_CURVE_TYPE means ra tag curve type
	RA_TAG_CURVE_TYPE = RA_INTEGER | 2
	// RA_TAG_QTA_IMG_HASH means ra tag qta img hash
	RA_TAG_QTA_IMG_HASH = RA_BYTES
	// RA_TAG_TA_IMG_HASH means ra tag ta img hash
	RA_TAG_TA_IMG_HASH = RA_BYTES | 1
	// RA_TAG_QTA_MEM_HASH means ra tag qta mem hash
	RA_TAG_QTA_MEM_HASH = RA_BYTES | 2
	// RA_TAG_TA_MEM_HASH means ra tag ta mem hash
	RA_TAG_TA_MEM_HASH = RA_BYTES | 3
	// RA_TAG_RESERVED means ra tag reserved
	RA_TAG_RESERVED = RA_BYTES | 4
	// RA_TAG_AK_PUB means ra tag ak pub
	RA_TAG_AK_PUB = RA_BYTES | 5
	// RA_TAG_SIGN_DRK means ra tag sign drk
	RA_TAG_SIGN_DRK = RA_BYTES | 6
	// RA_TAG_SIGN_AK means ra tag sign ak
	RA_TAG_SIGN_AK = RA_BYTES | 7
	// RA_TAG_CERT_DRK means ra tag cert drk
	RA_TAG_CERT_DRK = RA_BYTES | 8
	// RA_TAG_CERT_AK means ra tag cert ak
	RA_TAG_CERT_AK = RA_BYTES | 9
	// RA_ALG_RSA_3072 means the code name of
	// RSA algorithm with thr key length of 3072
	RA_ALG_RSA_3072 = 0x20000
	// RA_ALG_RSA_4096 means the code name of
	// RSA algorithm with thr key length of 4096
	RA_ALG_RSA_4096 = 0x20001 // PSS padding
	// RA_ALG_SHA_256 means the code name of
	// SHA256 algorithm
	RA_ALG_SHA_256 = 0x20002
	// RA_ALG_SHA_384 means the code name of
	// SHA384 algorithm
	RA_ALG_SHA_384 = 0x20003
	// RA_ALG_SHA_512 means the code name of
	// SHA512 algorithm
	RA_ALG_SHA_512 = 0x20004
	// RA_ALG_ECDSA means the code name of
	// ECDSA algorithm
	RA_ALG_ECDSA = 0x20005
	// RA_ALG_ED25519 means the code name of
	// ED25519 algorithm
	RA_ALG_ED25519 = 0x20006
	// RA_ALG_SM2_DSA_SM3 means the code name of
	// DSA SM3 algorithm
	RA_ALG_SM2_DSA_SM3 = 0x20007
	// RA_ALG_SM3 means the code name of
	// SM3 algorithm
	RA_ALG_SM3 = 0x20008
	// RA_ALG_DAA_GRP_FP512BN means the code name of
	// DAA GRP FP512BN algorithm
	RA_ALG_DAA_GRP_FP512BN = 0x20009
	// x509 cert template default value
	strChina      = "China"
	strCompany    = "Company"
	strCommonName = "AK Server"
	// RA_SCENARIO_NO_AS means the code name of ra scenario without as
	RA_SCENARIO_NO_AS = 0
	// RA_SCENARIO_AS_NO_DAA means the code name of ra scenario without daa
	RA_SCENARIO_AS_NO_DAA = 1
	// RA_SCENARIO_AS_WITH_DAA means the code name of ra scenario with daa
	RA_SCENARIO_AS_WITH_DAA = 2
)

type (
	ra_data_offset struct {
		data_len    uint32
		data_offset uint32
	}
	ra_params struct {
		tags uint32
		buf  interface{}
	}
	certificate struct {
		version     uint32
		timestamp   uint64
		tag         string
		param_count uint32
		params      []ra_params
	}
	daacre struct {
		akcre struct {
			A *FP512BN.ECP
			B *FP512BN.ECP
			C *FP512BN.ECP
			D *FP512BN.ECP
		}
		sigma struct {
			u *FP512BN.BIG
			j *FP512BN.BIG
		}
	}
)

var (
	m            sync.Mutex
	serialNumber int64 = 1
)

func verifyAKCert(oldAKCert []byte) (drkpub *rsa.PublicKey, akpub []byte, err error) {
	// STEP1: get data used for verify
	var c_cert, c_certdrk, c_akpub C.buffer_data
	c_cert.size = C.uint(len(oldAKCert))
	up_old_cert := C.CBytes(oldAKCert)
	defer C.free(up_old_cert)
	c_cert.buf = (*C.uchar)(up_old_cert)
	C.tee_get_akcert_data(&c_cert, &c_akpub, &c_certdrk)
	drkcertbyte := []byte(C.GoBytes(unsafe.Pointer(c_certdrk.buf), C.int(c_certdrk.size)))
	// STEP2: get data used for re-sign
	akpub = []byte(C.GoBytes(unsafe.Pointer(c_akpub.buf), C.int(c_akpub.size)))
	// STEP3: parse device cert
	drkcertBlock, _ := pem.Decode(drkcertbyte)
	drkcert, err := x509.ParseCertificate(drkcertBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	drkpub = drkcert.PublicKey.(*rsa.PublicKey)
	log.Print("Server: Parse drk cert succeeded.")
	// STEP4: verify device cert signature
	err = verifyDRKSig(drkcert)
	if err != nil {
		return nil, nil, errors.New("verify drk signature failed")
	}
	log.Print("Server: Verify drk signature ok.")
	// STEP5: verify ak cert signature & QCA
	up_qca_ref := C.CBytes([]byte(config.GetBaseValue()))
	defer C.free(up_qca_ref)
	rs := C.tee_verify_akcert(&c_cert, 3, (*C.char)(up_qca_ref))
	if !bool(rs) {
		return nil, nil, errors.New("verify ak signature failed")
	}
	log.Print("Server: Verify ak signature & QCA ok.")
	return drkpub, akpub, nil
}

// The input parameter is the AK certificate issued by the target platform device certificate
// GenerateNoDAAAKCert after receiving the AK certificate, parses and extracts the signed data fields,
// signature fields, and DRK certificate fields
// Parse the DRK certificate
// Use huawei Level-2 certificate to check the DRK certificate.
// If the DRK certificate passes the check, the DRK certificate is trusted
// Use the DRK certificate to check the AK certificate.
// If the AK certificate passes the check, the AK certificate is trusted
// Re-sign the AK certificate using the AS private key
// Return the re-signed AK certificate
func GenerateNoDAAAKCert(oldAKCert []byte) ([]byte, error) {
	_, akpubbyte, err := verifyAKCert(oldAKCert)
	if err != nil {
		return nil, err
	}
	b := big.NewInt(0)
	b.SetBytes(akpubbyte)
	akpub := &rsa.PublicKey{
		N: b,
		E: 0x10001,
	}
	// STEP6: get as private key and as cert
	asprivkey := config.GetASPrivKey()
	ascert := config.GetASCert()
	// STEP7: re-sign ak cert
	newCertDer, err := signForAKCert(oldAKCert, ascert, akpub, asprivkey)
	if err != nil {
		return nil, err
	}
	log.Print("Server: re-sign ak cert ok.")
	newCertBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertDer,
	}
	newCertPem := pem.EncodeToMemory(newCertBlock)
	return newCertPem, nil
}

func parseCustomCert(certByte []byte) (*certificate, error) {
	if certByte == nil {
		return nil, errors.New("empty cert")
	}
	var c *certificate = &certificate{}
	c.params = make([]ra_params, 0)
	var offset int = 4
	c.version = bytesToUint(certByte[:offset]).(uint32)
	c.timestamp = bytesToUint(certByte[offset : offset+8]).(uint64)
	offset += 8
	c.tag = string(certByte[offset : offset+32])
	offset += 32
	c.param_count = bytesToUint(certByte[offset : offset+4]).(uint32)
	offset += 4
	for i := 0; i < int(c.param_count); i++ {
		tags := bytesToUint(certByte[offset : offset+4]).(uint32)
		switch tags {
		case RA_TAG_SIGN_TYPE:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_SIGN_TYPE
			c.params[i].buf = bytesToUint(certByte[offset+4 : offset+12]).(uint64)
			offset += 12
		case RA_TAG_HASH_TYPE:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_HASH_TYPE
			c.params[i].buf = bytesToUint(certByte[offset+4 : offset+12]).(uint64)
			offset += 12
		case RA_TAG_QTA_IMG_HASH:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_QTA_IMG_HASH
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		case RA_TAG_QTA_MEM_HASH:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_QTA_MEM_HASH
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		case RA_TAG_RESERVED:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_RESERVED
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		case RA_TAG_AK_PUB:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_AK_PUB
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		case RA_TAG_SIGN_DRK:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_SIGN_DRK
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		case RA_TAG_CERT_DRK:
			c.params = append(c.params, ra_params{})
			c.params[i].tags = RA_TAG_CERT_DRK
			c.params[i].buf = ra_data_offset{
				data_len:    bytesToUint(certByte[offset+4 : offset+8]).(uint32),
				data_offset: bytesToUint(certByte[offset+8 : offset+12]).(uint32),
			}
			offset += 12
		}
	}
	return c, nil
}

func bytesToUint(b []byte) interface{} {
	bytesBuffer := bytes.NewBuffer(b)
	length := len(b)
	switch length {
	case UINT32_BYTES_LENGTH:
		var x uint32
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
		return x
	case UINT64_BYTES_LENGTH:
		var x uint64
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
		return x
	}
	return nil
}

func verifyDRKSig(c *x509.Certificate) error {
	hwcert := config.GetHWCert()
	err := c.CheckSignatureFrom(hwcert)
	if err != nil {
		return err
	}
	return nil
}

// AS will uses its private key/cert to re-sign AK cert,
// and convert AK cert to x509 format
func signForAKCert(cb []byte, parent *x509.Certificate, pub interface{}, priv interface{}) ([]byte, error) {
	var ACtemplate = x509.Certificate{
		NotBefore:      time.Now(),
		NotAfter:       time.Now().AddDate(1, 0, 0),
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		IsCA:           false,
		MaxPathLenZero: true,
	}
	var akcertDer []byte
	// parse ak cert signed by device cert
	c, err := parseCustomCert(cb)
	if err != nil {
		return nil, err
	}
	// set AK Cert id, version and signature field
	m.Lock()
	id := serialNumber
	serialNumber++
	m.Unlock()
	ACtemplate.SerialNumber = big.NewInt(id)
	// extract sign algorithm
	alg_sign := extractSignAlg(c)
	if alg_sign == ZERO_VALUE {
		return nil, errors.New("sign type not found")
	}
	switch alg_sign {
	case RA_ALG_RSA_4096:
		ACtemplate.PublicKeyAlgorithm = x509.RSA
	default:
		return nil, errors.New("signature algorithm not support yet")
	}
	// extract hash algorithm
	alg_hash := extractHashAlg(c)
	if alg_hash == ZERO_VALUE {
		return nil, errors.New("hash type not found")
	}
	switch alg_hash {
	case RA_ALG_SHA_256:
		ACtemplate.SignatureAlgorithm = x509.SHA256WithRSAPSS
	default:
		return nil, errors.New("hash algorithm not support yet")
	}

	akcertDer, err = x509.CreateCertificate(rand.Reader, &ACtemplate, parent, pub, priv)
	if err != nil {
		return nil, err
	}

	return akcertDer, nil
}

func extractHashAlg(c *certificate) uint64 {
	for i := 0; i < int(c.param_count); i++ {
		if c.params[i].tags == RA_TAG_HASH_TYPE {
			alg_hash := c.params[i].buf.(uint64)
			return alg_hash
		}
	}
	return 0
}

func extractSignAlg(c *certificate) uint64 {
	for i := 0; i < int(c.param_count); i++ {
		if c.params[i].tags == RA_TAG_SIGN_TYPE {
			alg_sign := c.params[i].buf.(uint64)
			return alg_sign
		}
	}
	return 0
}

// GenerateDAAAKCert verifies old ak cert
// and gets daa grp private key configuration to generate daa credential.
func GenerateDAAAKCert(oldAKCert []byte) ([]byte, error) {
	drkcertpubkey, akpubbyte, err := verifyAKCert(oldAKCert)
	if err != nil {
		return nil, err
	}
	skxstr, skystr := config.GetDAAGrpPrivKey()
	// generate cert
	sig, err := makeDAACredential(akpubbyte, skxstr, skystr, drkcertpubkey)
	if err != nil {
		return nil, errors.New("make daa credential failed")
	}
	return sig, nil
}

func str2chunk(str string) (*FP512BN.BIG, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return nil, errors.New("string to hex failed")
	}
	res := FP512BN.FromBytes(bytes)
	return res, nil
}

func (dcre *daacre) combineu(P1 *FP512BN.ECP, Qs *FP512BN.ECP, R_B *FP512BN.ECP, R_D *FP512BN.ECP, n *FP512BN.BIG) {
	var buffer bytes.Buffer

	P1tmp := ecp2bytes(P1)
	P1bytes := []byte{0x01, 0x00, 0x00, 0x00, P1tmp[35], 0x01, 0x00, 0x00, 0x00, P1tmp[71]}
	Qsbytes := ecp2bytes(Qs)
	Abytes := ecp2bytes(dcre.akcre.A)
	Bbytes := ecp2bytes(dcre.akcre.B)
	Cbytes := ecp2bytes(dcre.akcre.C)
	Dbytes := ecp2bytes(dcre.akcre.D)
	RBbytes := ecp2bytes(R_B)
	RDbytes := ecp2bytes(R_D)

	buffer.Write(int2bytes(len(P1bytes)))
	buffer.Write(P1bytes)
	buffer.Write(int2bytes(len(Qsbytes)))
	buffer.Write(Qsbytes)
	buffer.Write(int2bytes(len(Abytes)))
	buffer.Write(Abytes)
	buffer.Write(int2bytes(len(Bbytes)))
	buffer.Write(Bbytes)
	buffer.Write(int2bytes(len(Cbytes)))
	buffer.Write(Cbytes)
	buffer.Write(int2bytes(len(Dbytes)))
	buffer.Write(Dbytes)
	buffer.Write(int2bytes(len(RBbytes)))
	buffer.Write(RBbytes)
	buffer.Write(int2bytes(len(RDbytes)))
	buffer.Write(RDbytes)
	comb := buffer.Bytes()

	hash := sha512.New()
	hash.Write(comb)
	ubytes := hash.Sum(nil)
	u := FP512BN.FromBytes(ubytes)
	zero := FP512BN.NewBIG()
	dcre.sigma.u = FP512BN.Modadd(u, zero, n) // u.Mod(n)
}

func int2bytes(n int) []byte {
	x := int32(n)
	bytesBuffer := new(bytes.Buffer)
	err := binary.Write(bytesBuffer, binary.LittleEndian, x)
	if err != nil {
		return nil
	}
	return bytesBuffer.Bytes()
}

func ecp2bytes(E *FP512BN.ECP) []byte {
	MB := int(FP512BN.MODBYTES)
	var b [2*FP512BN.MODBYTES + 2*4]byte
	xbytes, ybytes := b[4:4+MB], b[8+MB:]
	E.GetX().ToBytes(xbytes)
	E.GetY().ToBytes(ybytes)

	copy(b[:], int2bytes(len(xbytes)))
	copy(b[4+MB:], int2bytes(len(ybytes)))
	return b[:]
}

func bytes2ecp(b []byte) *FP512BN.ECP {
	MB := int(FP512BN.MODBYTES)
	xbytes, ybytes := b[4:4+MB], b[8+MB:]

	px := FP512BN.FromBytes(xbytes[:])
	py := FP512BN.FromBytes(ybytes[:])
	return FP512BN.NewECPbigs(px, py)
}

func encryptAESGCM(plaintext []byte, key []byte) ([]byte, []byte, []byte, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}
	cipher, tag := core.GCM_ENCRYPT(key, nonce, nil, plaintext)
	return cipher, tag, nonce, nil
}

func cipher1(K []byte, Qs *FP512BN.ECP, drkpubk *rsa.PublicKey) ([]byte, error) {
	// size||Qs||size||K
	var buffer bytes.Buffer
	Qsbytes := ecp2bytes(Qs)

	buffer.Write(int2bytes(len(Qsbytes)))
	buffer.Write(Qsbytes)
	buffer.Write(int2bytes(len(K)))
	buffer.Write(K)
	cleartext1 := buffer.Bytes()
	// RSA4096, padding: TEE_ALG_RSAES_PKCS1_OAEP_MGF1_SHA256
	c1, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, drkpubk, cleartext1, nil)
	if err != nil {
		return nil, errors.New("daa cipher1 failed")
	}
	return c1, nil
}

func (dcre *daacre) cipher2(K []byte, Qs *FP512BN.ECP) ([]byte, []byte, []byte, error) {
	// size||AKCert.A||size||AKCert.B||size||AKCert.C||size||AKCert.D||size||u||size||j
	var buffer bytes.Buffer
	var ubytes, jbytes [FP512BN.MODBYTES]byte

	Abytes := ecp2bytes(dcre.akcre.A)
	Bbytes := ecp2bytes(dcre.akcre.B)
	Cbytes := ecp2bytes(dcre.akcre.C)
	Dbytes := ecp2bytes(dcre.akcre.D)
	dcre.sigma.u.ToBytes(ubytes[:])
	dcre.sigma.j.ToBytes(jbytes[:])

	buffer.Write(int2bytes(len(Abytes)))
	buffer.Write(Abytes)
	buffer.Write(int2bytes(len(Bbytes)))
	buffer.Write(Bbytes)
	buffer.Write(int2bytes(len(Cbytes)))
	buffer.Write(Cbytes)
	buffer.Write(int2bytes(len(Dbytes)))
	buffer.Write(Dbytes)
	buffer.Write(int2bytes(len(ubytes)))
	buffer.Write(ubytes[:])
	buffer.Write(int2bytes(len(jbytes)))
	buffer.Write(jbytes[:])
	cleartext2 := buffer.Bytes()
	// AES-GCM key256bit
	c2, tag, iv, err := encryptAESGCM(cleartext2, K)
	if err != nil {
		return nil, nil, nil, errors.New("daa cipher2 failed")
	}
	return c2, tag, iv, nil
}

func makeDAACredential(akprip1 []byte, skxstr string, skystr string, drkpubk *rsa.PublicKey) ([]byte, error) {
	rnd := core.NewRAND()
	var rndraw [128]byte
	if _, err := io.ReadFull(rand.Reader, rndraw[:]); err != nil {
		return nil, errors.New("daa K generation failed")
	}
	rnd.Seed(len(rndraw), rndraw[:])
	r := FP512BN.Random(rnd)
	l := FP512BN.Random(rnd)
	// daa private key
	skx, err := str2chunk(skxstr)
	if err != nil {
		return nil, errors.New("daa private key data conversion failed")
	}
	sky, err := str2chunk(skystr)
	if err != nil {
		return nil, errors.New("daa private key data conversion failed")
	}

	dcre := new(daacre)
	// A=[r]P_1
	P1 := FP512BN.ECP_generator()
	dcre.akcre.A = P1.Mul(r)
	// B=[y]A
	dcre.akcre.B = dcre.akcre.A.Mul(sky)
	// D=[ry]Q_s
	n := FP512BN.NewBIGints(FP512BN.CURVE_Order)
	ry := FP512BN.Modmul(r, sky, n) // n = bnp256_order
	Qs := bytes2ecp(akprip1)
	dcre.akcre.D = Qs.Mul(ry)
	// tmp=A+D
	tmp := FP512BN.NewECP()
	tmp.Copy(dcre.akcre.A)
	tmp.Add(dcre.akcre.D)
	// C=[x]tmp
	dcre.akcre.C = tmp.Mul(skx)
	// R_B=[l]P1
	R_B := P1.Mul(l)
	// R_D=[l]Qs
	R_D := Qs.Mul(l)
	// u=sha256(P1,Qs,AKCert,R_B,R_D)
	dcre.combineu(P1, Qs, R_B, R_D, n)
	// j=l+yru (mod n), sigma=(u,j)
	yru := FP512BN.Modmul(ry, dcre.sigma.u, n)
	dcre.sigma.j = FP512BN.Modadd(yru, l, n)
	// use DRK to encrypt K -> ENCdrk(Qs||K)
	K := FP512BN.Random(rnd)
	var Kbytes [FP512BN.MODBYTES]byte
	K.ToBytes(Kbytes[:])
	cip1, err := cipher1(Kbytes[:], Qs, drkpubk)
	if err != nil {
		return nil, errors.New("daa cipher1 generation failed")
	}
	// choose a K to encrypt cert -> ENC(certAK||sigma)
	cip2, tag, iv, err := dcre.cipher2(Kbytes[:], Qs)
	if err != nil {
		return nil, errors.New("daa cipher2 generation failed")
	}
	var buffer bytes.Buffer
	buffer.Write(int2bytes(len(cip1)))
	buffer.Write(cip1)
	buffer.Write(tag)
	buffer.Write(iv)
	buffer.Write(int2bytes(len(cip2)))
	buffer.Write(cip2)
	cip := buffer.Bytes()
	return cip, nil
}

// GenerateAKCert generates ak cert according to scenario.
func GenerateAKCert(oldAKCert []byte, scenario int32) ([]byte, error) {
	switch scenario {
	case RA_SCENARIO_AS_NO_DAA:
		newCert, err := GenerateNoDAAAKCert(oldAKCert)
		if err != nil {
			log.Print("NoDAA scenario: Generate AK Cert failed!")
			return nil, err
		}
		log.Print("NoDAA scenario: Generate AK Cert succeeded!")
		return newCert, nil
	case RA_SCENARIO_AS_WITH_DAA:
		newCert, err := GenerateDAAAKCert(oldAKCert)
		if err != nil {
			log.Print("DAA scenario: Generate AK Cert failed!")
			return nil, err
		}
		log.Print("DAA scenario: Generate AK Cert succeeded!")
		return newCert, nil
	}
	return nil, errors.New("do not need to access as")
}
