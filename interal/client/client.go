package client

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"updatePlus/interal/pkg/mode"
)

const snapshotFilePath = "snapshot.json"
const tagFileDefault = "tagFile"

func scanDir(tagFilePath string) (res map[string]mode.SimpleFileInfo, err error) {
	err = filepath.WalkDir(tagFilePath,
		func(path string, dirEntry os.DirEntry, err error) error {
			if err != nil {
				return errors.New(fmt.Sprint("读取文件目录时，发生异常->%s", err.Error()))
			}

			if !dirEntry.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return errors.New(fmt.Sprintf("读取目录文件%s时发生异常->%s", path, err.Error()))
				}

				shaHash := sha256.New()
				if _, err := io.Copy(shaHash, file); err != nil {
					return errors.New(fmt.Sprintf("计算文件%s哈希时发生异常->%s", path, err.Error()))
				}

				hashCode := shaHash.Sum(nil)
				fileInfo, err := file.Stat()
				if err != nil {
					return errors.New(fmt.Sprintf("获取文件%s的info属性时发生异常->%s", path, err.Error()))
				}

				res[path] = mode.SimpleFileInfo{
					Name:           file.Name(),
					Size:           fileInfo.Size(),
					HashCode:       fmt.Sprintf("%x", hashCode),
					LastUpdateTime: fileInfo.ModTime(),
				}
			}

			return nil
		})

	return res, err
}

func generaSnapshot(tagFilePath string) (err error) {
	simpleFileInfo, err := scanDir(tagFilePath)
	if err != nil {
		return errors.Join(
			errors.New(fmt.Sprintf("扫描文件目录%s时发生异常", tagFilePath)),
			err)
	}

	snapFile, err := os.Create("snapshot.json")
	if err != nil {
		return errors.New("创建snapshot.json时发生异常->" + err.Error())
	}

	jsonBytes, err := json.MarshalIndent(simpleFileInfo, "", "    ")
	if err != nil {

	}

	_, err = snapFile.Write(jsonBytes)
	if err != nil {

	}

	err = snapFile.Sync()
	if err != nil {

	}
	snapFile.Close()
	return nil
}

func Client() {
	startTime := time.Now()
	fmt.Printf("客户端启动->%v\n", startTime)

	pTagFile := flag.String("tagFile", tagFileDefault, "目标文件目录")
	flag.Parse()

	// 若文件不存在则返回
	tagFileInfo, err := os.Stat(*pTagFile)
	if err != nil {
		fmt.Printf("打开目标文件失败，发生错误%s", err.Error())
		return
	}

	// 若目标文件不是一个目录则退出
	if !tagFileInfo.IsDir() {
		fmt.Println("目标文件不是一个目录")
		return
	}

	// 尝试找Snapshot.json文件
	_, err = os.Stat(snapshotFilePath)
	if err != nil {
		fmt.Printf("打开Snapshot.json文件失败，重新执行生成程序")
		generaSnapshot(*pTagFile)
	}

	defer func() {
		fmt.Printf("程序运行结束，耗时->%d\n", time.Since(startTime))
	}()

}
