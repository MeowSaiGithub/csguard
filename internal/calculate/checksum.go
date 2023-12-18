package calculate

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const largeFileSizeThreshold = 100 * 1024 * 1024 // 100Mb

type Checksum interface {
	SetInputFile() *string
	SetInputFolder() *string
	SetOutputFile() *string
	CalculateInputValidation() error
	CalculateChecksum() error
	GetOutputFile() *string
	CreateOutput() error
	GetChecksum() *map[string]string
	SetChecksumFolder() *string
	ValidateInputValidation(cs string) error
	ValidateChecksum() error
	CreateValidateOutputTxt() error
	GetValidation() *map[string]string

	calculateSmallMd5(file string) (string, error)
	calculateLargeMd5(file string) (string, error)
	calculateMD5() error
	getAlgorithm() algo
	//SetAlgorithm() *string
	validateChecksumMd5() error
	loadChecksumFromFile() error
}

type checksum struct {
	mode                string
	inputFile           string
	algorithm           string
	inputFolder         string
	output              string
	checksum            map[string]string
	checksumFile        string
	validateMap         map[string]string
	validateChecksumMap map[string]string
}

type algo int

const (
	md5Algorithm algo = iota
)

func NewChecksumProvider() Checksum {
	return &checksum{
		checksum:            make(map[string]string),
		validateMap:         make(map[string]string),
		validateChecksumMap: make(map[string]string),
	}
}

func (c *checksum) calculateSmallMd5(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	cs := md5.Sum(data)
	return hex.EncodeToString(cs[:]), nil
}

func (c *checksum) calculateLargeMd5(file string) (string, error) {
	data, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = data.Close(); err != nil {
			log.Printf("failed to close file: %v\n", err)
		}
	}()
	hash := md5.New()
	if _, err = io.Copy(hash, data); err != nil {
		return "", err
	}
	cs := hex.EncodeToString(hash.Sum(nil))
	return cs, nil
}

func (c *checksum) ValidateMd5(checksum, input string) bool {
	return checksum == input
}

func (c *checksum) SetInputFile() *string {
	return &c.inputFile
}

func (c *checksum) SetInputFolder() *string {
	return &c.inputFolder
}

func (c *checksum) SetChecksumFolder() *string {
	return &c.checksumFile
}

func (c *checksum) SetOutputFile() *string {
	return &c.output
}
func (c *checksum) GetOutputFile() *string {
	return &c.output
}

//func (c *checksum) SetAlgorithm() *string {
//	return &c.algorithm
//}

func (c *checksum) getAlgorithm() algo {
	switch c.algorithm {
	case "md5":
		return md5Algorithm
	case "ssa":
		return 2
	default:
		return md5Algorithm
	}
}

func (c *checksum) calculateMD5() error {
	if c.inputFolder != "" {
		err := filepath.WalkDir(c.inputFolder, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			fileinfo, err := d.Info()
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			fp := filepath.Join(filepath.Dir(c.inputFolder), path)

			if fileinfo.Size() > largeFileSizeThreshold {
				c.checksum[fp], err = c.calculateLargeMd5(path)
				return err
			}
			c.checksum[fp], err = c.calculateSmallMd5(path)
			return err
		})
		return err
	}

	fileinfo, err := os.Stat(c.inputFile)
	if err != nil {
		return err
	}
	if fileinfo.Size() > largeFileSizeThreshold {
		c.checksum[c.inputFile], err = c.calculateLargeMd5(c.inputFile)
	} else {
		c.checksum[c.inputFile], err = c.calculateSmallMd5(c.inputFile)
	}
	return err
}

func (c *checksum) CalculateChecksum() error {
	switch c.getAlgorithm() {
	case md5Algorithm:
		return c.calculateMD5()
	case 2:
		return nil
	default:
		return c.calculateMD5()
	}
}

func (c *checksum) CreateOutput() error {
	ext := strings.Split(c.output, ".")
	switch ext[len(ext)-1] {
	case "txt":
		return createOutputTxt(c.output, &c.checksum)
	case "json":
		return createOutputJson(c.output, c.checksum)
	case "yaml":
		return createOutputYaml(c.output, c.checksum)
	default:
		return fmt.Errorf("invalid output format. supported format 'txt', 'json', 'yaml'")
	}
}

