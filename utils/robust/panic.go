package robust



//Panic 抛出异常
func Panic(err error) {

	if err != nil {
		panic(err)
	}
}
