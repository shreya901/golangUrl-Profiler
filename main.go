package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

//const ....
const (
	Links    string = "GET /links HTTP/1.1\r\nHost:my-worker-1.shreyasdomain.workers.dev\r\n Accept:application/json\r\n\r\n"
	Root     string = "GET / HTTP/1.1\r\nHost:my-worker-1.shreyasdomain.workers.dev\r\n Accept:text/html\r\n\r\n"
	Linksurl string = "http://my-worker-1.shreyasdomain.workers.dev/links"
	Rooturl  string = "http://my-worker-1.shreyasdomain.workers.dev"
)

func main() {
	info()
	flags()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "Simple Url profiler"
	app.Usage = "A simple Url Profiler Tool"
	app.Version = "1.0.0"

}

func flags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "url, u",
			Value: "http://my-worker-1.shreyasdomain.workers.dev",
			Usage: "Url to request",
		},
		&cli.StringFlag{
			Name:  "profile, p",
			Value: "1",
			Usage: "Number of times the url is requested to profile it",
		},
	}
}
func commands() {

	app.Commands = []*cli.Command{
		{
			Name:    "getResponse",
			Aliases: []string{"gR"},
			Usage:   "Requests the url passed in the url flag",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "url",
					Aliases: []string{"u"},
					Value:   "http://my-worker-1.shreyasdomain.workers.dev",
					Usage:   "Url to request",
				},
			},
			Action: func(c *cli.Context) error {

				// conn, err := net.Dial("tcp4", "my-worker-1.shreyasdomain.workers.dev:80")

				h, g, err := urlParse(c.String("url"))

				if err != nil || h == ":80" {
					return errors.New("Invalid url")
				}
				conn, err := net.Dial("tcp4", h)
				if err != nil {
					log.Println("dial error:", err)
					return err
				}

				defer conn.Close()
				fmt.Fprintf(conn, g)

				err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				if err != nil {
					log.Println("SetReadDeadline failed:", err)
					return err
				}

				recvBuf := make([]byte, 4096)
				length := 0
				for {
					n, err := conn.Read(recvBuf[:])
					length += n
					if err == io.EOF {
						break
					}
					if n == 0 {
						break
					}
					if err != nil {
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							log.Println("read timeout:", err)
							return err

						}
						log.Println("read error:", err)
						return err

					}
					fmt.Println(string(recvBuf))

				}
				if length == 0 {
					return errors.New("The provided path is hosted by the server")
				}
				return nil

			},
		},
		{
			Name:    "doProfile",
			Aliases: []string{"dP"},
			Usage:   "Profiles the url passed in the url flag by the times passed in the profile flag",

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "url",
					Aliases: []string{"u"},
					Value:   "http://my-worker-1.shreyasdomain.workers.dev",
					Usage:   "Url to request",
				},
				&cli.StringFlag{
					Name:    "profile",
					Aliases: []string{"p"},
					Value:   "1",
					Usage:   "Number of times the url is requested to profile it",
				},
			},
			Action: func(c *cli.Context) error {
				count, _ := strconv.Atoi(c.String("profile"))
				errorcodes := make([]string, count)
				minSize := math.MaxInt64
				maxSize := math.MinInt64
				var minTime int64
				var maxTime int64
				var totalTime int64
				var meanTime int64
				minTime = math.MaxInt64
				maxTime = math.MinInt64
				totalTime = 0
				meanTime = 0
				success := 0
				failure := 0

				h, g, err := urlParse(c.String("url"))

				if err != nil || h == ":80" {
					return errors.New("Invalid url")
				}
				for i := 0; i < count; i++ {

					s, t, l := makeRequest(h, g)
					totalTime += t
					if minSize > l {
						minSize = l
					}
					if maxSize < l {
						maxSize = l

					}
					if minTime > t {
						minTime = t
					}
					if maxTime < t {
						maxTime = t
					}
					s1 := strings.Split(s, " ")

					if s1[1] != "200" {
						failure++
						errorcodes = append(errorcodes, s1[1])
					} else {
						success++
					}

				}
				meanTime = totalTime / int64(count)
				successPercent := (float64(success) / float64(count)) * 100.0
				fmt.Printf("Total number of times the url is requested for Profiling : %v \n", count)
				fmt.Printf("Maxinum time for requesting: %v \n", maxTime)
				fmt.Printf("Mininum time for requesting: %v \n", minTime)
				fmt.Printf("Mean time for requesting: %v \n", meanTime)
				fmt.Printf("Maximum response length in bytes: %v \n", maxSize)
				fmt.Printf("Minimum response length in bytes: %v \n", minSize)
				fmt.Printf("Percentage of Success: %v \n", successPercent)
				fmt.Printf("Error codes for failures: %v \n", errorcodes)

				return nil
			},
		},
	}
}

func makeRequest(host string, getString string) (string, int64, int) {
	start := time.Now()
	s, t, l := func() (string, int64, int) {
		conn, err := net.Dial("tcp4", host)
		if err != nil {
			log.Println("dial error:", err)
			return "Dial Error", int64(time.Since(start) / time.Millisecond), 0

		}
		defer conn.Close()
		fmt.Fprintf(conn, getString)

		err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Println("SetReadDeadline failed:", err)
			return "Read Deadline failed", int64(time.Since(start) / time.Millisecond), 0
		}

		statusBuf := make([]byte, 15)
		n, _ := conn.Read(statusBuf[:])
		status := string(statusBuf)

		length := n
		recvBuf := make([]byte, 4096)
		for {
			n, err := conn.Read(recvBuf[:])
			length += n
			if err == io.EOF {
				break
			}
			if n == 0 {
				break
			}
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					log.Println("read timeout:", err)
					return "Read Timeout", int64(time.Since(start) / time.Millisecond), length

				}
				log.Println("read error:", err)
				return "Read Error", int64(time.Since(start) / time.Millisecond), length

			}

		}

		return status, int64(time.Since(start) / time.Millisecond), length

	}()
	return s, t, l

}

func urlParse(urlString string) (string, string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}
	if u.Path == "" {
		u.Path = "/"
	}

	return u.Host + ":80", "GET " + u.Path + " HTTP/1.1\r\nHost:" + u.Host + "\r\nAccept:text/html\r\n\r\n", nil
}
