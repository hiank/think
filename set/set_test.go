package set_test

// import (
// 	"context"
// 	"path/filepath"
// 	"testing"

// 	"github.com/hiank/think/set"
// 	"gotest.tools/v3/assert"
// )

// func TestFilepathAbsSuffix(t *testing.T) {
// 	path1, _ := filepath.Abs("testdata/config")
// 	path, _ := filepath.Abs("testdata/config/")
// 	assert.Equal(t, path, path1)
// 	assert.Equal(t, rune(path[len(path)-1]), rune('g'))
// }

// func TestFilepathAbsNotExsit(t *testing.T) {
// 	_, err := filepath.Abs("notExsitFolder")
// 	assert.Assert(t, err == nil, "Abs并不会检查路径是否存在", err)
// 	// t.Log(path)
// }

// func TestArrayAppend(t *testing.T) {
// 	arr := make([]string, 0, 4)
// 	valArr := []string{
// 		"1",
// 		"2",
// 		"3",
// 		"4",
// 		"5",
// 	}
// 	arr = append(arr, valArr...)
// 	assert.Equal(t, len(arr), 5)
// 	assert.Equal(t, cap(arr), 8)
// }

// func TestModSetSignUpFolders(t *testing.T) {
// 	arr := []string{
// 		"testdata/folders/one",
// 		"testdata/folders/one/two",
// 		"testdata/folders/onepop/tmp/",
// 		"testdata/folders/onepop",
// 		"testdata/folders/notfolder",
// 	}
// 	err := set.Config.SignUpFolder(arr...)
// 	assert.Assert(t, err == nil, "如果存在合法的目录，不会返回错误")

// 	oms := set.NewOutModSet(set.Config)
// 	folders := oms.GetFolders()
// 	assert.Equal(t, len(folders), 2)

// 	path, _ := filepath.Abs(arr[0])
// 	assert.Equal(t, path, folders[0])

// 	path, _ = filepath.Abs(arr[3])
// 	assert.Equal(t, path, folders[1])

// 	set.ResetConfigMod()
// }

// func TestConfigModOnStart(t *testing.T) {
// 	set.Config.SignUpFolder("testdata/configs/deep")

// 	val := &testJsonData{}
// 	err := set.Config.SignUpValue(set.JSON, val)
// 	assert.Assert(t, err == nil, err)

// 	set.Config.OnStart(context.Background())
// 	assert.Equal(t, val.Value, "value 3")
// }
