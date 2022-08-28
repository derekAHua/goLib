package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"testing"
	"time"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/28 09:25
// @Version 1.0

type Student struct {
	StudentID uint64
	Name      string
	Age       uint64
	Phone     string
}

func Test_defaultSafeErrGroup_Run(t *testing.T) {
	ctx := &gin.Context{}

	a := []Student{
		{StudentID: 100000048307},
		{StudentID: 100000048296},
		{StudentID: 100000048211},
		{StudentID: 100000048293},
		{StudentID: 100000048292},
		{StudentID: 100000048291},
		{StudentID: 100000048290},
		{StudentID: 100000048269},
		{StudentID: 100000048279},
		{StudentID: 100000048275},
	}

	t.Log(fmt.Sprintf("开始Goroutine数量:%v", runtime.NumGoroutine()))

	startTime := time.Now()
	studentMap, err := getStudentMap(ctx, a)
	if err != nil {
		t.Log(err)
	}

	t.Log(fmt.Sprintf("耗费时间：%v", time.Now().Sub(startTime).Milliseconds()))
	t.Logf("结果长度：%d", len(studentMap))

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 100)
		t.Log(fmt.Sprintf("结束Goroutine数量:%v", runtime.NumGoroutine()))
	}
}

func getStudentMap(ctx *gin.Context, studentList []Student) (ret map[uint64]Student, err error) {
	ret = make(map[uint64]Student, len(studentList))

	ids := make([]uint64, 0, len(studentList))
	for _, v := range studentList {
		ids = append(ids, v.StudentID)
	}

	g := NewErrGroup(ctx)

	mName := make(map[uint64]string, len(ids))
	g.Go(func() error {
		time.Sleep(time.Millisecond * 100)
		for _, v := range ids {
			mName[v] = fmt.Sprintf("name_%d", v)
		}
		return nil
	})

	mAge := make(map[uint64]uint64, len(ids))
	g.Go(func() error {
		time.Sleep(time.Millisecond * 200)
		for _, v := range ids {
			mAge[v] = v
		}
		//panic("出错啦")
		return nil
	})

	mPhone := make(map[uint64]string, len(ids))
	g.Go(func() error {
		time.Sleep(time.Millisecond * 300)
		for _, v := range ids {
			mPhone[v] = fmt.Sprintf("phone_%d", v)
		}
		return nil
	})

	g.Run(studentList,
		func(item interface{}) (ret interface{}, err error) { // 入参item
			ret, err = f(ctx, item.(Student)) // 获取学生信息
			return
		}, func(result interface{}) {
			stu := result.(Student)
			ret[stu.StudentID] = stu
		}, // WithWorkers(3),
	)

	err = g.Error()
	if err != nil {
		return
	}

	for _, v := range ret {
		v.Age = mAge[v.StudentID]
		v.Name = mName[v.StudentID]
		v.Phone = mPhone[v.StudentID]

		ret[v.StudentID] = v
	}

	return
}

// 具体业务逻辑
func f(_ *gin.Context, stuId Student) (stu Student, err error) {
	num := RandNum(200)
	time.Sleep(time.Millisecond * time.Duration(num))
	fmt.Println(fmt.Sprintf("sleep:%dms", num))

	// test panic
	if num > 150 {
		// test panic
		//a := 0
		//fmt.Println(1 / a)
	}

	if num > 150 {
		//err = base.Error{ErrNo: -1, ErrMsg: "获取学生信息失败!"}
		//zlog.Errorf(ctx, "获取学生信息失败,err=[%v],studentId=[%v]", err, stuId.StudentID)
		//return stu, err
	}

	stu = Student{StudentID: stuId.StudentID}
	return
}
