package signals

const (
	GetModelName  string = "GetModelName"
	GetHandleFunc string = "GetHandleFunc"
)

//type SignalRigsterGroup struct {
//	redis.RedSignalGroup
//	mysql.MysqlSignalGroup
//}
//
//var SignalGroupAPP = new(SignalRigsterGroup)

//var SignalAPPTest = initialize.InitSignalHandle(SignalGroupAPP)
//
//
//
//func main() {
//	ctx, cancle := context.WithCancel(context.Background())
//	defer cancle()
//	InitSignalHandle(ctx)
//
//	time.Sleep(5 * time.Second)
//	User := test.SchoolDao{
//		Id:        0,
//		ScholName: "重大",
//	}
//	SignalAPPTest.SendSignal(User)
//	//SignalAPPTest.Handler(user)
//	time.Sleep(3 * time.Second)
//	cancle()
//	time.Sleep(10 * time.Second)
//
//}
