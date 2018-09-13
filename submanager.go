package main

/**
Autor: LandStalker (inocencio)
Data: 02/09/2018
Descricao: Programa para adicionar ou subtrair tempo nas legendas em formato SRT
           Uso: ./submanager -file:<arquivo_legenda.srt> -time=<tempo_em_ms>
           Ex.: ./submanager -file:"lawrence of arabia.srt" -time=-2000

EN Version: This is a short program to increment or reduce srt subtitle files. To use it call this program passing
these arguments ./submanager -file:<file.srt> -time=<time in ms>

Make sure you entered time argument as millescond. Positive value will increate the time to subtitle shows up and
negative value will decrese this time instead.

                                        `:+shhso+/
                 ``                  ./hddy+..:/+o.
             `/osyyso+/:`          -odNdhy+  .yNd-+:
           .oddhhdddhysooso/-`  `-smmhs+//:   -oh-/+:
          /hdysdhs+/. `-yo+smhyyhmNNmdhhhhyso+/-` ohh:
        `-mhyyhmo///   -dMo--/////+syyss+///:/+osyydmNo.
    ````-ddysyms+///.`  .-.`-/+osys/:-:.``:so.```..-:+yso/.
 `:sdddhyNysyhN//////:.````-/syo/-```:so.``.:.`````````-+sh+`
.ymdyyhmNmsshdy///////////+sy+-.``````..``````````.-:+oooosys:
hmhyddhdmdsshd+//////////ys+-``````````........-/oys+/-..``.so:
Mdhhm-.+ddsyhh////////+sh+-```.-/+ooooooooossso+/-.`````````:ss
Mhhdd``+dhshhy///////oss-``.-/oo+:--------....``````````````.oy.
Ndhhm-.+dhsdho/////+ys/.`.+ss/-.`````````````````````````````+y/
omdhhmdmmhsdy+////oh+--/os/-.````````````````````````````````+h+
 +hmddmNNdsdh+//+sh/:oso:..``````````````````````````````````+y-
  `-+oosymshhy/oyyoso/-`````````````````````````````````````.oy
       `:Nyyhdhddy+-.```````````````````````````````````````:so
        .dhysymms:.````````````````````````````````````````.sy+
         .Nhyssyhhho:..`````````````````````````.--:://+osshdmy
          +hdyssssyyyyo+-..````````````...-:+osyyyyhhyyyyhddyo-
           :sddhyysssyyyhys+:.`````-:/oshdddhhhysssssyhhdh+`
            `.-oyyhhhhhyyyyhhhdhhhdhhhhhhhyyyyyhhhhhhyy/-`
                  `-:/+oyhddddddddddddhhyysso+////:-.
*/

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
	"time"

	"github.com/logrusorgru/aurora"
)

var srtfilepath string

type TextPart struct {
	num   string
	time  string
	lines []string
}

func main() {

	isNotFlag := false

	timePtr := flag.Int("time", 0, "O tempo é dado em millisegundos (1000ms = 1s) "+
		"Valores acima de 0 atrasa a legenda, abaixo de 0, adianta.")
	filePtr := flag.String("file", "", "Nome do arquivo SRT.")

	flag.Parse()

	//fmt.Print(aurora.Green(">> SubStr Manager v0.1 <<\n\n").Bold())
	fmt.Println(aurora.Colorize(">> SubStr Manager v0.1 <<", aurora.GrayFg|aurora.GreenBg|aurora.BoldFm))

	if *timePtr == 0 {
		fmt.Println(aurora.Red("Informe um 'time'. Ex: -time=-1000"))
		isNotFlag = true
	}

	if *filePtr == "" {
		fmt.Println(aurora.Red("Informe um 'file'. Ex: -file=\"minha legenda.srt\""))
		isNotFlag = true
	}

	if isNotFlag {
		os.Exit(1)
	}

	ex, err := os.Executable()
	checkError(err)

	fmt.Println("CDir:          ", aurora.Cyan(filepath.Dir(ex)))
	fmt.Println("Time:          ", aurora.Cyan(*timePtr), aurora.Cyan("ms"))
	fmt.Println("File:          ", aurora.Cyan(*filePtr))

	srtfilepath = *filePtr

	if s.HasSuffix(srtfilepath, ".srt") {
		//todos os parametros OK? Entao converte o tempo do arquivo SRT.
		strShifter(srtfilepath, *timePtr)
	} else {
		//arquivo invalido
		fmt.Println(aurora.Red("Error: O arquivo "), srtfilepath, aurora.Red(" não é uma extensão srt."))
	}
}

