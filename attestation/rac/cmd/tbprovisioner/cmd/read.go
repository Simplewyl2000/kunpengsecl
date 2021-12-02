/*
Copyright (c) Huawei Technologies Co., Ltd. 2021.
kunpengsecl licensed under the Mulan PSL v2.
You can use this software according to the terms and conditions of
the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at:
    http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.

Author: wucaijun
Create: 2021-12-02
Description: Command line tool for tpm provision process.
*/
package cmd

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gitee.com/openeuler/kunpengsecl/attestation/rac/ractools"
	"github.com/spf13/cobra"
)

const (
	// for command flags const long and strings.
	LongIndex      = "index"
	ShortIndex     = "i"
	LongFile       = "file"
	ShortFile      = "f"
	LongType       = "type"
	ShortType      = "t"
	LongSimulator  = "simulator"
	ShortSimulator = "s"
	// for input parameter const strings.
	ConstNVRAM = "NVRAM"
	ConstPCR   = "PCR"
	ConstPEM   = "PEM"
	ConstDER   = "DER"
	// for error messages.
	errOpenTPM    = "can't open TPM, %v\n"
	errReadNVRAM  = "read NVRAM(0x%08X) error: %v\n"
	errReadPCR    = "read PCR(%d) error: %v\n"
	errReadFile   = "read file(%s) error: %v\n"
	errWriteNVRAM = "write NVRAM(0x%08X) error: %v\n"
	errWriteFile  = "write file(%s) error: %v\n"
	// for information message.
	constWriteNVRAM = "NVRAM(0x%08X): %v\n"
	constWritePCR   = "PCR(%d): %v\n"
	// long help string
	constIndexList = `
    0x01C00002   RSA 2048 EK Certificate
    0x01C00003   RSA 2048 EK Nonce
    0x01C00004   RSA 2048 EK Template
    0x01C0000A   ECC NIST P256 EK Certificate
    0x01C0000B   ECC NIST P256 EK Nonce
    0x01C0000C   ECC NIST P256 EK Template
    0x01C00012   RSA 2048 EK Certificate (H-1)
    0x01C00014   ECC NIST P256 EK Certificate (H-2)
    0x01C00016   ECC NIST P384 EK Certificate (H-3)
    0x01C00018   ECC NIST P512 EK Certificate (H-4)
    0x01C0001A   ECC SM2_P256 EK Certificate (H-5)
    0x01C0001C   RSA 3072 EK Certificate (H-6)
    0x01C0001E   RSA 4096 EK Certificate (H-7)`
	constRIndexHelp = `the NVRAM index to read from (DEFAULT: 0x01C00002)` + constIndexList
	constWIndexHelp = `the NVRAM index to write to (DEFAULT: 0x01C00002)` + constIndexList
	constEIndexHelp = `the NVRAM index to erase (DEFAULT: 0x01C00002)` + constIndexList
)

var (
	// readCmd represents the read command
	readCmd = &cobra.Command{
		Use:   "read [nvram|pcr]",
		Short: "Read TPM NVRAM|PCR content",
		Long: `Use this command to read TPM resources(nvram, pcr, etc) and save to
file(Default: StdOut) with the define format(Default: PEM).
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println(cmd.Long)
				fmt.Println(cmd.UsageString())
				os.Exit(1)
			}
			obj := strings.ToUpper(args[0])
			if obj != ConstNVRAM && obj != ConstPCR {
				fmt.Println(cmd.Long)
				fmt.Println(cmd.UsageString())
				os.Exit(1)
			}
			tp, err := ractools.OpenTPM(!rSim)
			if err != nil {
				fmt.Printf(errOpenTPM, err)
				os.Exit(1)
			}
			defer tp.Close()
			switch obj {
			case ConstNVRAM:
				nvramHandle(tp)
			case ConstPCR:
				pcrHandle(tp)
			}
		},
	}
	rIndex uint32
	rFile  string
	rType  string
	rSim   bool
)

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	readCmd.Flags().Uint32VarP(&rIndex, LongIndex, ShortIndex, 0, constRIndexHelp)
	readCmd.Flags().StringVarP(&rFile, LongFile, ShortFile, "", "the file to write content into")
	readCmd.Flags().StringVarP(&rType, LongType, ShortType, "", "the output type PEM/DER (DEFAULT: PEM)")
	readCmd.Flags().BoolVarP(&rSim, LongSimulator, ShortSimulator, false, "use the simulator to test (DEFAULT: false)")
}

func nvramHandle(tp *ractools.TPM) {
	if rIndex == 0 {
		rIndex = ractools.IndexRsa2048EKCert
	}
	buf, err := tp.ReadEKCert(rIndex)
	if err != nil {
		fmt.Printf(errReadNVRAM, rIndex, err)
		os.Exit(1)
	}
	var out []byte
	if rType == ConstDER {
		out = buf
	} else {
		block := &pem.Block{Bytes: buf}
		out = pem.EncodeToMemory(block)
	}
	rType = strings.ToUpper(rType)
	if rFile == "" {
		fmt.Printf(constWriteNVRAM, rIndex, out)
	} else {
		err = ioutil.WriteFile(rFile, out, 0644)
		if err != nil {
			fmt.Printf(errWriteFile, rFile, err)
		}
	}
}

func pcrHandle(tp *ractools.TPM) {
	//fmt.Printf(errReadPCR, rIndex, err)
	//fmt.Printf(constWritePCR, rIndex, buf)
}