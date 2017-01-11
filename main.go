package main

import (
	"bytes"
	"flag"
	"fmt"
	hdfs "github.com/xmwilldo/gowfs"
	"os"
	"os/exec"
	"strings"
	//"io/ioutil"
)

var (
	HDFSHOST string
	HDFSPORT string
	HDFSUSER string
	BASEDIR  string
	HDFSURI  string
)

func init() {
	HDFSURI = os.Getenv("BSI_HDFS_HDFSDEMO_URI")
	HDFSHOST = os.Getenv("BSI_HDFS_HDFSDEMO_HOST")
	HDFSPORT = os.Getenv("BSI_HDFS_HDFSDEMO_PORT")
	HDFSUSER = os.Getenv("BSI_HDFS_HDFSDEMO_USERNAME")
	BASEDIR = os.Getenv("BSI_HDFS_HDFSDEMO_NAME")
	fmt.Printf("HDFSURI: %s\nHDFSHOST: %s\nHDFSPORT: %s\nHDFSUSER: %s\nBASEDIR: %s\n", HDFSURI, HDFSHOST, HDFSPORT, HDFSUSER, BASEDIR)

	HDFSUSER = changeUsername(HDFSUSER)
}

func main() {

	createdir := flag.String("mkdir", "", "create dir path.")
	lsdir := flag.String("ls", "", "list the dir.")
	load := flag.String("load", "", "load file.")
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

	//
	//
	//path := hdfs.Path{}
	//path.Name = "/abc"
	//ab, _ := fs.GetFileStatus(path)
	//fmt.Println(ab.PathSuffix, ab.Length)
	//return

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
	} else if *createdir != "" {
		isCreated, err := createDirectory(fs, *createdir, 0700)
		if err != nil {
			fmt.Println("createDirectory err:", err)
			return
		}

		fmt.Println("make dir: ", isCreated)
	} else if *load != "" {
		in := bytes.NewBuffer(nil)
		cmd := exec.Command("sh")
		cmd.Stdin = in

		//fmt.Println(*load)
		abc := strings.Split(*load, "/")
		createFile := abc[len(abc)-1]
		//fmt.Println(createFile)

		_, err := in.WriteString("curl -o /tmp/" + createFile + " " + *load + "\n")
		if err != nil {
			fmt.Println(err)
		}
		in.WriteString("exit\n")

		if err := cmd.Run(); err != nil {
			fmt.Println("cmd run err:", err)
			return
		}

		data, err := os.Open("/tmp/" + createFile)
		//data,err := ioutil.ReadFile("/home/wm/Desktop/test.tar.gz")
		//buffer := bytes.NewBuffer(data)

		ok, err := fs.Create(
			data,
			hdfs.Path{Name: "/" + createFile},
			false,
			0,
			0,
			0700,
			0,
		)
		if err != nil {
			fmt.Println("Create file err:", err)
			return
		}
		fmt.Println(ok)

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
	in.WriteString("curl -i -v --negotiate -u : -b /tmp/cookiejar.txt -c /tmp/cookiejar.txt " + HDFSURI + "?op=liststatus\n")
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
	config.Addr = HDFSHOST + ":" + HDFSPORT
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