func (p *TextPart) convertTime(timestamp int) {
	times := s.Split(p.time, "-->")
	p.time = ""

	for counter, e := range times {
		timesdiv := s.Split(s.TrimSpace(e), ",")

		var hours int
		var minutes int
		var seconds int
		var ms int
		var err error

		for i, _e := range timesdiv {
			if i == 0 {
				//time completo
				t := s.Split(_e, ":")

				for j, v := range t {
					if j == 0 {
						//hora
						hours, err = convertStrToInt(v)
					} else if j == 1 {
						//minuto
						minutes, err = convertStrToInt(v)
					} else {
						//segundo
						seconds, err = convertStrToInt(v)
					}

					checkError(err)
				}
			} else {
				//restante em millisegundos
				ms, err = convertStrToInt(_e)
				//cria uma data no padrao Go Lang
				result := time.Date(2018, 9, 7, hours, minutes, seconds, 0, time.UTC)
				//adiciona os millisegundos restantes a data criada
				result = result.Add(time.Duration(ms) * time.Millisecond)
				result = result.Add(time.Duration(timestamp) * time.Millisecond)

				//formata o tempo ja convertido para o padrao do formato SRT com base num padrao definido
				//ref: https://golang.org/pkg/time/#example_Time_Format
				formattedTime := s.Replace(result.Format("15:04:05.000"), ".", ",", 1)

				if counter > 0 {
					p.time += " --> "
				}

				//adiciona o tempo convertido no lugar do anterior
				p.time += formattedTime
			}
		}
	}
}

/**
Converte uma string em inteiro.
@param srt
*/
func convertStrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

/**
Ajusta o tempo em millesegundos da legenda. Valores positivos, aumenta o tempo da legenda em relação ao
vídeo enquanto valores negativos reduz o tempo da legenda.
*/
func strShifter(filename string, timestamp int) {
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()

	//variaveis
	scanner := bufio.NewScanner(file)
	var line string
	var text []string
	var parts []TextPart
	partCounter := 0
	i := 0
	isNewPart := true
	textPart := new(TextPart)

	//le o arquivo linha a linha (para leitura direta usa-se ioutil.ReadFile(filename) no lugar)
	for scanner.Scan() {
		//pega a linha corrente do arquivo
		line = scanner.Text()

		//cria um novo trecho
		if isNewPart {
			i++
			isNewPart = false
			partCounter = 0
			textPart = new(TextPart)
		}

		//novo trecho do arquivo?
		if len(line) == 0 {
			textPart.lines = text
			text = nil

			//salva o trecho atual no slice
			parts = append(parts, *textPart)
			isNewPart = true
			continue
		}

		if partCounter == 0 {
			//0 = numeracao da legenda
			textPart.num = line
			partCounter = 1
		} else if partCounter == 1 {
			//1 = tempo da legenda
			textPart.time = line
			partCounter = 2
		} else if partCounter == 2 {
			//2 = linha(s)
			text = append(text, line)
		}
	}

	fmt.Println("Parts:         ", aurora.Cyan(i))

	//formata o texto de saida para o arquivo de saida
	var buffer bytes.Buffer
	br := "\n"

	for _, e := range parts {
		e.convertTime(timestamp)

		buffer.WriteString(e.num + br)
		buffer.WriteString(e.time + br)
		fmt.Println("entrada")

		for _, _e := range e.lines {
			buffer.WriteString(_e + br)
		}

		buffer.WriteString(br)
	}

	//cria o arquivo
	file, err = os.Create(srtfilepath)
	//escreve no arquivo criado
	n, err := file.Write(buffer.Bytes())
	file.Sync()

	fmt.Println("Written Bytes: ", aurora.Cyan(n))
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
