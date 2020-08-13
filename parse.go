package configurer

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"net"
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
	ConfFile string              `json:"conf_file,omitempty" yaml:"conf_file,omitempty" toml:"conf_file,omitempty"`
	Files    []string            `json:"files,omitempty" yaml:"files,omitempty" toml:"files,omitempty"`
	Hosts    map[string]hostItem `json:"hosts,omitempty" yaml:"hosts,omitempty" toml:"hosts,omitempty"`
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

func getIps() error {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return err
	}

	ips = make([]string, 0)
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return nil
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

func resetConf() {
	cnf.ConfFile = nil
	cnf.CommonConfItem = nil
	
	cnf.ConfFile = make(map[string]string, 0)
}

// ParseConf 解析配置的配置
func ParseConf(fname string, reload bool) error {
	if false == reload {
		flag.Parse()
	}

	// 重置配置
	resetConf()

	err := getIps()
	if nil != err {
		return err
	}

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
			// 解析conf-file并设置到common-conf-item中
			if "" != item.ConfFile {
				t, err := parsePath(item.ConfFile)
				if nil != err {
					return err
				}

				fileFullName := t.dirPath + t.fileName + "." + t.fileExt
				if "" == t.dirPath {
					fileFullName = k + "/" + fileFullName
				}
				tfi, err := os.Open(fileFullName)
				if err != nil {
					return err
				}
				defer tfi.Close()
				td, err := edcoder.NewDecoder(edcoder.SetDecoderExt(t.fileExt), edcoder.SetDecoderReader(tfi))
				if nil != err {
					return err
				}
				var c map[string]interface{}
				err = td.Decode(&c)
				if nil != err {
					return err
				}
				for ck, cv := range c {
					cnf.CommonConfItem[ck] = cv
				}
			}
			
			if 0 < len(item.Files) {
				// 解析files和hosts
				for _, iv := range item.Files {
					t, err := parsePath(iv)
					if nil != err {
						return err
					}

					// 支持扩展配置
					if 0 < len(item.Hosts) {
						if "" == cmdKey {
							var IP string
							var HI hostItem
							for _, iv := range ips {
								if "" == IP {
									for iik := range item.Hosts {
										if iv == iik {
											IP = iv
											HI = item.Hosts[iik]
											break
										}
									}
								}
								break
							}

							if "" == IP { // 无匹配则默认基本配置文件
								fileFullName := iv
								if "" == t.dirPath {
									fileFullName = k + "/" + fileFullName
								}

								err = addConfFile(fileFullName, t.fileName)
								if nil != err {
									return err
								}
							} else { // 扩展配置
								fileFullName := t.dirPath + t.fileName + HI.Suffix + "." + t.fileExt
								if "" == t.dirPath {
									fileFullName = k + "/" + fileFullName
								}

								err = addConfFile(fileFullName, t.fileName)
								if nil != err {
									return err
								}
							}
						} else { // 扩展配置
							flag := false
							for ik, iv := range item.Hosts {
								if "" != iv.Suffix && cmdKey == ik {
									flag = true
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

							if false == flag {
								return errors.New("no matched cmdKey")
							}
						}
					} else { // 则默认基本配置文件
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
