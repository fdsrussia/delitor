package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var usrAnsfer bool = true   // отбивка о том, что все ок!
var myStr string = ""       // строка для перебора
var checkPoint bool = false // Нахождение чекпоинта, если чекпоинт найден, используется функция slcApnd()
var resultSlice []string = make([]string, 0)

// UsrHello - приветствие
func UsrHello() {
	fmt.Println("")
	fmt.Println("Программа Делитор версия 1.0.1 АО\"Академия Документооборота\"")
	fmt.Println("Предназначена для удаления файлов из файлового тома 1С")
	fmt.Println("Выгрузите отчет из Вашего 1С:Документооборот в формате txt, назовите его del")
	fmt.Println("Если Вы поместили файл del.txt в папку с программой, то введите цифру 1 для продолжения, или цифру 2 для выхода.")
}

// UsrCheckinput - проверка ввода пользователя
func UsrValue() bool {
	for true {
		fmt.Print("Введите значение: ")
		reader := bufio.NewReader(os.Stdin)  // Загоняем в буфер
		result, _ := reader.ReadString('\n') // парсим буфер на значение
		//Далее проверяем
		if strings.TrimSpace(result) == "1" {
			fmt.Println("Вы продолжили программу нажав 1") // Получили значение стартовать
			time.Sleep(time.Second * 1)
			break // Все ок - продолжаем программу
		} else if strings.TrimSpace(result) == "2" {
			fmt.Println("Вы завершили программу нажав 2")
			usrAnsfer = false // Выходим из программы
			break
		} else {
			fmt.Println("Вы ввели -->>", strings.TrimSpace(result), "<<-- повторите ввод!")
			time.Sleep(time.Second * 1)
			continue // Запрашиваем повторно ввод
		}
	}
	return usrAnsfer
}

// SearchFile - проверка нахождения файла в папке
func SearchFile() bool {
	checkTmp := true // для текущего цикала
	for checkTmp {
		filePath := "del.txt"                  // По условию файл называется - del.txt находится в папке программы
		absPath, err := filepath.Abs(filePath) // Находим текущий путь и делаем абсолбтный путь
		if err != nil {                        // Проверяем на ошибки
			panic(err)
		}
		file, err := os.Open(absPath) //Проверяем нахождение файла
		if err != nil {
			fmt.Println("Но файл с именем del.txt не найден в папке программы!")
			checkTmp = UsrValue()
		} else {
			fmt.Println("Файл найден..... начинаю анализ....\n")
			time.Sleep(time.Second * 1)
			fmt.Println("ожидайте.....\n")
			checkTmp = false
		}
		err = file.Close() // Закрываем файл
	}
	return usrAnsfer
}

// ScanFile - поиск вхождения "Отсутствуют данные в томе на диске"
func ScanFile() bool {
	filePath := "del.txt"                // По условию файл называется - del.txt находится в папке программы
	absPath, _ := filepath.Abs(filePath) // Находим текущий путь и делаем абсолбтный путь
	file, _ := os.Open(absPath)          //открываем файл
	scanner := bufio.NewScanner(file)    // начинаем сканировать построчно
	for scanner.Scan() {
		myStr = scanner.Text()
		if strings.Contains(myStr, "Лишние файлы (есть на диске, но сведения о них отсутствуют)") { // Если найходим "Отсутствуют данные в томе на диске" то добавляем в массив
			time.Sleep(time.Second * 1)
			fmt.Println("*Первый чекпоинт найден...")
			time.Sleep(time.Second * 1)
			continue
		} else if strings.Contains(myStr, "Имя	Полный путь	Время изменения	Размер") { // Если найходим "Отсутствуют данные в томе на диске" то добавляем в массив
			time.Sleep(time.Second * 1)
			fmt.Println(`**Второй чекпоинт найден, начинаю считать...`)
			time.Sleep(time.Second * 1)
			checkPoint = true
			continue
		} else if checkPoint { // Начинаю смотреть строки
			indexStar := strings.Index(myStr, "\\")
			if indexStar == -1 { // Если -1 то скорее всего пустая строка - такого быть не должно.
				fmt.Println("У Вас в файле пустая строка, Вы точно не меняли файл в ручную? 1С - не делает пустых строк! Выгрузите заного файл и повторите операцию.")
				time.Sleep(time.Second * 3)
				usrAnsfer = false
				break
			} else {
				strJob(myStr, indexStar)
			}
		}

	}

	return usrAnsfer
}

// strJob - функция сначала получает имя (поиск и обрезка до //) -> путь обрезка имени спереди и сзади (это и есть путь)
func strJob(myStr string, indexStar int) {
	nameFile := strings.TrimSpace(myStr[:indexStar])      // получаем имя файла
	indexLast := strings.LastIndex(myStr, nameFile)       // получаем последний индекс
	tmpWayResult := myStr[indexStar:indexLast] + nameFile // получаем полный путь файла
	resultSlice = append(resultSlice, tmpWayResult)
}

// Result Scan - просто инфа
func ResultScan() { // Просто инфа
	time.Sleep(time.Second * 1)
	fmt.Println("***Сканирование закончено")
	time.Sleep(time.Second * 1)
	fmt.Println("Найдено", len(resultSlice), "файлов.")
	time.Sleep(time.Second * 1)
	fmt.Println("--->>> Начинаю удаление файлов! <<<---")
}

// DelFiles - удаляемя файлы
func DelFiles() bool {
	checkDel := 0                   // считаем удаленне файлы
	checkNoDel := 0                 // считаем файлы не удаленные
	sliceNoDel := make([]string, 0) // сюда добавляем неудаленные файлы

	for i := 0; i < len(resultSlice); i++ {

		if err := os.Remove(resultSlice[i]); err != nil {
			sliceNoDel = append(sliceNoDel, resultSlice[i])
			checkNoDel++
			continue
		} else {
			checkDel++
		}
	}
	if len(sliceNoDel) <= 0 {
		fmt.Println("Удалено", checkDel, "файлов из", len(resultSlice))
		fmt.Println("Программа закроется через 5 секунд. Хорошего дня.")
		time.Sleep(time.Second * 7)

		usrAnsfer = true
	} else {
		fmt.Println("Списко файлов которые не удалось удалить:", checkNoDel, "из", len(resultSlice), ", показываю первые 10.")
		time.Sleep(time.Second * 2)
		for i := 0; i < 10; i++ {
			fmt.Println("Файл №", i+1, "имя -", sliceNoDel[i])
		}
	}
	return usrAnsfer

}
