package configurer

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/naruto-june/edcoder"
)

type hostItem struct {
	Suffix string `json:"suffix,omitempty" yaml:"suffix,omitempty" toml:"suffix,omitempty"`
}

type envConfItem struct {
	Files []string            `json:"files,omitempty" yaml:"files,omitempty" toml:"files,omitempty"`
	Hosts map[string]hostItem `json:"hosts,omitempty" yaml:"hosts,omitempty" toml:"hosts,omitempty"`
}

type configure struct {
	CCI map[string]interface{} `json:"common_conf_item,omitempty" yaml:"common_conf_item,omitempty" toml:"common_conf_item,omitempty"`
	CCF map[string]string      `json:"common_conf_file,omitempty" yaml:"common_conf_file,omitempty" toml:"common_conf_file,omitempty"`
	ICF map[string]envConfItem `json:"individual_conf_file,omitempty" yaml:"individual_conf_file,omitempty" toml:"individual_conf_file,omitempty"`
}

type filePathInfo struct {
	dirPath  string
	fileName string
	fileExt  string
}

func parsePath(filePath string) (f *filePathInfo, e error) {
	filePath = filepath.ToSlash(filePath)
	dirPath, fileFullName := filepath.Split(filePath)
	fileExt := path.Ext(filePath)
	fileName := strings.TrimSuffix(fileFullName, fileExt)
	if "" == fileExt || "" == fileName {
		return nil, errors.New("illegal file format")
	}
	fileExt = fileExt[1:]

	return &filePathInfo{dirPath: dirPath, fileExt: fileExt, fileName: fileName}, nil
}

func readFile(fileName string) (string, error) {
	fi, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	var chunks []byte
	r := bufio.NewReader(fi)
	buf := make([]byte, 4096)
	if nil == buf {
		return "", errors.New("memory error(make []byte)")
	}
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if 0 == n {
			break
		}

		chunks = append(chunks, buf[:n]...)
	}

	return string(chunks), nil
}

func addConfFile(fileFullName, fileName string) error {
	_, exist := cnf.ConfFile[fileName]
	if exist {
		return errors.New("repeated config file name")
	}

	content, err := readFile(fileFullName)
	if nil != err {
		return err
	}

	cnf.ConfFile[fileName] = content
	return nil
}

// ParseConf 解析配置的配置
func ParseConf(fname string) error {
	flag.Parse()

	fmt.Println(cmdKey)

	tmp, err := parsePath(fname)
	if nil != err {
		return err
	}
	fi, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fi.Close()

	d, err := edcoder.NewDecoder(edcoder.SetDecoderExt(tmp.fileExt), edcoder.SetDecoderReader(fi))
	if nil != err {
		return err
	}

	var mcnf configure
	err = d.Decode(&mcnf)
	if nil != err {
		return err
	}

	// 公共配置项
	cnf.CommonConfItem = mcnf.CCI

	// 公共文件配置的路径
	for k := range mcnf.CCF {
		t, err := parsePath(k)
		if nil != err {
			return err
		}

		err = addConfFile(k, t.fileName)
		if nil != err {
			return err
		}
	}

	// 个性化配置文件的配置
	for k := range mcnf.ICF {
		if k == cmdEnv {
			item := mcnf.ICF[k]
			if 0 < len(item.Files) {
				for _, iv := range item.Files {
					t, err := parsePath(iv)
					if nil != err {
						return err
					}

					// 支持扩展配置
					if 0 < len(item.Hosts) {
						for ik, iv := range item.Hosts {
							if "" != iv.Suffix && cmdKey == ik {
								fileFullName := t.dirPath + t.fileName + iv.Suffix + "." + t.fileExt
								if "" == t.dirPath {
									fileFullName = k + "/" + fileFullName
								}

								err = addConfFile(fileFullName, t.fileName)
								if nil != err {
									return err
								}
							}
						}
					} else {
						fileFullName := iv
						if "" == t.dirPath {
							fileFullName = k + "/" + fileFullName
						}

						err = addConfFile(fileFullName, t.fileName)
						if nil != err {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
