package main

import (
	"archive/tar"
	"bufio"
	"github.com/hotei/dcompress"
	"io"
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
	Remove_duplicate int
}

const Extention = "tar.Z"
const DirectoryConfigFile = "config.json"
const FileListFile = "fileList.txt"

func checkInputDirecotry(config Configuration, fileListMap map[string]bool) {
	files, err := ioutil.ReadDir(config.Input_directory)
	if err != nil {
		log.Fatal(err)
	}
	existingFiles := make(map[string]string)

	for _, f := range files {
		fileName := f.Name()
		if strings.Contains(fileName, Extention) {
			log.Println("processing " + fileName)
			processCompressFile(filepath.Join(config.Input_directory, fileName), config, fileListMap, existingFiles)
		}
	}

}

func processCompressFile(fileName string, config Configuration, fileListMap map[string]bool, existingFiles map[string]string) {
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
				key := findKey(name, fileListMap)
				if key != "" {
					log.Println("found a file ", name)
					target := filepath.Join(config.Output_directory, getFileName(name))
					newFileName := getFileName(name)
					createFileFlag := true
					if config.Remove_duplicate == 1 {

						existingFileName := existingFiles[key]
						if existingFileName != "" {
							if existingFileName < newFileName {
								existingFiles[key] = newFileName
								existingTarget := filepath.Join(config.Output_directory, existingFileName)
								os.Remove(existingTarget)
								log.Println("removed file", existingFileName, "because", newFileName, "is newer")
							} else {
								createFileFlag = false
							}
						} else {
							existingFiles[key] = newFileName
						}
					}
					if createFileFlag {
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

func getFileName(fileName string) string {
	parts := strings.Split(fileName, "/")
	return parts[len(parts)-1]
}

func findKey(fileName string, fileListMap map[string]bool) string {
	for k, _ := range fileListMap {
		if strings.Contains(fileName, k) {
			return k
		}
	}
	return ""
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
		if !seen[line+"_1"] {
			seen[line+"_1"] = true
		}
		if !seen[line+"_2"] {
			seen[line+"_2"] = true
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

	log.Println("input in config file", configuration.Input_directory)
	log.Println("output in config file", configuration.Output_directory)
	log.Println("remove duplicate in config file", configuration.Remove_duplicate)
	createDir(configuration.Output_directory)
	fileListMap := readFileList(FileListFile)
	checkInputDirecotry(configuration, fileListMap)
}