func createOutputTxt(fileName string, data *map[string]string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("failed to close: %v\n", err)
		}
	}()
	writer := bufio.NewWriter(f)
	for file, cs := range *data {
		str := fmt.Sprintf("%s %s\n", file, cs)
		_, err = writer.WriteString(str)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (c *checksum) CreateValidateOutputTxt() error {
	ext := strings.Split(c.output, ".")
	switch ext[len(ext)-1] {
	case "txt":
		return createOutputTxt(c.output, &c.validateMap)
	case "json":
		return createOutputJson(c.output, c.validateMap)
	case "yaml":
		return createOutputYaml(c.output, c.validateMap)
	default:
		return fmt.Errorf("invalid output format. supported format 'txt', 'json', 'yaml'")
	}
}

func (c *checksum) GetChecksum() *map[string]string {
	return &c.checksum
}

func (c *checksum) GetValidation() *map[string]string {
	return &c.validateMap
}

func (c *checksum) ValidateChecksum() error {
	if c.checksumFile != "" {
		if err := c.loadChecksumFromFile(); err != nil {
			return err
		}
	}
	err := c.validateChecksumMd5()
	return err
}

func (c *checksum) loadChecksumFromFile() error {
	ext := strings.Split(c.checksumFile, ".")
	var err error
	switch ext[len(ext)-1] {
	case "txt":
		c.validateChecksumMap, err = loadFromTxt(c.checksumFile)
	case "json":
		c.validateChecksumMap, err = loadFromJson(c.checksumFile)
	case "yaml":
		c.validateChecksumMap, err = loadFromYaml(c.checksumFile)
	default:
		return fmt.Errorf("invalid input format. supported format 'txt', 'json', 'yaml'")
	}
	return err
}

func loadFromTxt(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	data := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		} else {
			return nil, fmt.Errorf("invalid line: %s\n", line)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return data, err
}

func loadFromJson(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data []fileStruct
	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	mdata := make(map[string]string)
	for _, d := range data {
		mdata[d.FileName] = d.Value
	}

	return mdata, err
}

func loadFromYaml(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data []fileStruct
	if err = yaml.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	mdata := make(map[string]string)
	for _, d := range data {
		mdata[d.FileName] = d.Value
	}
	return mdata, err
}

func (c *checksum) validateChecksumMd5() error {
	if len(c.validateChecksumMap) != 0 {
		for f, _ := range c.validateChecksumMap {

			f = filepath.Clean(f)

			fileinfo, err := os.Stat(f)
			if err != nil {
				if os.IsNotExist(err) {
					c.validateMap[f] = "missing"
					continue
				}
				return err
			}

			if fileinfo.IsDir() {
				return nil
			}

			if fileinfo.Size() > largeFileSizeThreshold {
				c.checksum[f], err = c.calculateLargeMd5(f)
			} else {
				c.checksum[f], err = c.calculateSmallMd5(f)
			}
			if err != nil {
				return err
			}
			if checkChecksum(c.validateChecksumMap[f], c.checksum[f]) {
				c.validateMap[f] = "OK"
				continue
			}
			c.validateMap[f] = "NOT OK"
		}
		return nil
	}
	return nil
}

func checkChecksum(input, cs string) bool {
	return input == cs
}

func createOutputJson(filename string, mdata map[string]string) error {
	var data []fileStruct
	for key, value := range mdata {
		data = append(data, fileStruct{
			FileName: key,
			Value:    value,
		})
	}
	jdata, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jdata, 0644)
	return err
}

func createOutputYaml(filename string, mdata map[string]string) error {
	var data []fileStruct
	for key, value := range mdata {
		data = append(data, fileStruct{
			FileName: key,
			Value:    value,
		})
	}
	ydata, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, ydata, 0644)
	return err
}

type fileStruct struct {
	FileName string `json:"file_name" yaml:"file_name"`
	Value    string `json:"value,omitempty" yaml:"value,omitempty"`
}

func (c *checksum) CalculateInputValidation() error {
	if c.inputFile == "" && c.inputFolder == "" {
		return fmt.Errorf("either --input-file or --input-folder must be set")
	}
	if c.inputFile != "" {
		fileInfo, err := os.Stat(c.inputFile)
		if err != nil {
			return fmt.Errorf("error checking input file: %s", err)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("--input-file must be a file, not a directory. use --input-folder for directory")
		}
	}
	if c.inputFolder != "" {
		fileInfo, err := os.Stat(c.inputFile)
		if err != nil {
			return fmt.Errorf("error checking input folder: %s", err)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("--input-folder must be a folder, not a file. use --input-file for single file")
		}
	}
	return nil
}

func (c *checksum) ValidateInputValidation(cs string) error {
	if c.inputFile == "" && c.checksumFile == "" {
		return fmt.Errorf("either --input-file or --checksum-file must be set")
	}
	if c.inputFile != "" {
		if cs == "" {
			return fmt.Errorf("--checksum is empty")
		}
		fileInfo, err := os.Stat(c.inputFile)
		if err != nil {
			return fmt.Errorf("error checking input file: %s", err)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("--input-file must be a file, not a directory")
		}
		c.validateChecksumMap[c.inputFile] = cs
	}
	if c.checksumFile != "" {
		fileInfo, err := os.Stat(c.checksumFile)
		if err != nil {
			return fmt.Errorf("error checking checksum-file: %s", err)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("--checksum-file must be a file, not a directory")
		}
	}
	return nil
}