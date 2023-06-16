package logger

import "fmt"

func InitLogger() {

}

type Logger interface {
}

func Error(v ...any) {
	fmt.Println(v...)
}

func Info(v ...any) {
	fmt.Println(v...)
}
