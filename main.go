package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ifeanyiBatman/uplink/config"
	"golang.ngrok.com/ngrok/v2"
	//"github.com/joho/godotenv"
)

func main() {
	var cfg config.Config
	
	args := os.Args[1:]
	if len(args) == 1 {
		if args[0] == "register" {
			register(&cfg)
			return
		}
		if args[0] == "users" {
			config.ListUsers()
			return
		}
	} else if len(args) == 2 {
		if args[0] == "login" {
			config.Login(args[1])
		}
	}
	
	
	port := flag.String("port", "8080", "Usage: -port (xxxx) default 8080")
	flag.Parse()
	

	
	
	username, err := config.GetCurrentUser()
	if err != nil {
		err = register(&cfg)
		}

	cfg, err = config.GetUserConfig(username)
	if err != nil {
		fmt.Printf("Error getting user err: %s", err)
		return
	}

	cfg.Port = *port

	// dotenv implementation
	//err := godotenv.Load()
	//if err != nil{
	//	fmt.Println(err)
	//}
	//cfg.domain, cfg.authToken = os.Getenv("NGROK_DOMAIN") , os.Getenv("NGROK_AUTHTOKEN")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = run(ctx, cfg)
		if err != nil {
			fmt.Printf("Error running the tunnel %v", err)
		}
	}()

	fmt.Println("Tunnel is running. Type 'quit' to shutdown.")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type quit to shutdown")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "quit" {
			fmt.Println("shutting tunnel down...")
			cancel()
			return
		}
	}
}

func run(ctx context.Context, cfg config.Config) error {
	agent, err := ngrok.NewAgent(ngrok.WithAuthtoken(cfg.AuthToken))
	if err != nil {
		fmt.Println("Error creating agent", err)
		return err
	}

	ln, err := agent.Forward(ctx,
		ngrok.WithUpstream(cfg.GetLocalhost()),
		ngrok.WithURL(cfg.Domain),
	)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	fmt.Printf("Your machine here %s has been forwarded to %s \n", cfg.GetLocalhost(), cfg.Domain)

	<-ln.Done()
	return nil
}

func register (cfg *config.Config) error {
	fmt.Println("Please setup Ngrok... https://dashboard.ngrok.com/get-started/setup")

	fmt.Print("Enter your authtoken: ")
	reader := bufio.NewReader(os.Stdin)
	authtoken,_ := reader.ReadString('\n')
	cfg.AuthToken = strings.TrimSpace(authtoken)
	
	fmt.Print("Enter your domain: ")
	reader = bufio.NewReader(os.Stdin)
	domain, err := reader.ReadString('\n')
	cfg.Domain = strings.TrimSpace(domain)
	
	fmt.Print("Enter your username: ")
	reader = bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	username = strings.TrimSuffix(username, "\n")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return err
	}
	err = cfg.SetUser(username)
	if err != nil {
		fmt.Printf("Error setting you up %s :(%v)\n", username,err)
		return err
	}
	fmt.Printf("your account have been setup %s", username)
	return nil
}