package main

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"log"
	"strings"
	"github.com/alexflint/go-arg"

)



/*
Idea:


Takes a file with subdomains and their Ip separated by space
stores it it map, runs masscan on all ip and then match it back to domain to give results in the form -


domain:port
domain:port





*/

func check(err error){
	if err != nil {
        log.Fatal(err)
    }
}

var hip = make(map[string]string)

func masscan(masscan string){
	println("Running Masscan on hosts.")
	cmd := exec.Command("sudo",masscan, "-iL", "/tmp/masscan_input.txt", "--rate","10000", "-p","80,443,8080,8000,9001,3000,4443", "-oL", "/tmp/masscan_output.txt")
	err := cmd.Run()
	check(err)
	println("Masscan done.")
}


func main(){
	// Read file having domains and Ip
	
	var host,ip,port string
	var line []string
	var args struct {
		Input string
		Masscan string
	}
	arg.MustParse(&args)


	// readign the input file
	f,err := os.Open(args.Input)
	check(err)

	// creating a input file for masscan
	masscan_input,err := os.Create("/tmp/masscan_input.txt")
	check(err)

	w := bufio.NewWriter(masscan_input)
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		line = strings.Split(scanner.Text()," ")
		host,ip = line[0],line[1]
		hip[ip]=host
		_, err := w.WriteString(ip+"\n")
    	check(err)
	}
	w.Flush()
	masscan(args.Masscan) // Do masscan
	

	f,err = os.Open("/tmp/masscan_output.txt")
	scanner = bufio.NewScanner(f)
	for scanner.Scan(){
		line = strings.Split(scanner.Text()," ")
		if len(line) > 2  {
			port,ip = line[2], line[3]
			fmt.Println(hip[ip]+":"+port)
		}
	}

}
