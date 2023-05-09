//Actividad 5.2 Programación Concurrente
//Andrés Alejandro Ramírez Fernández A00831316
//José Carlos de la Torre Hernández A01235953
//Julian Lawrence Gil Soares A00832272

/*Reflexión:
4. Tiempo
Los resultados que obtuvimos al calcular el tiempo fueron interesantes, sacamos el promedio de 10 
ejecucciones y encontramos que en promedio la funcion que corre de manera concurrente tarda 83.5 ms. 
Algo interesante que observamos en nuestro codigo es que en maquinas con menos cores la funcion concurrente
puede llegar a tardar mas en ejecutar que la funcion sequencial esto se debe a que 
si el numero de threads que lanza el programa excede el numero de cores de la computadora
en realidad no se corre de manera concurrente, secciona una core(o las que sean necesarias) y lo corre secuencial.
5. Complejidad
La complejidad del programa es O(n), si tomamons el archivo como n. Cada vez que corre la funcion, 
lexer, solo se recorre con un for de manera sequencial. Cuando se corre de manera concurrente, la complejidad es O(n)/num de threads. 
Implicación éticas:
El resaltador de sintaxis tiene un implicacion etica neutral. Ya que lo unico que hace es facilitar la codificacion de un progrma, 
los programas que ayudara a que crear pueden ser tanto buenos como malos, pero eso no esta en el alcance de este. Depende de como se usa el codigo.
*/


package main
 
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
  "log"
  "time"
)

