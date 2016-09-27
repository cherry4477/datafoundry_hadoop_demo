package main

import (
	"fmt"
	hdfs "github.com/xmwilldo/gowfs"
	"os"
	"bytes"
	"os/exec"
	"flag"
	"strings"
)

var (
	HDFSHOST string
	HDFSPORT string = ":50070"
	HDFSUSER string
	BASEDIR string

)

func init() {
	HDFSHOST = os.Getenv("BSI_HDFS_HDFSDEMO_HOST")
	HDFSUSER = os.Getenv("BSI_HDFS_HDFSDEMO_USERNAME")
	BASEDIR = os.Getenv("BSI_HDFS_HDFSDEMO_NAME")
	fmt.Printf("HDFSHOST: %s\nHDFSPORT: %s\nHDFSUSER: %s\nBASEDIR: %s\n", HDFSHOST, HDFSPORT, HDFSUSER, BASEDIR)

	HDFSUSER = changeUsername(HDFSUSER)
}

func main() {

	createdir := flag.String("mkdir", "default", "create dir path.")
	lsdir := flag.String("ls", "", "list the dir.")
	flag.Parse()

	err := initCookie()
	if err != nil {
		fmt.Println("initCookie err:", err)
		return
	}

	config := newHdfsConfig()

	fs, err := hdfs.NewFileSystem(*config)
	if err != nil {
		fmt.Println("NewFileSystem err:", err)
		return
	}

	if *lsdir != "" {
		path := hdfs.Path{}
		path.Name = "/"
		files, err := fs.ListStatus(path)
		if err != nil {
			fmt.Println("GetContentSummary err: ", err)
			return
		}

		for _, file := range files {
			fmt.Println(file.PathSuffix)
		}
	} else {
		isCreated, err := createDirectory(fs, *createdir, 0700)
		if err != nil {
			fmt.Println("createDirectory err:", err)
			return
		}

		fmt.Println("make dir: ", isCreated)
	}

}

func initCookie() error {
	isExist := isExistfile("/tmp/cookiejar.txt")
	if isExist {
		err := os.Remove("/tmp/cookiejar.txt")
		if err != nil {
			fmt.Println("Remove cookie err:", err)
			return err
		}
	}

	in := bytes.NewBuffer(nil)
	cmd := exec.Command("sh")
	cmd.Stdin = in
	//execstr := "curl -i -v --negotiate -u : -b /tmp/cookiejar.txt -c /tmp/cookiejar.txt http://hadoop-1.jcloud.local:50070/webhdfs/v1"+BASEDIR+"?op=liststatus\n"
	in.WriteString("curl -i -v --negotiate -u : -b /tmp/cookiejar.txt -c /tmp/cookiejar.txt http://hadoop-1.jcloud.local:50070/webhdfs/v1"+BASEDIR+"?op=liststatus\n")
	//in.WriteString("curl -i -v --negotiate -u : -b /tmp/cookiejar.txt -c /tmp/cookiejar.txt http://hadoop-1.jcloud.local:50070/webhdfs/v1?op=liststatus\n")
	in.WriteString("exit\n")

	//fmt.Println("execstr:", execstr)
	if err := cmd.Run(); err != nil {
		fmt.Println("cmd run err:", err)
		return err
	}
	return nil
}

func isExistfile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func newHdfsConfig() *hdfs.Configuration {
	config := hdfs.NewConfiguration()
	config.Addr = HDFSHOST+HDFSPORT
	config.User = HDFSUSER
	config.BasePath = BASEDIR
	config.MaxIdleConnsPerHost = 64

	return config
}

func createDirectory(fs *hdfs.FileSystem, name string, fileMode os.FileMode) (bool, error) {
	path := hdfs.Path{}
	path.Name = name

	isCreated, err := fs.MkDirs(path, fileMode)
	if err != nil {
		return isCreated, err
	}

	return isCreated, nil
}

func changeUsername(username string) string {
	return strings.Split(username, "@")[0]
}
