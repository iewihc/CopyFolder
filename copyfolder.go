package GO_CopyFolder

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type copyer struct {
	dest string
	wg sync.WaitGroup
}

func Copy(src, dest string)  {
	c := &copyer{dest: dest}
	c.wg.Add(1)

	go c.walk(src)

	c.wg.Wait()
}

func (c *copyer) walk(folderPath string) {
	files,err:=ioutil.ReadDir(folderPath)
	if err!=nil{
		log.Fatal(err)
	}
	for _,f := range files{
		if c.isFolder(f) {
			c.wg.Add(1)
			go c.walk(folderPath+"/"+f.Name())
			continue
		}
		os.MkdirAll(c.dest+"/"+folderPath,os.ModePerm)
		fmt.Println("輸出一下"+c.dest+"/"+folderPath)
		c.wg.Add(1)
		go c.copy(folderPath+"/"+f.Name(),c.dest+"/"+folderPath+"/"+f.Name())

	}
	c.wg.Done()
}

func (c *copyer) isFolder(f os.FileInfo) bool {
	switch mode := f.Mode();{
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func (c *copyer) copy(src,dst string)  {
	in,err:=os.Open(src)
	if err!=nil {
		log.Fatal(err)
	}
	defer in.Close()
	out,err:=os.Create(dst)
	if err!=nil {
		log.Fatal(err)
	}
	defer out.Close()

	_,err =io.Copy(out,in)
	if err!=nil {
		log.Fatal(err)
	}
	out.Close()
	c.wg.Done()
}

