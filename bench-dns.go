package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	file    = flag.String("file", "", "Local domain list")
	url     = flag.String("url", "", "URL to raw domain list")
	ip      = flag.String("ip", "", "DNS server IP address")
	iterate = flag.Int("iterate", 1, "Repeat counts")
)

// DNS lookup setup
var r = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(10000),
		}

		ip := *ip + ":53"
		return d.DialContext(ctx, network, ip)
	},
}

func DownloadFile(filepath string, url string) error {

	//
	path := "download"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func CheckDns(domain string, wg *sync.WaitGroup) {

	// iprecords, _ := net.LookupIP(domain)

	iprecords, _ := r.LookupHost(context.Background(), domain)
	for _, ip := range iprecords {
		fmt.Println("============= - " + domain)
		fmt.Println(ip)
	}

	r.LookupHost(context.Background(), domain)
	wg.Done()
}

func OpenFile(f string) {

	count := 0
	sum := 0
	comment := "#"

	// Open downloaded file
	file, err := os.Open(f)
	if err != nil {
		log.Fatalf("File does not exist. Please try to use url argument or create domain list manually.")
		fmt.Println("error opening file: err:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	//
	var wg sync.WaitGroup
	//

	for i := 0; i < *iterate; i++ {
		sum += i
		// loop lines from domains,txt
		for _, domain := range text {

			if !strings.Contains(domain, comment) {

				wg.Add(1)
				// Enable thread
				go CheckDns(domain, &wg)

				// Standard
				// CheckDns(domain, &wg)

			}
			count++
		}

		wg.Wait()
		// fmt.Println("=============")
		// fmt.Println("Sum iterations: " + fmt.Sprint(sum))
	}

	// Count lines in text file
	domains := 0
	for _, domain := range text {

		if !strings.Contains(domain, comment) {

			domains++

		}
	}

	fmt.Println("Domains: " + fmt.Sprint(domains))
	fmt.Println("Iterations: " + fmt.Sprint(*iterate))
	fmt.Println("Total checked domains (with iterations): " + fmt.Sprint(count))
}

func checkIPAddress(ip string) bool {
	valid := false
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		valid = true
	}
	return valid
}

func main() {

	flag.Parse()

	if *file == "" && *url == "" {
		fmt.Println("You need use one argument - url or file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *ip == "" {
		fmt.Println("Please set DNS IP address")
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		fmt.Println("Checking IP address")
		if checkIPAddress(*ip) {
			fmt.Println("IP address - OK")
		} else {
			fmt.Println("IP address does not correct. Exit.")
			os.Exit(1)
		}
	}

	if *file != "" && *url == "" {
		fmt.Println("Domain list file name: ", *file)
		OpenFile(*file)

	} else {
		fmt.Println("Will try download file from URL to download.txt in to script folder")
		fmt.Println("url: ", *url)

		// Download file
		err := DownloadFile("domains.txt", *url)
		if err != nil {
			// panic(err)
			fmt.Println("Unsupported URL format. Exit.")
			os.Exit(1)
		}

		OpenFile("domains.txt")
	}

}
