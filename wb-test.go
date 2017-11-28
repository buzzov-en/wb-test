package main

import(
	"fmt"
	"io/ioutil"
	"net/http"	
	"sync"
	"strings"
)

func main(){
	var wg = sync.WaitGroup{} 			//использую это для синхронизации созданных горутинов
	var str string
	queue := make(chan string) 			//канал, на который я посылаю "задания"
	total := 0
	i, k := 0, 5				        //тут задается параметр k

	for true {					//читаю stdin построчно, пока он не закончится
		_, err := fmt.Scanln(&str)  
		if err != nil { break }		
		if i < k { 				//под каждый новый url создаю новый горутин, но не больше, чем k
			go myTask(queue, &total, &wg) 	//таким образом, если k = 50, а строк с url 100, то будет создано 50 горутин
			i++				//но если k = 1000, а строк 100, то будет создано 100 горутин 	
		}
		queue <- str				//отправляю url обрабатываться на свободную горутину
	}
	close(queue)
	wg.Wait()                             		 //жду, пока созданные горутины закончат работать
	defer fmt.Println("Total", total)   		 //вывожу общий результат
}

func myTask( queue <- chan string, total *int, wg *sync.WaitGroup){      //в запущенном горутине функция слушает канал и принимает URL из очереди
	//fmt.Println("routine created")
	wg.Add(1)					//для синхронизации с main(), даем понять что появился работающий горутин
	for job:= range queue{												
		*total = *total + counter(job)			//для отправки запроса и подсчета тут вызывается отдельная функция
	}
	wg.Done()					 //для синхронизации с main(), даем понять, что работа горутина закончена
}

func counter(url string) int{				//функция, совершающая отправку запроса и подсчет вхождений строки
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Incorrect input")
		return 0
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't process html")
		return 0
	}
	defer res.Body.Close()  
	s := string(body[:])
	n := strings.Count(s, "Go")
	fmt.Println("Count for", url + ":", strings.Count(s, "Go")) 	//вывод результата выполнения функции подсчета вхождений строки
	return n;
}
