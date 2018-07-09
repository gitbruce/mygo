package main

import (
	"archive/tar"
	"github.com/hotei/dcompress"
	"io"
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Input_directory  string
	Output_directory string
}

const Extention = "tar.Z"
const DirectoryConfigFile = "config.json"
const FileListFile = "fileList.txt"

func checkInputDirecotry(inputDir string, ouputDir string, fileListMap map[string]bool) {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()
		if strings.Contains(fileName, Extention) {
			log.Println("processing " + fileName)
			processCompressFile(filepath.Join(inputDir, fileName), ouputDir, fileListMap)
		}
	}

}

func processCompressFile(fileName string, ouputDir string, fileListMap map[string]bool) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	zfile, err := dcompress.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	tarReader := tar.NewReader(zfile)
	i := 0
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			{
//				log.Println("(", i, ")", "Name: ", name, " with size: ", header.Size)
				if (fileListMap[name]) {
					log.Println("found a file ", name)
					target := filepath.Join(ouputDir, name)
					outFile, err := os.Create(target)
					if err != nil {
						log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
					}
					defer outFile.Close()
					if _, err := io.Copy(outFile, tarReader); err != nil {
						log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
					}
				} 
			}

		default:
			log.Println("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
}

func readFileList(fileName string) map[string]bool {
	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer file.Close()
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
		}
	}
	log.Println("loaded file list ", len(seen))
	return seen
}

func createDir(outputDir string) {
	os.MkdirAll(outputDir, os.ModePerm)
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf(DirectoryConfigFile, &configuration)
	if err != nil {
		log.Println("error in configuration file with: " + err.Error())
		os.Exit(500)
	}

	log.Println("input: " + configuration.Input_directory)
	log.Println("output: " + configuration.Output_directory)
	createDir(configuration.Output_directory)
	fileListMap := readFileList(FileListFile)
	checkInputDirecotry(configuration.Input_directory, configuration.Output_directory, fileListMap)
}
