package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	lowerDomainName string
	lowerServiceName string
	upperServiceName string
)

var regex  = "[a-zA-Z]*"

func main(){
	var domainName, serviceName string
	flag.StringVar(&domainName,"d","","domain like web, srv, api")
	flag.StringVar(&serviceName,"s","","service name ")
	flag.Parse()
	if domainName == "" ||  serviceName == ""{
		fmt.Println("required parameter is missing")
		flag.Usage()
		os.Exit(1)
	}
	//校验输入的名称
	re,err := regexp.Compile(regex)
	matched := re.Match([]byte(domainName))

	if !matched {
		fmt.Println("domain name just accept alphabet")
		return
	}

	matched = re.Match([]byte(serviceName))
	if !matched{
		fmt.Println("service name just accept alphabet")
		return
	}


	//大小写转化
	lowerDomainName = strings.ToLower(domainName)
	lowerServiceName = strings.ToLower(serviceName)
	upperServiceName = strings.ToUpper(string(serviceName[0])) + string([]byte(serviceName)[1:])
	err = clone()
	if err != nil{
		fmt.Printf("%v\n",err)
	}

	err = replace(domainName,serviceName)

	if err != nil{
		fmt.Println("%v\n",err)
	}

}


func clone()error{

	cmd :=exec.Command("git","clone","https://github.com/dayan-be/template-service.git")
	//TODO:下载下来的文件夹后面多了个斜杠,需要找原因
	return cmd.Run()

}

func replace(domainName, serviceName string)error{
	//重命名文件夹
	absPath,err := filepath.Abs("./template-service/")
	if err != nil{
		return err
	}

	path := filepath.Dir(absPath)
	newPath := path +string(filepath.Separator)+domainName+"-"+serviceName
	err =os.Rename(absPath,newPath)
	if err != nil{
		return err
	}

	//删除git
	os.Chdir(newPath)
	err = os.RemoveAll(".git")
	if err != nil{
		return err
	}
	//修改文件内容
	recursionDirectory(newPath)
	return nil
}


//递归遍历目录
func recursionDirectory(path string) error{
	rd, err := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			//fmt.Printf("[%s]\n", path+string(filepath.Separator)+fi.Name())
			recursionDirectory(path+ string(filepath.Separator) + fi.Name())
		} else {
			replaceFileContent(path+string(filepath.Separator)+fi.Name())
		}
	}
	return err
}

//替换内容
func replaceFileContent(fileName string){
	f,err := ioutil.ReadFile(fileName)
	if err != nil{
		fmt.Println(err)
		return
	}

	strContent := string(f)
	strContent = strings.Replace(strContent,"$(DOMAIN)",lowerDomainName,-1)
	strContent = strings.Replace(strContent,"$(SERVICE)",upperServiceName,-1)
	strContent = strings.Replace(strContent, "$(LOWER-SERVICE)", lowerServiceName,-1)
	file , _ := os.Create(fileName)
	defer file.Close()
	file.WriteString(strContent)
}



