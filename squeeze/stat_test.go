package squeeze

import (
	"bufio"
	"os"
	"testing"
)

var (
	tempDir = os.TempDir()
	dirTestStat = tempDir + "/test_stat"
	fileNameInputTest = dirTestStat + "/TestInput.txt"
	fileNameOutputTest = dirTestStat + "/TestOutput.txt"
	fileNameInputEmptyTest = dirTestStat + "/TestInputEmpty.txt"
	fileNameInputEmptyOtherTest = dirTestStat + "/TestInputEmptyOther.txt"
	fileNameInputWithoutStatTest = dirTestStat + "/TestInputWithoutStat.txt"
	fileNameOutputWithoutStatTest = dirTestStat + "/TestOutputWithoutStat.txt"
	fileNameInputErrorTest = dirTestStat + "/TestInputError.txt"
	fileNameOutputErrorTest = dirTestStat + "/TestOutputError.txt"
	string1 = "This is a test string"
	string2 = "This is a another test string"
)

func TestStat(t *testing.T) {
	setUpStat()

	t.Run("InputNoExist", func(t *testing.T) {
		_, err := GetMapStat("FileName", 0, "")
		if !os.IsNotExist(err) {
			t.Errorf("File exist")
		}
	})

	t.Run("GetMapStatEmpty", func(t *testing.T) {
		paths := []string {
			fileNameInputEmptyTest, fileNameInputEmptyOtherTest,
		}
		for _, path := range paths {
			stat, err := GetMapStat(path, 15, "Jan _2 15:04:05")

			if os.IsNotExist(err) {
				t.Errorf("File not exist")
			}

			if len(stat) != 0 {
				t.Fail()
			}
		}
	})

	t.Run("GetMapStat", func(t *testing.T) {
		type DataStruct struct {
			input string
			repeat1, index1, repeat2, index2 uint
		}
		arrayDataStruct := []DataStruct {
			{fileNameInputTest, 2, 1, 1, 2},
			{fileNameInputWithoutStatTest, 1, 0, 1, 1},
		}

		for _, data := range arrayDataStruct {
			stat, err := GetMapStat(data.input, 15, "Jan _2 15:04:05")

			if os.IsNotExist(err) {
				t.Errorf("File not exist")
			}

			if len(stat) != 2 {
				t.Errorf("Length not equals 2")
			}

			if stat1, ok := stat[string1]; ok {
				if !(stat1.repeat == data.repeat1 && stat1.index == data.index1) {
					t.Fail()
				}
			} else {
				t.Fail()
			}

			if stat2, ok := stat[string2]; ok {
				if !(stat2.repeat == data.repeat2 && stat2.index == data.index2) {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		}
	})

	t.Run("OutputNoExist", func(t *testing.T) {
		err := ReturnResult("", nil)
		if !os.IsNotExist(err) {
			t.Errorf("File exist")
		}
	})

	t.Run("ReturnResult", func(t *testing.T) {
		type DataStruct struct {
			input, output string
			len int
		}
		arrayDataStruct := []DataStruct {
			{fileNameInputTest, fileNameOutputTest, 6},
			{fileNameInputWithoutStatTest, fileNameOutputWithoutStatTest, 2},
		}

		for _, data := range arrayDataStruct {
			stat, _ := GetMapStat(data.input, 15, "Jan _2 15:04:05")
			err := ReturnResult(data.output, stat)

			if os.IsNotExist(err) {
				t.Errorf("File not exist")
			}

			if getLenFile(data.output) != data.len {
				t.Fail()
			}
		}
	})

	t.Run("ReturnResultError", func(t *testing.T) {
		stat, _ := GetMapStat(fileNameInputErrorTest, 15, "Jan _2 15:04:05")
		err := ReturnResult(fileNameOutputErrorTest, stat)

		if err.Error() != "empty result file" {
			t.Errorf("Error")
		}

		_, err = os.Stat(fileNameOutputErrorTest)
		if !os.IsNotExist(err) {
			t.Fail()
		}
	})

	tearDownStat()
}

func setUpStat() {
	_ = os.Mkdir(dirTestStat,0777)

	fileInput, err1 := os.OpenFile(fileNameInputTest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err1 == nil {
		_, _ = fileInput.WriteString("Oct 19 01:20:58 " + string1 + "\n")
		_, _ = fileInput.WriteString("Oct 19 01:21:50 " + string1 + "\n")
		_, _ = fileInput.WriteString("Oct 19 12:45:53 " + string2)

		defer fileInput.Close()
	}

	fileInputEmpty, err2 := os.OpenFile(fileNameInputEmptyTest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err2 == nil {
		defer fileInputEmpty.Close()
	}

	fileInputEmptyOther, err3 := os.OpenFile(fileNameInputEmptyOtherTest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err3 == nil {
		_, _ = fileInputEmptyOther.WriteString(string1)

		defer fileInputEmptyOther.Close()
	}

	fileInputWithoutStat, err4 := os.OpenFile(fileNameInputWithoutStatTest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err4 == nil {
		_, _ = fileInputWithoutStat.WriteString("Oct 19 01:20:58 " + string1 + "\n")
		_, _ = fileInputWithoutStat.WriteString("Oct 19 12:45:53 " + string2)

		defer fileInputWithoutStat.Close()
	}

	fileInputError, err5 := os.OpenFile(fileNameInputErrorTest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err5 == nil {
		_, _ = fileInputError.WriteString("Oct 9 01:20:58 " + string1)

		defer fileInputError.Close()
	}
}

func tearDownStat() {
	_ = os.RemoveAll(dirTestStat)
}

func getLenFile(path string) int {
	file, _ := os.Open(path)
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}