package init

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

const (
	GetModelName  string = "GetModelName"
	GetHandleFunc string = "GetHandleFunc"
)

type SignalSend struct {
	Subscribe       map[string][]chan interface{}
	SubscribeMaster map[string]chan interface{}
	Context         context.Context
	SizeOfChan      int
}

//var SignalAPP *SignalSend

var SignalAPPTest = InitSignal()

func InitSignal() *SignalSend {
	SignalAPP := new(SignalSend)
	SignalAPP = &SignalSend{
		Subscribe:       make(map[string][]chan interface{}, 10),
		SubscribeMaster: make(map[string]chan interface{}, 10),
		SizeOfChan:      20,
	}

	return SignalAPP
}

func (r *SignalSend) SetContext(ctx context.Context) {
	r.Context = ctx
}

// SendSignal 接收要发送的数据,数据类型为结构体地址
func (r *SignalSend) SendSignal(data interface{}) {
	modelName := reflect.TypeOf(data).Elem().Name()
	if _, ok := r.Subscribe[modelName]; !ok {
		zap.L().Error("该模型未被chan订阅，请先订阅在使用", zap.String("model_name", modelName))
		//return errors.New("该模型未被chan订阅，请先订阅在使用")
		//panic("该模型未被chan订阅，请先订阅在使用")
	}
	for _, ch := range r.Subscribe[modelName] {
		ch <- data
	}
	return
}

// Register 初始化订阅
func (r *SignalSend) Register(obj interface{}, name string) error {
	modelName, ok := obj.(string)
	if !ok {
		//zap.L().Error("Register初始化 断言错误，请注意注册的类以string类型返回")
		panic("Register初始化 断言错误，请注意注册的类以string类型返回")
	}
	if _, ok = r.SubscribeMaster[name]; !ok {
		ch := make(chan interface{}, r.SizeOfChan)
		r.SubscribeMaster[name] = ch
		if _, ok = r.Subscribe[modelName]; !ok {
			r.Subscribe[modelName] = make([]chan interface{}, 0)
			r.Subscribe[modelName] = append(r.Subscribe[modelName], ch)
			return nil
		} else {
			r.Subscribe[modelName] = append(r.Subscribe[modelName], ch)
		}
	}
	return nil
}

// RegisterFunc 订阅者传入订阅结构体与处理函数
func (r *SignalSend) RegisterFunc(obj interface{}, name string, fun func(data interface{}) error) {
	err := r.Register(obj, name)
	if err != nil {

		return
	}
	go r.Handler(name, fun)
}

// Handler 处理函数
func (r *SignalSend) Handler(name string, fun func(data interface{}) error) {
	ch := r.SubscribeMaster[name]
	for {
		select {
		case data := <-ch:
			go func() {
				_ = fun(data)
			}()
		case <-r.Context.Done():
			fmt.Println("上下文结束")
			return
		}
	}
}

func InitSignalHandle(ctx context.Context, SignalGroupAPP interface{}) *SignalSend {
	SignalAPPTest.SetContext(ctx)
	t := reflect.TypeOf(SignalGroupAPP).Elem()
	v := reflect.ValueOf(SignalGroupAPP).Elem()
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			for j := 0; j < v.Field(i).NumField(); j++ {
				v1 := v.Field(i).Field(j).Type()
				obj, ok := v1.MethodByName(GetModelName)
				if !ok {
					panic("signals init err!")
				}
				resultObj := obj.Func.Call([]reflect.Value{v.Field(i).Field(j)})
				funHandle, ok := v1.MethodByName(GetHandleFunc)
				if !ok {
					panic("signals init err! ")
				}
				resultHandle := funHandle.Func.Call([]reflect.Value{v.Field(i).Field(j)})
				SignalAPPTest.RegisterFunc(resultObj[0].Interface(), v1.Name(), resultHandle[0].Interface().(func(data interface{}) error))
			}
		}
	}
	return SignalAPPTest

}
