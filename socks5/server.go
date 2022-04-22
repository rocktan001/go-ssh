// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"flag"
	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"log"
	"os/exec"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
	cmd  *exec.Cmd
}

func (p *program) Start(s service.Service) error {

	root, _ := osext.ExecutableFolder()
	p.cmd = exec.Command(root + "/socks5.exe")
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	err := p.cmd.Run()
	if err != nil {
		logger.Warningf("Error running: %v", err)
	}
	<-p.exit
	return nil

}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	close(p.exit)
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "socks5",
		DisplayName: "socks5",
		Description: "socks5 rocktan001",
		// Dependencies: []string{
		//     "Requires=network.target",
		//     "After=network-online.target syslog.target"},
		Option: options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
