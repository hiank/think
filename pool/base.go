package pool

//Timer 计时器
type Timer interface {

	Update()				//NOTE: 更新状态
	SetInterval(int64)		//NOTE: 设置超时时间间隔
	TimeOut() bool 			//NOTE: 判断是否超时
}

// //Identifier 身份信息
// type Identifier interface {

// 	GetKey() string				//NOTE: 
// 	GetToken() string	 		//NOTE:		
// }