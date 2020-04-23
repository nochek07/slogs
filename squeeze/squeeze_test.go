package squeeze

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	flags = Flags{
		ExtFlag:         ".txt",
		DatePatternFlag: "Jan _2 15:04:05",
		RecursionFlag:   false,
		RemoveFlag:      false,
		DateLengthFlag:  15,
		SizeFileAsync:   1,
	}
	dirTestSqueeze = tempDir + "/test_squeeze/"
	fileName1 = "test1.txt"
	fileName2 = "test2.txt"
	fileName3 = "test3.txt"
	fullFileName1 = dirTestSqueeze + fileName1
	fullFileName2 = dirTestSqueeze + fileName2
	fullFileName3 = dirTestSqueeze + fileName3
)

func TestSqueeze(t *testing.T) {
	setUpSqueeze()

	t.Run("findFilesErrorPath", func(t *testing.T) {
		_, err := findFiles(":", flags)
		if err.Error() != "path error" {
			t.Errorf("Error")
		}
		Squeeze(":", flags)
	})

	t.Run("findFilesOncePath", func(t *testing.T) {
		files, err := findFiles(fullFileName1, flags)
		if err != nil {
			t.Errorf("Error")
		}
		fmt.Println(files)
		if !(len(files) == 1 && files[0] == fullFileName1) {
			t.Fail()
		}
	})

	t.Run("Squeeze", func(t *testing.T) {
		Squeeze(dirTestSqueeze, flags)

		newFullFileName1 := dirTestSqueeze + prefix + fileName1
		newFullFileName2 := dirTestSqueeze + prefix + fileName2
		_, err1 := os.Stat(newFullFileName1)
		_, err2 := os.Stat(newFullFileName2)
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			t.Fail()
		}

		if getLenFile(newFullFileName1) != 6 || getLenFile(newFullFileName2) != 7 {
			t.Fail()
		}

		newFullFileName3 := dirTestSqueeze + prefix + fileName3
		_, err3 := os.Stat(newFullFileName3)
		if !os.IsNotExist(err3) {
			t.Fail()
		}
	})

	t.Run("SqueezeRemove", func(t *testing.T) {
		flags.RecursionFlag = true
		flags.RemoveFlag = true
		Squeeze(dirTestSqueeze, flags)

		_, err1 := os.Stat(fullFileName1)
		_, err2 := os.Stat(fullFileName2)
		if !os.IsNotExist(err1) || !os.IsNotExist(err2) {
			t.Fail()
		}

		_, err3 := os.Stat(fullFileName3)
		if os.IsNotExist(err3) {
			t.Fail()
		}
	})

	tearDownSqueeze()
}

func setUpSqueeze() {
	_ = os.Mkdir(dirTestSqueeze,0777)

	fileInput1, err1 := os.OpenFile(fullFileName1, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err1 == nil {
		_, _ = fileInput1.WriteString("Oct 19 01:20:58 " + string1 + "\n")
		_, _ = fileInput1.WriteString("Oct 19 01:21:50 " + string1 + "\n")
		_, _ = fileInput1.WriteString("Oct 19 12:45:53 " + string2)

		defer fileInput1.Close()
	}

	fileInput2, err2 := os.OpenFile(fullFileName2, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err2 == nil {
		seconds := 0
		for i := 0; i < 20000; i++ {
			_, _ = fileInput2.WriteString(getDate(&seconds) + " " + string1 + "\n")
			_, _ = fileInput2.WriteString(getDate(&seconds) + " " + string2 + "\n")
		}
		defer fileInput2.Close()
	}

	fileInput3, err3 := os.OpenFile(fullFileName3, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err3 == nil {
		_, _ = fileInput3.WriteString("Oct 9 01:20:58 " + string1)

		defer fileInput3.Close()
	}
}

func tearDownSqueeze() {
	_ = os.RemoveAll(dirTestSqueeze)
}

func getDate(seconds *int) string {
	date := time.Date(2020, time.April, 1, 10, 0, *seconds, 0, time.UTC)
	*seconds++
	return date.Format("Jan 02 15:04:05")
}