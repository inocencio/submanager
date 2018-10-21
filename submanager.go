package main

/**
Autor: LandStalker (inocencio)
Data: 02/09/2018
Descricao: Programa para adicionar ou subtrair tempo nas legendas em formato SRT
           Uso: ./submanager -file:<arquivo_legenda.srt> -time=<tempo_em_ms>
           Ex.: ./submanager -file:"lawrence of arabia.srt" -time=-2000

EN Version: This is a short program to increment or reduce srt subtitle files. To use it call this program passing
these arguments ./submanager -file:<file.srt> -time=<time in ms>

Make sure you entered time argument as millesecond. Positive value will increase the time to subtitle shows up and
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
	"github.com/dustin/go-humanize"
	"github.com/paulrademacher/climenu"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
	"time"

	"github.com/logrusorgru/aurora"
)

//var srtfilepath string

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

	fmt.Println(aurora.Colorize(">> SubStr Manager v0.2 <<", aurora.GrayFg|aurora.GreenBg|aurora.BoldFm))

	if *timePtr == 0 && *filePtr == "" {
		//no arguments found, seek for video files and check the correspondent srt file
		files, err := ioutil.ReadDir("./")
		videos := make([]string, 0)
		subtitles := make([]string, 0)

		checkError(err)

		//find the videos out
		for _, f := range files {
			n := s.ToLower(f.Name())

			if !s.Contains(n,"sample") && (
				s.HasSuffix(n,".avi") ||
				s.HasSuffix(n,".mp4") ||
				s.HasSuffix(n, ".mkv") ||
				s.HasSuffix(n, ".mov")) {
					videos = append(videos, f.Name())
			}
		}

		//find srt files
		for _, v := range videos {
			if v != "" {
				ext := filepath.Ext(v)
				fn := s.Split(v, ext)[0]
				fn += ".srt"

				for _, f := range files {
					if s.EqualFold(f.Name(), fn) {
						subtitles = append(subtitles, f.Name())
					}
				}
			}
		}

		if len(subtitles) == 0 {
			fmt.Println(aurora.Red("No SRT files found!\n"),
				"\nMake sure there is a SRT file and this file must have the same name of video's file.")
			os.Exit(0)
		}

		//show srt files and menu
		fmt.Println("SRT file(s) found:")
		for _, f := range subtitles {
			if f != "" {
				//fmt.Println("File:          ", aurora.Cyan(f))
				fmt.Println(aurora.Cyan(f))
			}
		}
		fmt.Println()

		m := climenu.NewButtonMenu("","Choose an option to sync the SRT subtitle file")
		m.AddMenuItem("Custom", "custom")
		m.AddMenuItem("-2000 ms (rush)", "-2000")
		m.AddMenuItem("-1500 ms (rush)", "-1500")
		m.AddMenuItem("-1000 ms (rush)", "-1000")
		m.AddMenuItem(" -500 ms (rush)", "-500")
		m.AddMenuItem("  500 ms (delay)", "500")
		m.AddMenuItem(" 1000 ms (delay)", "1000")
		m.AddMenuItem(" 1500 ms (delay)", "1500")
		m.AddMenuItem(" 2000 ms (delay)", "2000")

		action, escaped := m.Run()

		if escaped {
			os.Exit(0)
		}
		fmt.Println()

		time := 0

		//time from menu or entry?
		if action != "custom" {
			time, _ = strconv.Atoi(action)
		} else {
			response := climenu.GetText("Enter time in milliseconds", "0")

			if response == "0" {
				fmt.Println("No entry time found!")
				os.Exit(0)
			}

			time, _ = strconv.Atoi(response)
		}

		//sync subtitles
		for _, f := range subtitles {
			if f != "" {
				strShifter(f, time)
			}
		}
	} else {
		//some arguments have been entered
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

		if s.HasSuffix(*filePtr, ".srt") {
			//todos os parametros OK? Entao converte o tempo do arquivo SRT.
			strShifter(*filePtr, *timePtr)
		} else {
			//arquivo invalido
			fmt.Println(aurora.Red("Error: O arquivo "), *filePtr, aurora.Red(" não é uma extensão srt."))
		}
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
@param filename, timestamp
*/
func strShifter(filename string, timestamp int) {
	fmt.Println("Time:          ", aurora.Cyan(timestamp), aurora.Cyan("ms"))
	fmt.Println("File:          ", aurora.Cyan(filename))

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

		for _, _e := range e.lines {
			buffer.WriteString(_e + br)
		}

		buffer.WriteString(br)
	}

	//cria o arquivo
	//file, err = os.Create(srtfilepath)
	file, err = os.Create(filename)
	//escreve no arquivo criado
	n, err := file.Write(buffer.Bytes())
	checkError(err)
	file.Sync()

	//fmt.Println("Written Bytes: ", aurora.Cyan(n))
	fmt.Println("KBytes Written:", aurora.Cyan(humanize.Bytes(uint64(n))))
}

func checkError(e error) {
	if e != nil {
		log.Println("Error: ", e)
		panic(e)
	}
}
