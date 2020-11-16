package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// refer to:
// https://golang.org/src/net/mail/message.go?h=debugT#L36
var debug = debugT(false)

type debugT bool

func (d debugT) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

type userConfig struct {
	LoginURL   string
	LoginID    string
	LoginPass  string
	ChannelURL string
	CmdFile    string
}

func (c *userConfig) loadConfig() error {
	dir, err := getExecDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(c)

	return err
}

func main() {
	log.Println("------- start cmd2s -------")
	startTime := time.Now()

	funcTime := time.Now()
	// Load user config
	userConf := &userConfig{}
	err := userConf.loadConfig()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	debug.Printf("userConf: %+v", userConf)
	log.Printf("-> loadConfig(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	funcTime = time.Now()
	// Read command file
	cmds, err := readCommands(userConf.CmdFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf("-> cmdCount: %+v", len(cmds))
	log.Printf("-> readCommands(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	// Create new context for displaying with chrome
	// set options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
	)

	// create context
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	funcTime = time.Now()
	// Log in to slack via AAD
	err = chromedp.Run(taskCtx, loginToSlack(userConf.LoginURL, userConf.LoginID, userConf.LoginPass))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf("-> loginToSlack(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	funcTime = time.Now()
	// Send slash command to Slack
	err = chromedp.Run(taskCtx, sendCmdToSlack(userConf.ChannelURL, cmds))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf("-> sendCmdToSlack(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	funcTime = time.Now()
	// Check the command execution result
	var resp string
	err = chromedp.Run(taskCtx, checkResultToSlack(userConf.ChannelURL, `/feed list`, &resp))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	debug.Printf("checkResult: %+v", resp)
	log.Printf("-> checkResultToSlack(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	funcTime = time.Now()
	// Write check result
	err = writeCheckResult(resp)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf("-> writeCheckResult(): %+v sec\n", (time.Now().Sub(funcTime)).Seconds())

	log.Println("==========================")
	log.Printf("-> total: %+v sec\n", (time.Now().Sub(startTime)).Seconds())
	log.Println("------- end cmd2s -------")
}

func writeCheckResult(text string) error {
	dir, err := getExecDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, time.Now().Format("20060102_150405_")+"check_result.txt")

	err = ioutil.WriteFile(path, []byte(text), 0666)
	if err != nil {
		return err
	}

	return err
}

func checkResultToSlack(url string, cmd string, resp *string) chromedp.Tasks {
	const cmdForm string = `//*[contains(@class, 'ql-editor')]`
	const sendBtn string = `//button[@data-qa='texty_send_button']`
	const onlyMsg string = `//div[@class='c-virtual_list__item' and contains(@id, 'xxxxx')]`

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(sendBtn),
		chromedp.SendKeys(cmdForm, "/feed "),
		chromedp.SendKeys(cmdForm, "list"),
		chromedp.Click(sendBtn),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(onlyMsg),
		chromedp.Text(onlyMsg, resp),
	}
}

func sendCmdToSlack(url string, cmds []string) chromedp.Tasks {
	const cmdForm string = `//*[contains(@class, 'ql-editor')]`
	const sendBtn string = `//button[@data-qa='texty_send_button']`

	tasks := chromedp.Tasks{}
	tasks = append(tasks, chromedp.Navigate(url))
	tasks = append(tasks, chromedp.WaitVisible(sendBtn))
	tasks = append(tasks, chromedp.Navigate(url))
	tasks = append(tasks, chromedp.WaitVisible(sendBtn))

	// Run multiple commands
	for _, value := range cmds {
		cmd, args := strings.Split(value, " ", 2)
		tasks = append(tasks, chromedp.SendKeys(cmdForm, cmd+" "))
		tasks = append(tasks, chromedp.SendKeys(cmdForm, args))
		tasks = append(tasks, chromedp.Click(sendBtn))
	}

	return tasks
}

func loginToSlack(url string, id string, pass string) chromedp.Tasks {
	const idForm string = `//input[@id='email']`
	const passForm string = `//input[@id='password']`
	const signInBtn string = `//button[@id='signin_btn']`
	const sendBtn string = `//button[@data-qa='texty_send_button']`

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(idForm),
		chromedp.SendKeys(idForm, id),
		chromedp.SendKeys(passForm, pass),
		chromedp.Click(signInBtn),
	}
}

func readCommands(path string) ([]string, error) {
	cmds := []string{}
	f, err := os.Open(path)
	if err != nil {
		return cmds, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		text := strings.TrimSpace(sc.Text())
		debug.Printf("cmd: %+v", text)

		if strings.HasPrefix(text, "#") {
			continue
		}
		if text == "" {
			continue
		}
		cmds = append(cmds, text)
	}

	return cmds, err
}

func getExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}

func getConfigPath() (string, error) {
	dir, err := getExecDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), err
}
