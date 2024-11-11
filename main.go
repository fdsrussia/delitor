package main

func main() {
	var usrAnsfer bool = true

	for usrAnsfer {
		UsrHello()             // приветствие
		usrAnsfer = UsrValue() // получение значения от пользователя
		if usrAnsfer {
			usrAnsfer = SearchFile() // поиск файла
		} else {
			break
		}
		if usrAnsfer {
			usrAnsfer = ScanFile()
		} else {
			break
		}
		if usrAnsfer == false {
			usrAnsfer = true
			continue
		}
		ResultScan()
		usrAnsfer = DelFiles()

	}
}
