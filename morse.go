package main

import (
	"github.com/albenik/go-serial"
	"time"
)

type Morse struct {
	port    serial.Port
	msgChan chan string
	wpm     float64
}

func NewMorse(port serial.Port, msgChan chan string, wpm float64) *Morse {
	return &Morse{port, msgChan, wpm}
}

func (m *Morse) Run() {
	for {
		msg := <-m.msgChan
		m.send(msg)
	}
}

func (m *Morse) send(msg string) {
	for _, c := range msg {
		m.sendChar(c)
	}
}

func (m *Morse) dit() {
	m.port.SetDTR(true)
	m.waitDit()
	m.port.SetDTR(false)
	m.waitDit()
}

func (m *Morse) dah() {
	m.port.SetDTR(true)
	m.waitDah()
	m.port.SetDTR(false)
	m.waitDit()
}

func (m *Morse) waitDit() {
	time.Sleep(time.Duration(1200.0/m.wpm) * time.Millisecond)
}

func (m *Morse) waitDah() {
	time.Sleep(time.Duration(3600.0/m.wpm) * time.Millisecond)
}

func (m *Morse) sendChar(char rune) {
	switch char {
	case 'a', 'A':
		m.dit()
		m.dah()
	case 'b', 'B':
		m.dah()
		m.dit()
		m.dit()
		m.dit()
	case 'c', 'C':
		m.dah()
		m.dit()
		m.dah()
		m.dit()
	case 'd', 'D':
		m.dah()
		m.dit()
		m.dit()
	case 'e', 'E':
		m.dit()
	case 'f', 'F':
		m.dit()
		m.dit()
		m.dah()
		m.dit()
	case 'g', 'G':
		m.dah()
		m.dah()
		m.dit()
	case 'h', 'H':
		m.dit()
		m.dit()
		m.dit()
		m.dit()
	case 'i', 'I':
		m.dit()
		m.dit()
	case 'j', 'J':
		m.dit()
		m.dah()
		m.dah()
		m.dah()
	case 'k', 'K':
		m.dah()
		m.dit()
		m.dah()
	case 'l', 'L':
		m.dit()
		m.dah()
		m.dit()
		m.dit()
	case 'm', 'M':
		m.dah()
		m.dah()
	case 'n', 'N':
		m.dah()
		m.dit()
	case 'o', 'O':
		m.dah()
		m.dah()
		m.dah()
	case 'p', 'P':
		m.dit()
		m.dah()
		m.dah()
		m.dit()
	case 'q', 'Q':
		m.dah()
		m.dah()
		m.dit()
		m.dah()
	case 'r', 'R':
		m.dit()
		m.dah()
		m.dit()
	case 's', 'S':
		m.dit()
		m.dit()
		m.dit()
	case 't', 'T':
		m.dah()
	case 'u', 'U':
		m.dit()
		m.dit()
		m.dah()
	case 'v', 'V':
		m.dit()
		m.dit()
		m.dit()
		m.dah()
	case 'w', 'W':
		m.dit()
		m.dah()
		m.dah()
	case 'x', 'X':
		m.dah()
		m.dit()
		m.dit()
		m.dah()
	case 'y', 'Y':
		m.dah()
		m.dit()
		m.dah()
		m.dah()
	case 'z', 'Z':
		m.dah()
		m.dah()
		m.dit()
		m.dit()
	case '0':
		m.dah()
		m.dah()
		m.dah()
		m.dah()
		m.dah()
	case '1':
		m.dit()
		m.dah()
		m.dah()
		m.dah()
		m.dah()
	case '2':
		m.dit()
		m.dit()
		m.dah()
		m.dah()
		m.dah()
	case '3':
		m.dit()
		m.dit()
		m.dit()
		m.dah()
		m.dah()
	case '4':
		m.dit()
		m.dit()
		m.dit()
		m.dit()
		m.dah()
	case '5':
		m.dit()
		m.dit()
		m.dit()
		m.dit()
		m.dit()
	case '6':
		m.dah()
		m.dit()
		m.dit()
		m.dit()
		m.dit()
	case '7':
		m.dah()
		m.dah()
		m.dit()
		m.dit()
		m.dit()
	case '8':
		m.dah()
		m.dah()
		m.dah()
		m.dit()
		m.dit()
	case '9':
		m.dah()
		m.dah()
		m.dah()
		m.dah()
		m.dit()
	case '.':
		m.dit()
		m.dah()
		m.dit()
		m.dah()
		m.dit()
		m.dah()
	case ',':
		m.dah()
		m.dah()
		m.dit()
		m.dit()
		m.dah()
		m.dah()
	case '?':
		m.dit()
		m.dit()
		m.dah()
		m.dah()
		m.dit()
		m.dit()
	case '=':
		m.dah()
		m.dit()
		m.dit()
		m.dit()
		m.dah()
	case ' ':
		m.waitDah()
	default:
		return
	}
	m.waitDah()
}