//Función para leer cada archivo presente en el directorio
func dentroFolder(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func main() {
	in := os.Args[1] //Argumento que guarda los directorios de entrada
  out := os.Args[2] //Argumento que guarda los directorios de salida
  
  var archivo []string
	r := filepath.Walk(in, dentroFolder(&archivo))
	if r != nil {
		panic(r)
	}
	cont := 1

  //Sequential
 start := time.Now()
	for _, path := range archivo {
		if path != in {
			fmt.Println(path)
      lexerS(path, out, cont)
			cont++
		}
	}
  duration := time.Since(start)
  fmt.Println(duration, "Sequential")
  
	wg := sync.WaitGroup{}
  duration = 0
  //Concurrent
  start = time.Now()
	for _, path := range archivo {
		if path != in {
			wg.Add(1)
			fmt.Println(path)
			go lexer(path, out, cont, &wg)
			cont++
			wg.Wait()
		}
	}
  duration = time.Since(start)
  fmt.Println(duration, "Concurrent")
}

//recibe un string que contiene la letra que se va a comparar
//dependiendo de que carácter recibe regresa un número que //representa uno de los índices necesarios para trasladarnos a lo //largo de la matriz
func filter(c string)(int){
  
    if(c == "\n"){
      return 10
    }else if ( c == " " ){
      return 3
    }else if (c =="0" || c =="1" || c =="2" || c =="3" || c =="4" || c =="5" || c =="6" || c =="7" || c =="8" || c =="9"){
      return 0
    }else if (c >= "a" && c <= "z") || (c >= "A" && c <= "Z"){
      return 1
    }else if (c == "."){
      return 2
    }else if (c == "+" || c == "-" || c == "*" || c == "/" || c == "^" || c == "=" || c == "<" || c == ">" || c == "!"){
      return 4
    }else if (c == "[" || c == "]" || c == "," || c == "(" || c == ")" || c == "{" || c == "}"){
      return 5 
    }else if (c == "_"){
      return 6
    }else if (c == "#"){
      return 7
    }else if (c == "\"" || c=="'"){
      return 8
    }else if(c == ":"){
      return 9
    }else if(c == "\t"){
      return 11
    }else{
      return 12
    }
}

//Funcion para el resaltado de forma concurrente
func lexer(filename string, out string, num_file int, wg *sync.WaitGroup){

  numAr := fmt.Sprint(num_file)
  conq := "CON"
  out += "/output_" + numAr + conq + ".html"

	f, _ := os.Open(filename) // Archivo de Entrada
	fo, _ := os.Create(out) // Archivo de Salida
  
  
  fo.WriteString("<pre> <!DOCTYPE html> <html><style>body{font-size: large;background-color: #202021; font-family: mono;color: white;}</style> <head> <meta charset=\"utf-8\"> <meta name=\"viewport\" content=\"width=device-width\"> </head> <body> <div color: #ffffff> José Carlos de la Torre Hernández A01235953</div> <div color: #ffffff>Julian Gil A00832272</div><div color: #ffffff>Andres Ramirez A00831316</div><div color: #ffffff >Actividad Integradora 5.3</div><pre>")
  
    lineScanner := bufio.NewScanner(f) // Scanner para el Input
	lineScanner.Split(bufio.ScanLines) // Scanea Lineas   
  
  
  //tokens
  NUM := 100 //numbers
  OPE := 101 //operators
  SYM := 102 //symbols
  VAR := 103 //variable
  COM := 104 //coments
  INQ := 105 //in quotes
  COL := 106 //Colon
  END := 107 //end of entry
  NLN := 108 //new line
  ERR := 200 //error
  SPA := 110 //empty space
  TAB := 111 //tab

//    Num  var   .   " "  OPE   (    _    #    ""   :    \n  Tab  Empty
  MT := [7][13] int{   
    {  1,   2, SYM, SPA, OPE, SYM, SYM,   3,   4, COL, NLN, TAB, END}, // 0- Initial
    {  1, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM}, // 1- Numero
    {VAR,   2, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR}, // 2- Variable 
    {  3,   3,   3,   3,   3,   3,   3,   3,   3,   3, COM,   3,   3}, // 3- Comment
    {  4,   4,   4,   4,   4,   4,   4,   4,   6,   4,   4,   4,   4}, // 4- Quote
    {ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, END}, // 5- ERROR
    {INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ}, // 6- quote
  }

    state := 0
    lex := ""
    flag := true
    lista := []string{}
    cont:= 0
    con := []int{}
    Max := 0
  
  // Inicio del Scanner con For, Scanea Runas del Documento
	for lineScanner.Scan() {
        // Convierte las Runas a String
		line := lineScanner.Text()

        // For que le caracter x caracter
		for _, rune_l := range line {
			letter := string(rune_l)
      lista = append(lista, letter)
      cont++
    }
    lista =append(lista, "\n")
    con = append(con, cont)
    Max += cont
    cont = 0
  }
  
  for i:=0; i<len(lista); i++{
    k := lista[i]
    if state < 100{
      state = (MT[state][filter(k)])
      if state >= 100 && flag == true{ 
        flag = false
      }else{ 
        flag = true
      }
      if flag == true{
        if k != "\n" {
          lex += k
        }
      }
    }

    if state >=100{
      if state != 500{
        if state == NUM {
				  fo.WriteString("<font color=\"D36A0\">"+lex+"</font>")
			  }else if state == OPE {
          fo.WriteString("<font color=\"0098D6\">"+lex+"</font>")
			  } else if state == SYM {
          fo.WriteString("<font color=\"7E455B\">" + lex + "</font>") 
			  }else if state == VAR {
				  if lex == "import" || lex == "read" || lex == "try" || lex == "write" || lex == "for" || lex == "while" || lex == "switch" || lex == "case" || lex == "in" || lex == "" || lex == "and" || lex == "not" || lex == "or" || lex == "if" || lex == "elif" || lex == "else" || lex == "def" || lex == "print" || lex == "return" || lex == "this" || lex == "False" || lex == "True"|| lex == "class"|| lex == "except" || lex == "break" || lex == "append" || lex == "len" || lex == "range" || lex == "bool" || lex == "str" || lex == "int" || lex == "float" || lex == "except" || lex == "self" || lex == "global" || lex == "pow" || lex == "pass" || lex == "del"{
				    fo.WriteString("<font color=\"0DF2C2\">" + lex + "</font>") 
				  } else {
					  fo.WriteString("<font color=\"FFFFFF\">" + lex + "</font>")
				  }
        }else if state == COM {
				  fo.WriteString("<font color=\"BEC30A\">" + lex + "</font>") 
			  }else if state == INQ {
				  fo.WriteString("<font color=\"04EC19\">" + lex + "</font>") 
			  }else if state == COL {
				  fo.WriteString("<font color=\"#ffffff\">" + lex + "</font>") 
			  }else if state == NLN {
          fo.WriteString("<br>")
			  }else if state == SPA {
				  fo.WriteString(lex)
			  }else if state == TAB {
				  fo.WriteString("<t>")
			  }else if state == END {
				  fo.WriteString(lex) 
			  }else if state == ERR {
				  fo.WriteString(lex) 
			  }
      state = 0
      lex=""
      }
    if flag == false{
      i--
    }
    }
	}
  fo.WriteString("</body></html>")
  fo.Close()
  wg.Done()
}

//Funcion para el resaltado de forma secuencial
func lexerS(filename string, out string, num_file int){

  numAr := fmt.Sprint(num_file)
  seq := "SEQ"
  out += "/output_" + numAr + seq + ".html"

	f, _ := os.Open(filename) // Archivo de Entrada
	fo, _ := os.Create(out) // Archivo de Salida
  
  fo.WriteString("<pre> <!DOCTYPE html> <html><style>body{font-size: large;background-color: #202021; font-family: mono;color: white;}</style> <head> <meta charset=\"utf-8\"> <meta name=\"viewport\" content=\"width=device-width\"> </head> <body> <div style = color: #ffffff> José Carlos de la Torre Hernández A01235953</div> <div color: #ffffff>Julian Gil A00832272</div><div color: #ffffff>Andres Ramirez A00831316</div><div color: #ffffff >Actividad Integradora 5.3</div><pre>")
  
    lineScanner := bufio.NewScanner(f) // Scanner para el Input
	lineScanner.Split(bufio.ScanLines) // Scanea Lineas   
  
  
  //tokens
  NUM := 100 //numbers
  OPE := 101 //operators
  SYM := 102 //symbols
  VAR := 103 //variable
  COM := 104 //coments
  INQ := 105 //in quotes
  COL := 106 //Colon
  END := 107 //end of entry
  NLN := 108 //new line
  ERR := 200 //error
  SPA := 110 //empty space
  TAB := 111 //tab

//    Num  var   .   " "  OPE   (    _    #    ""   :    \n  Tab  Empty
  MT := [7][13] int{   
    {  1,   2, SYM, SPA, OPE, SYM, SYM,   3,   4, COL, NLN, TAB, END}, // 0- Initial
    {  1, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM, NUM}, // 1- Numero
    {VAR,   2, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR, VAR}, // 2- Variable 
    {  3,   3,   3,   3,   3,   3,   3,   3,   3,   3, COM,   3,   3}, // 3- Comment
    {  4,   4,   4,   4,   4,   4,   4,   4,   6,   4,   4,   4,   4}, // 4- Quote
    {ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, ERR, END}, // 5- ERROR
    {INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ, INQ}, // 6- quote
  }

    state := 0
    lex := ""
    flag := true
    lista := []string{}
    cont:= 0
    con := []int{}
    Max := 0
  
  // Inicio del Scanner con For, Scanea Runas del Documento
	for lineScanner.Scan() {
        // Convierte las Runas a String
		line := lineScanner.Text()

        // For que le caracter x caracter
		for _, rune_l := range line {
			letter := string(rune_l)
      lista = append(lista, letter)
      cont++
    }
    lista =append(lista, "\n")
    con = append(con, cont)
    Max += cont
    cont = 0
  }
  
  for i:=0; i< len(lista); i++{
    k := lista[i]
    //fmt.Println(k)
    if state < 100{
      state = (MT[state][filter(k)])
      if state >= 100 && flag == true{ 
        flag = false
      }else{ 
        flag = true
      }
      if flag == true{
        if k != "\n" {
          lex += k
        }
      }
    }

    //Si el estado es mayor a 100 entonces escribe a el HTML que creamos
    if state >=100{
      if state != 500{
        if state == NUM {
				  fo.WriteString("<font color=\"D36A0\">"+lex+"</font>")
			  }else if state == OPE {
          fo.WriteString("<font color=\"0098D6\">"+lex+"</font>")
			  } else if state == SYM {
          fo.WriteString("<font color=\"7E455B\">" + lex + "</font>") 
			  }else if state == VAR {
				  if lex == "import" || lex == "read" || lex == "try" || lex == "write" || lex == "for" || lex == "while" || lex == "switch" || lex == "case" || lex == "in" || lex == "" || lex == "and" || lex == "not" || lex == "or" || lex == "if" || lex == "elif" || lex == "else" || lex == "def" || lex == "print" || lex == "return" || lex == "this" || lex == "False" || lex == "True"|| lex == "class"|| lex == "except" || lex == "break" || lex == "append" || lex == "len" || lex == "range" || lex == "bool" || lex == "str" || lex == "int" || lex == "float" || lex == "except" || lex == "self" || lex == "global" || lex == "pow" || lex == "pass" || lex == "del"{
				    fo.WriteString("<font color=\"0DF2C2\">" + lex + "</font>") 
				  } else {
					  fo.WriteString("<font color=\"FFFFFF\">" + lex + "</font>")
				  }
        }else if state == COM {
				  fo.WriteString("<font color=\"BEC30A\">" + lex + "</font>") 
			  }else if state == INQ {
				  fo.WriteString("<font color=\"04EC19\">" + lex + "</font>") 
			  }else if state == COL {
				  fo.WriteString("<font color=\"#ffffff\">" + lex + "</font>")  
			  }else if state == NLN {
          fo.WriteString("<br>")
			  }else if state == SPA {
				  fo.WriteString(lex)
			  }else if state == TAB {
				  fo.WriteString("<t>")
			  }else if state == END {
				  fo.WriteString(lex) 
			  }else if state == ERR {
				  fo.WriteString(lex) 
			  }
      state = 0
      lex=""
      }
    if flag == false{
      i--
    }
    }
	}
  
  fo.WriteString("</body></html>")
  fo.Close()
}
