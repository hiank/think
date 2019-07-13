package conf_test

import (
	"time"
	"context"
	"sync"
	// "reflect"
	"github.com/hiank/think/conf"
	// "encoding/json"
	// "io/ioutil"
	"testing"
)

type TestJson struct {
	Wei  string `json:"testJson.wei"`
	Shao string `json:"testJson.shao"`
}

func TestChan(t *testing.T) {

	wait := new(sync.WaitGroup)
	wait.Add(2)
	ch := make(chan int)
	// idx := 0
	fun := func() {

		// t.Logf("idx : %d\n", idx)
		// idx++
		// num := <-ch
		// t.Logf("%d\n", num)
		select {
		case <-ch:
			t.Logf("read ch")
		}
		wait.Done()
	}

	go fun()
	go fun()

	close(ch)
	// ch <- 1
	// ch <- 2

	wait.Wait()
}

func TestConf(t *testing.T) {

	list := make([]*conf.Info, 0, 1)
	if info, err := conf.NewInfoByFile("conf_test.json", &TestJson{}); err == nil {

		list = append(list, info)
	}
	conf.Init(list...)

	item := conf.Get("sys")
	t.Logf("key : %v __ num : %v\n", item.ValueByName("key"), item.ValueByName("maxMessageGo"))

	item = conf.Get("testJson")
	t.Logf("%v __ %v\n", item.ValueByName("shao"), item.ValueByName("wei"))
}

func TestContext(t *testing.T) {

	ctx1, cancel := context.WithCancel(context.Background())
	ctx2, _ := context.WithCancel(ctx1)

	go func() {

		select {
		case <-ctx2.Done():
			t.Logf("done 2")
		case <-ctx1.Done():
			t.Logf("done 1")
		}
	}()

	
	time.Sleep(10000)
	cancel()

	time.Sleep(10000)

	// go func(ctx context.Context) {

	// 	ctx1, cancel := context.WithCancel(ctx)
	// 	go func()
	// }(ctx1)
}

func TestGoroutine(t *testing.T) {

	// wait := make(chan bool)

	go func() {

		// close(wait)		
		t.Log("wait")
	}()

	// <-wait
	t.Log("main")
}


// func TestMoreLoad(t *testing.T) {

// 	buf, err := ioutil.ReadFile("conf.json")
// 	if err != nil {
// 		t.Log("read conf.json error : " + err.Error())
// 	}

// 	sys := conf.Sys{}
// 	v := &sys

// 	err = json.Unmarshal(buf, v)
// 	if err != nil {
// 		t.Log("unmarshal error : " + err.Error())
// 	}

// 	rv := reflect.ValueOf(sys)
// 	rt := reflect.TypeOf(sys)

// 	// for i, num := 0, rt.NumField(); i < num; i++ {

// 	// 	f := rt.Field(i)
// 	// 	t.Logf("%s\n", f.Name)
// 	// }
// 	if _, ok := rt.FieldByName("Key"); ok {

// 		t.Logf("rv key : %v\n", rv.FieldByName("Key").Interface())
// 	} else {

// 		t.Log("no filed")
// 	}

// 	t.Logf("key : %s, num : %d\n", v.Key, v.MaxMessageGo)

// 	buf, err = ioutil.ReadFile("conf_test.json")
// 	if err != nil {
// 		t.Log("read conf_test.json error : " + err.Error())
// 	}

// 	err = json.Unmarshal(buf, v)
// 	if err != nil {
// 		t.Log("unmarshal error : " + err.Error())
// 	}

// 	t.Logf("key : %s, num : %d\n", v.Key, v.MaxMessageGo)

// 	// rt = reflect.TypeOf(sys)
// 	rv = reflect.ValueOf(sys)
// 	if _, ok := rt.FieldByName("Key"); ok {

// 		t.Logf("rv key : %v\n", rv.FieldByName("Key").Interface())
// 	} else {

// 		t.Log("no filed")
// 	}

// }
